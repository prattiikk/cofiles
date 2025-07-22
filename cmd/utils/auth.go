package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// Config holds authentication data
type Config struct {
	JWT    string `json:"jwt,omitempty"`
	Server string `json:"server,omitempty"`
}

// AuthResponse from server
type AuthResponse struct {
	URL  string `json:"url"`
	Code string `json:"code"`
}

// TokenResponse from server
type TokenResponse struct {
	Token string `json:"token"`
}

const (
	defaultServer = "http://localhost:3000"
	pollInterval  = 2 * time.Second
	pollTimeout   = 5 * time.Minute
)

// getConfigPath returns path to config file
func getConfigPath() string {
	home, _ := os.UserHomeDir()

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(home, "AppData", "Roaming", "cofiles", "config.json")
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "cofiles", "config.json")
	default:
		return filepath.Join(home, ".config", "cofiles", "config.json")
	}
}

// LoadConfig loads configuration from file
func LoadConfig() *Config {
	config := &Config{Server: defaultServer}

	data, err := os.ReadFile(getConfigPath())
	if err != nil {
		return config // Return default if file doesn't exist
	}

	json.Unmarshal(data, config)
	if config.Server == "" {
		config.Server = defaultServer
	}

	return config
}

// SaveConfig saves configuration to file
func SaveConfig(config *Config) error {
	configPath := getConfigPath()

	// Create directory if it doesn't exist
	os.MkdirAll(filepath.Dir(configPath), 0755)

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

// IsAuthenticated checks if user has a token
func IsAuthenticated() bool {
	config := LoadConfig()
	return config.JWT != ""
}

// GetJWT returns the stored token
func GetJWT() (string, error) {
	config := LoadConfig()
	if config.JWT == "" {
		return "", fmt.Errorf("not authenticated. Please run login first")
	}
	return config.JWT, nil
}

// GetAuthHeader returns Bearer token header
func GetAuthHeader() (string, error) {
	jwt, err := GetJWT()
	if err != nil {
		return "", err
	}
	return "Bearer " + jwt, nil
}

// ClearAuth removes stored authentication
func ClearAuth() error {
	config := LoadConfig()
	config.JWT = ""
	return SaveConfig(config)
}

// openBrowser opens URL in default browser
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd, args = "rundll32", []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd, args = "open", []string{url}
	default:
		cmd, args = "xdg-open", []string{url}
	}

	return exec.Command(cmd, args...).Start()
}

// requestAuthURL gets authentication URL from server
func requestAuthURL(serverURL string) (*AuthResponse, error) {
	resp, err := http.Post(serverURL+"/auth/request-url", "application/json", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to request auth URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %d", resp.StatusCode)
	}

	var authResp AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	return &authResp, err
}

// pollForToken waits for user to complete auth and gets token
func pollForToken(serverURL, code string) (*TokenResponse, error) {
	client := &http.Client{}
	pollURL := fmt.Sprintf("%s/auth/token/%s", serverURL, code)
	deadline := time.Now().Add(pollTimeout)

	for time.Now().Before(deadline) {
		resp, err := client.Get(pollURL)
		if err != nil {
			time.Sleep(pollInterval)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			time.Sleep(pollInterval)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			var tokenResp TokenResponse
			if json.Unmarshal(body, &tokenResp) == nil && tokenResp.Token != "" {
				return &tokenResp, nil
			}
		}

		if resp.StatusCode == http.StatusAccepted {
			fmt.Print(".")
			time.Sleep(pollInterval)
			continue
		}

		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("authentication code expired")
		}

		time.Sleep(pollInterval)
	}

	return nil, fmt.Errorf("authentication timed out")
}

// Authenticate performs the complete authentication flow
func Authenticate(serverURL string) error {
	config := LoadConfig()
	if serverURL == "" {
		serverURL = config.Server
	}

	fmt.Println("🔐 Starting authentication...")

	// Get auth URL
	authResp, err := requestAuthURL(serverURL)
	if err != nil {
		return err
	}

	fmt.Printf("🌐 Opening browser: %s\n", authResp.URL)

	// Open browser
	if err := openBrowser(authResp.URL); err != nil {
		fmt.Printf("❌ Failed to open browser: %v\n", err)
		fmt.Printf("Please open: %s\n", authResp.URL)
	}

	fmt.Print("⏳ Waiting for authentication")

	// Wait for token
	tokenResp, err := pollForToken(serverURL, authResp.Code)
	if err != nil {
		fmt.Println()
		return err
	}

	fmt.Println("\n✅ Authentication successful!")

	// Save token
	config.JWT = tokenResp.Token
	config.Server = serverURL

	if err := SaveConfig(config); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("🎉 You are now logged in!")
	return nil
}

// VerifyAuth checks if stored token is still valid
func VerifyAuth() error {
	config := LoadConfig()
	if config.JWT == "" {
		return fmt.Errorf("not authenticated")
	}

	req, err := http.NewRequest("GET", config.Server+"/auth/verify", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+config.JWT)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to verify auth: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ClearAuth()
		return fmt.Errorf("token invalid, please login again")
	}

	return nil
}
