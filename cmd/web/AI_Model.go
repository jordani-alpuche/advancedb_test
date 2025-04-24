package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Define a struct to match the expected JSON response structure
type InferenceResponse struct {
	Labels []string  `json:"labels"`
	Scores []float64 `json:"scores"`
}

// Function to interact with the Hugging Face Inference API
func (app *application) generateProductTag(productName, productDescription string) (string, error) {
	// Replace with your Hugging Face API token
	apiToken := "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"
	modelEndpoint := "https://api-inference.huggingface.co/models/facebook/bart-large-mnli"

	payload := map[string]interface{}{
		"inputs": fmt.Sprintf("Product Name: %s. Product Description: %s", productName, productDescription),
		"parameters": map[string][]string{
			"candidate_labels": {"Clothing", "Electronics", "Accessory", "Footwear", "Furniture", "Tool", "Supply", "Hardware", "Display", "Equipment"},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", modelEndpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("inference API returned an error: %s - %s", resp.Status, string(bodyBytes))
	}

	var response InferenceResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Labels) > 0 {
		// The first label in the list is usually the one with the highest score
		return response.Labels[0], nil
	}

	return "", fmt.Errorf("no labels found in the inference API response")
}