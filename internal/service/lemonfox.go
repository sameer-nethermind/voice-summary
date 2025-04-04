package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// callLemonFoxTranscription sends the uploaded file to LemonFox for transcription.
func (us *UserService) callLemonFoxTranscription(file io.Reader, filename string) (string, error) {
	//-------------------------------------------------------------------
	// Build multipart/form-data body for LemonFox
	//-------------------------------------------------------------------
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add file field
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("unable to create form file for LemonFox: %v", err)
	}

	// Copy file contents
	if _, err := io.Copy(part, file); err != nil {
		return "", fmt.Errorf("unable to copy file content: %v", err)
	}

	// Additional fields for LemonFox
	if err := writer.WriteField("language", "english"); err != nil {
		return "", fmt.Errorf("unable to add language field: %v", err)
	}
	if err := writer.WriteField("response_format", "json"); err != nil {
		return "", fmt.Errorf("unable to add response_format field: %v", err)
	}

	// Close writer
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("error closing LemonFox writer: %v", err)
	}

	//-------------------------------------------------------------------
	// Prepare and send HTTP request
	//-------------------------------------------------------------------

	req, err := http.NewRequest(
		"POST",
		"https://api.lemonfox.ai/v1/audio/transcriptions",
		&requestBody,
	)
	if err != nil {
		return "", fmt.Errorf("error creating LemonFox request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+us.config.LemonFoxAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request to LemonFox: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("lemonFox returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	//-------------------------------------------------------------------
	// Parse response
	//-------------------------------------------------------------------
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading LemonFox response: %v", err)
	}

	var lemonFoxJSON LemonFoxResponse
	if err := json.Unmarshal(respBody, &lemonFoxJSON); err != nil {
		return "", fmt.Errorf("error parsing LemonFox JSON: %v", err)
	}

	return lemonFoxJSON.Text, nil
}
