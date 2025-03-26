package middleware

import (
	"bytes"
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
			
			// Lire et logger le corps de la requête
			var bodyBytes []byte
			if c.Request.Body != nil {
				bodyBytes, _ = io.ReadAll(c.Request.Body)
				// Restaurer le corps pour que les handlers puissent le lire
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				// Logger le corps de la requête
				log.Printf("📄 CORPS DE LA REQUÊTE: %s", string(bodyBytes))
			}
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
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		userID, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set("token", tokenString)
		c.Set("userID", userID)
		c.Next()
	}
}
