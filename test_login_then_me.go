package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Premièrement, connectons-nous pour obtenir un token valide
	// Préparer les données de connexion
	loginData := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Créer la requête de connexion
	loginUrl := "http://localhost:8081/4efb0957-d14e-437f-8f01-a8db9f47405b/login"
	req, err := http.NewRequest("POST", loginUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Envoyer la requête
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Afficher le statut de la réponse
	fmt.Printf("Login Status: %s\n", resp.Status)

	// Si la connexion échoue, essayons d'abord de créer un utilisateur
	if resp.StatusCode != 200 {
		fmt.Println("Login failed, trying to register first...")
		
		// Préparer les données d'inscription
		registerData := map[string]string{
			"email":    "test@example.com",
			"password": "password123",
			"name":     "Test User",
		}
		jsonData, err = json.Marshal(registerData)
		if err != nil {
			log.Fatalf("Failed to marshal JSON: %v", err)
		}

		// Créer la requête d'inscription
		registerUrl := "http://localhost:8081/4efb0957-d14e-437f-8f01-a8db9f47405b/register"
		req, err = http.NewRequest("POST", registerUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Envoyer la requête
		resp, err = client.Do(req)
		if err != nil {
			log.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		// Afficher le statut de la réponse
		fmt.Printf("Register Status: %s\n", resp.Status)

		// Réessayer la connexion
		jsonData, _ = json.Marshal(loginData)
		req, _ = http.NewRequest("POST", loginUrl, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp, err = client.Do(req)
		if err != nil {
			log.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()
		fmt.Printf("Login Status after registration: %s\n", resp.Status)
	}

	// Décoder la réponse pour obtenir le token
	var loginResult map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&loginResult); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	// Extraire le token
	token, ok := loginResult["token"].(string)
	if !ok {
		log.Fatalf("Failed to get token from response: %+v", loginResult)
	}
	fmt.Printf("Obtained token: %s\n", token)

	// Maintenant, utilisons le token pour accéder à l'endpoint /me
	meUrl := "http://localhost:8081/4efb0957-d14e-437f-8f01-a8db9f47405b/me"
	req, err = http.NewRequest("GET", meUrl, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	
	// Définir l'en-tête d'autorisation avec le token
	req.Header.Set("Authorization", "Bearer "+token)

	// Envoyer la requête
	resp, err = client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Afficher le statut de la réponse
	fmt.Printf("ME Endpoint Status: %s\n", resp.Status)
	
	// Décoder le corps de la réponse
	var meResult map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&meResult); err != nil {
		log.Printf("Failed to decode response: %v", err)
	} else {
		fmt.Printf("ME Endpoint Response: %+v\n", meResult)
	}
}
