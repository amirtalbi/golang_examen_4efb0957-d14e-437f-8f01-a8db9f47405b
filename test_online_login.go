package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Login with the new user credentials
	loginData := map[string]string{
		"email":    "online@example.com",
		"password": "password123",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Create request
	url := "http://68.183.71.248/4efb0957-d14e-437f-8f01-a8db9f47405b/login"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

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
		
		// Save the token for later use
		if token, ok := result["token"].(string); ok {
			fmt.Printf("\nToken for use in other tests: %s\n", token)
		}
	}
}
