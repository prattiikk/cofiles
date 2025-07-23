package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/prattiikk/cofiles/cmd/utils/auth"
)

// DownloadFile downloads a file from the given server URL using the filename
func DownloadFile(fileName string) error {
	// Replace with your actual backend download URL
	downloadURL := fmt.Sprintf("http://localhost:3000/core/files/download/%s", fileName)

	// Create new HTTP GET request
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Get and set authorization header
	authHeader, err := auth.GetAuthHeader()
	if err != nil {
		return fmt.Errorf("authorization failed: %w", err)
	}
	req.Header.Set("Authorization", authHeader)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: server returned status %d", resp.StatusCode)
	}

	// Create the destination file in the current directory
	outFile, err := os.Create(filepath.Base(fileName))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	// Stream file content to disk
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	fmt.Printf("✅ File downloaded: %s\n", outFile.Name())
	return nil
}
