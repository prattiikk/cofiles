package cloud

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/prattiikk/cofiles/internal/auth"
)

type File struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	MimeType  string `json:"mimeType"`
	CreatedAt string `json:"createdAt"`
}

type FileListResponse struct {
	Success bool   `json:"success"`
	Files   []File `json:"files"`
}

// GetUserFiles fetches the list of files for the logged-in user
func GetUserFiles() ([]File, error) {
	authHeader, err := auth.GetAuthHeader()
	if err != nil || authHeader == "" {
		return nil, fmt.Errorf("authorization failed: %w", err)
	}

	config := auth.LoadConfig()
	if config.Server == "" {
		return nil, fmt.Errorf("server not configured. Please run 'cofiles login'")
	}

	url := config.Server + "/cloud/files"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var fileList FileListResponse
	if err := json.Unmarshal(bodyBytes, &fileList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !fileList.Success {
		return nil, fmt.Errorf("server returned success=false")
	}

	return fileList.Files, nil
}
