package cloud

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prattiikk/cofiles/internal/auth"
)

// DeleteFile deletes a file from the server using the file ID
func DeleteFile(fileID string) error {
	config := auth.LoadConfig()
	deleteURL := fmt.Sprintf("%s/cloud/delete/%s", config.Server, fileID)

	req, err := http.NewRequest(http.MethodDelete, deleteURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %w", err)
	}

	// Attach Authorization header
	authHeader, err := auth.GetAuthHeader()
	if err != nil {
		return fmt.Errorf("authorization failed: %w", err)
	}
	req.Header.Set("Authorization", authHeader)

	// Optional: Set timeout for HTTP client
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status %d: failed to delete file", resp.StatusCode)
	}

	return nil
}
