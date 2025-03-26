package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Token invalide
	invalidToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDMwNzk1MDEsImlhdCI6MTc0Mjk5MzEwMSwidXNlcl9pZCI6IjM5M2YwYTBhLTQ4NjQtNDE0MC04NmZhLTVkY2E0ZjQyZmZjYiJ9.INVALID_SIGNATURE"

	// Create request
	url := "http://68.183.71.248/4efb0957-d14e-437f-8f01-a8db9f47405b/me"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	
	// Set the Authorization header with the invalid token
	req.Header.Set("Authorization", "Bearer "+invalidToken)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Print response
	fmt.Printf("Status: %s\n", resp.Status)
	
	// Decode response body
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode response: %v", err)
	} else {
		fmt.Printf("Response: %+v\n", result)
	}
}
