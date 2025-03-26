package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// First, request a password reset token
	forgotPasswordData := map[string]string{
		"email": "test@example.com",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(forgotPasswordData)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Create request for forgot-password
	forgotUrl := "http://localhost:8080/4efb0957-d14e-437f-8f01-a8db9f47405b/forgot-password"
	req, err := http.NewRequest("POST", forgotUrl, bytes.NewBuffer(jsonData))
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
	fmt.Printf("Forgot Password Status: %s\n", resp.Status)

	// Now test the reset-password endpoint with a simulated token
	resetPasswordData := map[string]string{
		"token":        "test-token",
		"new_password": "newpassword123",
	}

	// Convert to JSON
	jsonData, err = json.Marshal(resetPasswordData)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Create request for reset-password
	resetUrl := "http://localhost:8080/4efb0957-d14e-437f-8f01-a8db9f47405b/reset-password"
	req, err = http.NewRequest("POST", resetUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err = client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Print response
	fmt.Printf("Reset Password Status: %s\n", resp.Status)
	
	// Decode response body if there is any
	if resp.ContentLength > 0 {
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Printf("Failed to decode response: %v", err)
		} else {
			fmt.Printf("Response: %+v\n", result)
		}
	}
}
