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
		"email":    "test@example.com",
		"password": "password123",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Create request
	url := "http://localhost:8080/4efb0957-d14e-437f-8f01-a8db9f47405b/login"
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
	}
}
