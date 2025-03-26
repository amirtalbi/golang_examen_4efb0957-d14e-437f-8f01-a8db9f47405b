package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/amirtalbi/examen_go/internal/service"
	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Début de l'entrée de log avec un séparateur visible
		log.Printf("\n\n==================================================")
		log.Printf("▶️ NOUVELLE REQUÊTE - %s", time.Now().Format(time.RFC3339))
		log.Printf("==================================================")

		// Informations sur la requête
		method := c.Request.Method
		path := c.Request.URL.Path
		log.Printf("📌 MÉTHODE: %s", method)
		log.Printf("📌 CHEMIN: %s", path)

		// Paramètres de la requête
		queryParams := c.Request.URL.Query()
		if len(queryParams) > 0 {
			log.Printf("📝 PARAMÈTRES:")
			for key, values := range queryParams {
				log.Printf("   - %s: %v", key, values)
			}
		}

		// Headers importants
		log.Printf("📋 HEADERS IMPORTANTS:")
		log.Printf("   - Content-Type: %s", c.GetHeader("Content-Type"))
		log.Printf("   - User-Agent: %s", c.GetHeader("User-Agent"))

		// Pour les requêtes avec un corps
		if method == "POST" || method == "PUT" || method == "PATCH" {
			log.Printf("📊 TYPE DE CONTENU: %s", c.ContentType())
		}

		// Calcul du temps de traitement
		startTime := time.Now()
		c.Next() // Exécution de la requête
		duration := time.Since(startTime)

		// Informations sur la réponse
		log.Printf("\n--------------------------------------------------")
		log.Printf("⬅️ RÉPONSE - Traitement en %v", duration)
		log.Printf("--------------------------------------------------")
		log.Printf("📊 STATUT: %d", c.Writer.Status())
		log.Printf("📊 TAILLE: %d bytes", c.Writer.Size())
		log.Printf("==================================================\n")
	}
}

func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("\n🔒 AUTHENTIFICATION - %s %s", c.Request.Method, c.Request.URL.Path)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("❌ ERREUR AUTH: Header d'autorisation manquant")
			log.Printf("❌ DÉTAILS: %s %s", c.Request.Method, c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("❌ ERREUR AUTH: Format d'autorisation invalide")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		userID, err := authService.ValidateToken(tokenString)
		if err != nil {
			log.Printf("❌ ERREUR AUTH: Token invalide - %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		log.Printf("✅ AUTH RÉUSSIE: Utilisateur %s authentifié", userID)
		c.Set("token", tokenString)
		c.Set("userID", userID)
		c.Next()
	}
}
