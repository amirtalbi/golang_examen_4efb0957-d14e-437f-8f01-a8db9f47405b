package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Create a new user with simple credentials
	registerData := map[string]string{
		"name":     "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(registerData)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Create request
	url := "http://localhost:8080/567088a9-6689-4e67-b5e5-ed40ad0a830c/register"
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
