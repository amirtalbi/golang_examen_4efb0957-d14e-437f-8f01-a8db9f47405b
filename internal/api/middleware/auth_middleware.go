package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/amirtalbi/examen_go/internal/service"
	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("\n\n==================================================")
		log.Printf("▶️ NOUVELLE REQUÊTE - %s", time.Now().Format(time.RFC3339))
		log.Printf("==================================================")

		method := c.Request.Method
		path := c.Request.URL.Path
		log.Printf("📌 MÉTHODE: %s", method)
		log.Printf("📌 CHEMIN: %s", path)

		// Journalisation des paramètres de requête
		queryParams := c.Request.URL.Query()
		if len(queryParams) > 0 {
			log.Printf("📝 PARAMÈTRES DE REQUÊTE:")
			for key, values := range queryParams {
				log.Printf("   - %s: %v", key, values)
			}
		}

		// Journalisation des en-têtes de requête
		log.Printf("📋 TOUS LES HEADERS:")
		for name, values := range c.Request.Header {
			log.Printf("   - %s: %v", name, values)
		}

		// Journalisation du corps de la requête pour toutes les méthodes
		log.Printf("📊 TYPE DE CONTENU: %s", c.ContentType())
		
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			
			// Détection du format JSON pour une meilleure lisibilité
			if c.ContentType() == "application/json" {
				// Tentative de formatage du JSON
				var prettyJSON bytes.Buffer
				if err := json.Indent(&prettyJSON, bodyBytes, "   ", "   "); err == nil {
					log.Printf("📄 CORPS DE LA REQUÊTE (JSON):\n%s", prettyJSON.String())
				} else {
					log.Printf("📄 CORPS DE LA REQUÊTE: %s", string(bodyBytes))
				}
			} else {
				log.Printf("📄 CORPS DE LA REQUÊTE: %s", string(bodyBytes))
			}
		} else {
			log.Printf("📄 CORPS DE LA REQUÊTE: <vide>")
		}

		startTime := time.Now()
		c.Next()
		duration := time.Since(startTime)

		// Journalisation de la réponse
		log.Printf("\n--------------------------------------------------")
		log.Printf("⬅️ RÉPONSE - Traitement en %v", duration)
		log.Printf("--------------------------------------------------")
		log.Printf("📊 STATUT: %d", c.Writer.Status())
		log.Printf("📊 TAILLE: %d bytes", c.Writer.Size())
		
		// Affichage des en-têtes de réponse
		log.Printf("📋 EN-TÊTES DE RÉPONSE:")
		for name, values := range c.Writer.Header() {
			log.Printf("   - %s: %v", name, values)
		}
		
		log.Printf("==================================================\n")
	}
}

func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Vérifier si l'en-tête d'autorisation existe
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("❌ Erreur d'authentification: En-tête d'autorisation manquant")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Missing Authorization header"})
			c.Abort()
			return
		}

		// Vérifier le format de l'en-tête d'autorisation
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("❌ Erreur d'authentification: Format d'en-tête invalide: %s", authHeader)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid Authorization format"})
			c.Abort()
			return
		}

		// Valider le token
		tokenString := parts[1]
		if tokenString == "" {
			log.Printf("❌ Erreur d'authentification: Token vide")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Empty token"})
			c.Abort()
			return
		}

		userID, err := authService.ValidateToken(tokenString)
		if err != nil {
			log.Printf("❌ Erreur d'authentification: Token invalide: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid token"})
			c.Abort()
			return
		}

		// Token valide, continuer
		log.Printf("✅ Authentification réussie pour l'utilisateur ID: %s", userID)
		c.Set("token", tokenString)
		c.Set("userID", userID)
		c.Next()
	}
}
