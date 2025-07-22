package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/prattiikk/cofiles/cmd/utils/auth"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// UploadResponse represents the server response structure
type UploadResponse struct {
	Message string `json:"message"`
	File    struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		URL       string `json:"url"`
		CreatedAt string `json:"createdAt"`
		RoomID    *string `json:"roomId"`
		UserID    string `json:"userId"`
	} `json:"file"`
}

// ErrorResponse represents error response structure
type ErrorResponse struct {
	Error string `json:"error"`
}

func UploadFile(filePath string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	// Prepare multipart body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return fmt.Errorf("could not create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("could not copy file data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("could not close writer: %w", err)
	}

	// Replace with your actual backend URL
	uploadURL := "http://localhost:3000/core/files/upload"

	// Get authorization header
	authHeader, err := auth.GetAuthHeader()
	if err != nil {
		return fmt.Errorf("authorization failed: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("upload request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %w", err)
	}

	// Check for success status codes (2xx range)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Parse successful response
		var uploadResp UploadResponse
		if err := json.Unmarshal(respBody, &uploadResp); err != nil {
			return fmt.Errorf("could not parse success response: %w", err)
		}

		// Print success message
		fmt.Printf("✅ %s\n", uploadResp.Message)
		fmt.Printf("📁 File: %s\n", uploadResp.File.Name)
		fmt.Printf("🔗 URL: %s\n", uploadResp.File.URL)
		fmt.Printf("📅 Created: %s\n", uploadResp.File.CreatedAt)
		fmt.Printf("🆔 File ID: %s\n", uploadResp.File.ID)

		return nil
	}

	// Handle error response
	var errorResp ErrorResponse
	if err := json.Unmarshal(respBody, &errorResp); err != nil {
		// If we can't parse the error response, show raw response
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, errorResp.Error)
}