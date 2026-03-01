package auth

import (
	"bytes"
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

// start flow endpoint response
type AuthResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURL string `json:"verification_url"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// TokenResponse from polling endpoint server
type PollingResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

const (
	defaultServer = "http://ec2-43-205-235-230.ap-south-1.compute.amazonaws.com"
	pollInterval  = 5 * time.Second
	pollTimeout   = 5 * time.Minute
)

// getConfigPath returns path to config file stored at respective locations based on OS
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
func startAuthFlow(serverURL string) (*AuthResponse, error) {
	fmt.Println("inside startAuthFlow")
	resp, err := http.Post(serverURL+"/cli/device/start", "application/json", nil)
	fmt.Printf("Received response: %v, error: %v\n", resp, err)
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
func pollForToken(serverURL, deviceCode string) (*PollingResponse, error) {
	client := &http.Client{}
	pollURL := fmt.Sprintf("%s/cli/device/token", serverURL)
	deadline := time.Now().Add(pollTimeout)

	for time.Now().Before(deadline) {

		// Build request body
		payload := map[string]string{
			"device_code": deviceCode,
		}

		jsonBody, _ := json.Marshal(payload)

		req, err := http.NewRequest("POST", pollURL, bytes.NewBuffer(jsonBody))
		if err != nil {
			time.Sleep(pollInterval)
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
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

		// SUCCESS → Token received
		if resp.StatusCode == http.StatusOK {
			var tokenResp PollingResponse
			if json.Unmarshal(body, &tokenResp) == nil && tokenResp.AccessToken != "" {
				return &tokenResp, nil
			}
		}

		// Authorization pending (your backend returns 428)
		if resp.StatusCode == 428 {
			fmt.Print(".")
			time.Sleep(pollInterval)
			continue
		}

		// Expired or invalid
		if resp.StatusCode == http.StatusBadRequest {
			return nil, fmt.Errorf("authentication failed or expired")
		}

		time.Sleep(pollInterval)
	}

	return nil, fmt.Errorf("authentication timed out")
}

// Authenticate performs the complete authentication flow
func Authenticate() error {

	config := LoadConfig()
	serverURL := defaultServer

	fmt.Println("🔐 Starting aaaauthentication...")

	// start Auth flow
	fmt.Println("starting auth flow")
	authResp, err := startAuthFlow(serverURL)
	if err != nil {
		return err
	}

	fmt.Printf("🌐 Opening browser: %s\n", authResp.VerificationURL)

	// Open browser
	if err := openBrowser(authResp.VerificationURL); err != nil {
		fmt.Printf("❌ Failed to open browser: %v\n", err)
		fmt.Printf("Please open: %s\n", authResp.VerificationURL)
	}

	fmt.Print("⏳ Waiting for authentication")

	// Wait for token
	tokenResp, err := pollForToken(serverURL, authResp.DeviceCode)
	if err != nil {
		fmt.Println()
		return err
	}

	fmt.Println("\n✅ Authentication successful!")

	// Save token
	config.JWT = tokenResp.AccessToken
	config.Server = serverURL
	fmt.Println(config);
	

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
