package main

import (
	"fmt"
	"log"

	"github.com/amirtalbi/examen_go/internal/config"
	"github.com/amirtalbi/examen_go/internal/domain/models"
	"github.com/amirtalbi/examen_go/internal/domain/repositories"
	"github.com/amirtalbi/examen_go/internal/service"
)

func main() {
	// Initialiser la configuration
	cfg := config.Load()

	// Créer un repository utilisateur en mémoire pour les tests
	userRepo := repositories.NewUserRepository()

	// Créer un service d'authentification
	authService := service.NewAuthService(userRepo, cfg)

	// Enregistrer un utilisateur de test
	registerRequest := models.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	_, err := authService.Register(registerRequest)
	if err != nil {
		log.Fatalf("Erreur lors de l'enregistrement de l'utilisateur: %v", err)
	}

	fmt.Println("\n=== Test du processus de réinitialisation de mot de passe avec JWT ===\n")

	// 1. Demander un token de réinitialisation
	fmt.Println("1. Demande d'un token de réinitialisation pour:", registerRequest.Email)
	resetToken, err := authService.ForgotPassword(registerRequest.Email)
	if err != nil {
		log.Fatalf("Erreur lors de la demande de réinitialisation: %v", err)
	}

	fmt.Println("✅ Token de réinitialisation reçu directement dans la réponse:", resetToken)

	// 2. Utiliser le token pour réinitialiser le mot de passe
	newPassword := "newpassword456"
	fmt.Println("\n2. Réinitialisation du mot de passe avec le token JWT")
	resetRequest := models.ResetPasswordRequest{
		Token:       resetToken,
		NewPassword: newPassword,
	}

	err = authService.ResetPassword(resetRequest)
	if err != nil {
		log.Fatalf("❌ Erreur lors de la réinitialisation du mot de passe: %v", err)
	}

	fmt.Println("✅ Mot de passe réinitialisé avec succès")

	// 3. Vérifier que le nouveau mot de passe fonctionne
	fmt.Println("\n3. Tentative de connexion avec le nouveau mot de passe")
	loginRequest := models.LoginRequest{
		Email:    registerRequest.Email,
		Password: newPassword,
	}

	_, err = authService.Login(loginRequest)
	if err != nil {
		log.Fatalf("❌ Erreur lors de la connexion avec le nouveau mot de passe: %v", err)
	}

	fmt.Println("✅ Connexion réussie avec le nouveau mot de passe")
	fmt.Println("\n=== Test complet réussi ===\n")
}