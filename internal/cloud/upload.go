package cloud

import (
	"encoding/json"
	"fmt"
	"github.com/prattiikk/cofiles/internal/auth"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// UploadResponse represents the server response structure
type UploadResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	MimeType   string `json:"mimeType"`
	CreatedAt  string `json:"createdAt"`
}

// ErrorResponse represents error response structure
type ErrorResponse struct {
	Error string `json:"error"`
}

func UploadFile(filePath string) error {

	// Validate file
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("could not stat file: %w", err)
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("path is a directory, not a file")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	// Detect MIME type
	ext := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Streaming pipe
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	// Start writing multipart body in separate goroutine
	go func() {
		defer pw.Close()
		defer writer.Close()

		part, err := writer.CreateFormFile("file", filepath.Base(filePath))
		if err != nil {
			pw.CloseWithError(err)
			return
		}

		if _, err := io.Copy(part, file); err != nil {
			pw.CloseWithError(err)
			return
		}
	}()

	// Load auth header
	authHeader, err := auth.GetAuthHeader()
	if err != nil || authHeader == "" {
		return fmt.Errorf("authorization failed: %w", err)
	}

	// Prepare request
	uploadURL := "http://ec2-43-205-235-230.ap-south-1.compute.amazonaws.com/files/upload"
	req, err := http.NewRequest("POST", uploadURL, pr)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("upload request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response body: %w", err)
	}

	// Handle success
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var uploadResp UploadResponse
		if err := json.Unmarshal(respBody, &uploadResp); err != nil {
			return fmt.Errorf("could not parse success response: %w", err)
		}
		return nil
	}

	// Handle error
	var errorResp ErrorResponse
	if err := json.Unmarshal(respBody, &errorResp); err == nil && errorResp.Error != "" {
		return fmt.Errorf("upload failed (%d): %s", resp.StatusCode, errorResp.Error)
	}

	return fmt.Errorf("upload failed (%d): %s", resp.StatusCode, string(respBody))
}