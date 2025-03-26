package main

import (
	"fmt"
	"log"
	"os"

	"github.com/amirtalbi/examen_go/internal/config"
	"github.com/amirtalbi/examen_go/pkg/auth"
	"github.com/joho/godotenv"
)

func main() {
	// Charger les variables d'environnement
	err := godotenv.Load()
	if err != nil {
		log.Printf("Erreur lors du chargement du fichier .env: %v", err)
	}

	// Initialiser la configuration
	cfg := config.Load()

	// Test de génération d'un reset token
	email := "test@example.com"
	fmt.Println("Génération d'un reset token pour:", email)
	
	resetToken, tokenUID, err := auth.GenerateResetToken(email, cfg.ResetTokenSecret, cfg.TokenExpiryHours)
	if err != nil {
		fmt.Printf("Erreur lors de la génération du reset token: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Reset Token généré: %s\n", resetToken)
	fmt.Printf("UID du token: %s\n", tokenUID)
	
	// Test de validation du reset token
	fmt.Println("\nValidation du reset token...")
	
	validatedEmail, validatedUID, err := auth.ValidateResetToken(resetToken, cfg.ResetTokenSecret)
	if err != nil {
		fmt.Printf("Erreur lors de la validation du reset token: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Email validé: %s\n", validatedEmail)
	fmt.Printf("UID validé: %s\n", validatedUID)
	
	// Vérification des valeurs
	if validatedEmail != email {
		fmt.Printf("ERREUR: L'email validé (%s) ne correspond pas à l'email original (%s)\n", validatedEmail, email)
	} else {
		fmt.Println("Succès: L'email a été correctement validé")
	}
	
	if validatedUID != tokenUID {
		fmt.Printf("ERREUR: L'UID validé (%s) ne correspond pas à l'UID original (%s)\n", validatedUID, tokenUID)
	} else {
		fmt.Println("Succès: L'UID a été correctement validé")
	}
	
	// Test avec un mauvais secret
	fmt.Println("\nTest avec un mauvais secret...")
	_, _, err = auth.ValidateResetToken(resetToken, "mauvais-secret")
	if err != nil {
		fmt.Printf("Erreur attendue avec un mauvais secret: %v\n", err)
	} else {
		fmt.Println("ERREUR: La validation aurait dû échouer avec un mauvais secret")
	}
	
	fmt.Println("\nTest terminé avec succès!")
}
