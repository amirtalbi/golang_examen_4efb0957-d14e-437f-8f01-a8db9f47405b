package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"net/http"

	"github.com/amirtalbi/examen_go/internal/domain/models"
	"github.com/amirtalbi/examen_go/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var request models.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Register(request)
	if err != nil {
		if err == service.ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request models.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Login(request)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var request models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.ForgotPassword(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	if token == "" {
		bytes := make([]byte, 8)
		rand.Read(bytes)
		token = hex.EncodeToString(bytes)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "If your email exists, you will receive a password reset link",
		"token":   token,
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var request models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	tokenInterface, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve token"})
		return
	}

	token := tokenInterface.(string)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve user ID"})
		return
	}

	// Blacklist the access token
	err := h.authService.BlacklistToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process logout"})
		return
	}

	log.Printf("🔑 DÉCONNEXION RÉUSSIE: Token invalidé pour l'utilisateur %s", userID)
	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Lire le corps de la requête raw pour le logger
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("❌ ERREUR REFRESH: Impossible de lire le corps de la requête: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Restaurer le corps pour que Gin puisse le lire à nouveau
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Logger le corps brut reçu
	log.Printf("📥 REFRESH TOKEN - Corps reçu: %s", string(bodyBytes))

	var request struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("❌ ERREUR REFRESH: Format JSON invalide: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Logger le token de rafraîchissement extrait
	log.Printf("🔑 REFRESH TOKEN - Token à vérifier: %s", request.RefreshToken)

	// Vérifier si le refresh token est valide
	response, err := h.authService.RefreshToken(request.RefreshToken)
	if err != nil {
		if err == service.ErrInvalidToken {
			log.Printf("❌ REFRESH ÉCHOUÉ: Token invalide ou expiré")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
			return
		}
		log.Printf("❌ REFRESH ÉCHOUÉ: Erreur interne: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh token"})
		return
	}

	log.Printf("✅ REFRESH RÉUSSI: Nouveau token généré pour l'utilisateur %s", response.User.ID)
	c.JSON(http.StatusOK, response)
}

// Ajouter test intégration pour l'authentification
