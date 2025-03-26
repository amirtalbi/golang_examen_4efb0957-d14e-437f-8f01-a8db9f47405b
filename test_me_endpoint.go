package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Get the token from the login response
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDMwNzkzNzYsImlhdCI6MTc0Mjk5Mjk3NiwidXNlcl9pZCI6IjVhMDZjMjVjLTU5YjQtNDNkNi05N2M1LWExMzBlNWEyNWQzNSJ9.AUUMZw5mYcXGpicmh0AvoQ912syrDxhXt-ZaXfmxE1g"

	// Create request
	url := "http://localhost:8080/4efb0957-d14e-437f-8f01-a8db9f47405b/me"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	
	// Set the Authorization header with the token
	req.Header.Set("Authorization", "Bearer "+token)

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
