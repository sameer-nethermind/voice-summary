package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	types "github.com/cyberhawk12121/Saarthi/internal/shared"
)

// callLlamaAPI uses the transcribed text as input to the Llama API.
func (us *UserService) callLlamaAPI(transcribedText string) ([]byte, int, error) {
	// Build the request payload
	payload := types.LlamaRequest{
		Messages: []map[string]string{
			{
				"role":    "user",
				"content": transcribedText,
			},
		},
		Functions: []map[string]interface{}{
			{
				"name":        "get_current_weather",
				"description": "Get the current weather in a given location",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"location": map[string]interface{}{
							"type":        "string",
							"description": "The city and state, e.g. San Francisco, CA",
						},
						"days": map[string]interface{}{
							"type":        "number",
							"description": "for how many days ahead you want the forecast",
						},
						"unit": map[string]interface{}{
							"type": "string",
							"enum": []string{"celsius", "fahrenheit"},
						},
					},
					"required": []string{"location", "days"},
				},
			},
		},
		Stream:       false,
		FunctionCall: "get_current_weather",
	}

	reqBodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, 0, fmt.Errorf("error marshaling llama request: %v", err)
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.llama-api.com/chat/completions",
		bytes.NewBuffer(reqBodyBytes),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("error creating llama request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+us.config.LlamaAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("error sending request to Llama API: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("error reading Llama API response: %v", err)
	}

	// Debug: check the response
	fmt.Printf("Llama API response: %s\n", respBody)

	return respBody, resp.StatusCode, nil
}
