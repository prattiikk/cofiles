package cloud

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/prattiikk/cofiles/internal/auth"
)

type DownloadResponse struct {
	Success   bool   `json:"success"`
	URL       string `json:"url"`
	FileName  string `json:"fileName"`
	MimeType  string `json:"mimeType"`
	FileSize  int64  `json:"fileSize"`
	ExpiresIn string `json:"expiresIn"`
}

// DownloadFile downloads a file from the server using the file ID
func DownloadFile(fileID string) error {
	config := auth.LoadConfig()
	downloadMetaURL := fmt.Sprintf("%s/cloud/download/%s", config.Server, fileID)

	req, err := http.NewRequest("GET", downloadMetaURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	authHeader, err := auth.GetAuthHeader()
	if err != nil {
		return fmt.Errorf("authorization failed: %w", err)
	}
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get signed URL: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		if json.Unmarshal(respBody, &errorResp) == nil {
			if msg, ok := errorResp["error"].(string); ok {
				return fmt.Errorf("server error (%d): %s", resp.StatusCode, msg)
			}
		}
		return fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var meta DownloadResponse
	if err := json.Unmarshal(respBody, &meta); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	if !meta.Success {
		return fmt.Errorf("server returned success=false")
	}

	fileResp, err := http.Get(meta.URL)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer fileResp.Body.Close()

	if fileResp.StatusCode != http.StatusOK {
		return fmt.Errorf("file download failed with status %d", fileResp.StatusCode)
	}

	outFile, err := os.Create(meta.FileName)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, fileResp.Body); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	return nil
}