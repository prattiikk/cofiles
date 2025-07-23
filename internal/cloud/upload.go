package cloud

import (
	"bytes"
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
	Success bool   `json:"success"`
	Message string `json:"message"`
	File    struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Size       int64  `json:"size"`
		MimeType   string `json:"mimeType"`
		UploadedAt string `json:"uploadedAt"`
	} `json:"file"`
}

// ErrorResponse represents error response structure
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func UploadFile(filePath string) error {
	// Check if file exists and get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("could not stat file: %w", err)
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("path is a directory, not a file")
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	// Determine MIME type
	ext := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Prepare multipart body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileName := filepath.Base(filePath)

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return fmt.Errorf("could not create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("could not copy file data: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("could not close writer: %w", err)
	}

	// Load auth and prepare request
	authHeader, err := auth.GetAuthHeader()
	if err != nil {
		return fmt.Errorf("authorization failed: %w", err)
	}
	uploadURL := "http://localhost:3000/cloud/upload"
	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("upload request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read and handle response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %w", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var uploadResp UploadResponse
		if err := json.Unmarshal(respBody, &uploadResp); err != nil {
			return fmt.Errorf("could not parse success response: %w", err)
		}
		return nil
	}

	// Handle error response
	var errorResp ErrorResponse
	if err := json.Unmarshal(respBody, &errorResp); err != nil {
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(respBody))
	}
	return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, errorResp.Error)
}
