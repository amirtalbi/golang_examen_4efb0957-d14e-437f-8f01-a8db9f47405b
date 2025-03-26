package handlers

import (
	"bytes"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		return
	}

	if request.Email == "" || request.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		return
	}

	response, err := h.authService.Login(request)
	if err != nil {
		log.Printf("Login error: %v", err)
		errorMsg := err.Error()
		if errorMsg == "user not found" || errorMsg == "password mismatch" {
			log.Printf("Returning 401 for error: %s", errorMsg)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			log.Printf("Returning 400 for error: %s", errorMsg)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		}
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

	// Appel au service pour générer un token JWT
	resetToken, err := h.authService.ForgotPassword(request.Email)
	if err != nil {
		// Pour des raisons de sécurité, nous ne révélons pas si l'email existe ou non
		// Nous retournons un message générique même en cas d'erreur
		c.JSON(http.StatusOK, gin.H{
			"message": "If your email exists, you will receive a password reset token",
		})
		return
	}

	// Retourner le token JWT dans la réponse
	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset token generated successfully",
		"token":   resetToken,
	})
}



func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var request models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Processing reset password request with token: %s", request.Token)
	
	// Vérifier si le token est valide
	err := h.authService.ResetPassword(request)
	if err != nil {
		if err == service.ErrInvalidToken {
			// Token invalide ou expiré - renvoyer 401 Unauthorized
			log.Printf("❌ Tentative de réinitialisation avec un token invalide: %s", request.Token)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		} else if err == service.ErrUserNotFound {
			log.Printf("❌ Utilisateur non trouvé pour le token: %s", request.Token)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		} else {
			log.Printf("❌ Erreur lors de la réinitialisation du mot de passe: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
			return
		}
	}

	// Succès - mot de passe réinitialisé
	log.Printf("✅ Mot de passe réinitialisé avec succès pour le token: %s", request.Token)
	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Récupérer le token depuis le contexte
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve token"})
		return
	}

	// Récupérer l'ID de l'utilisateur depuis le contexte
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve user ID"})
		return
	}

	// Récupérer le refresh token du corps de la requête
	var request struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("❌ ERREUR LOGOUT: Format JSON invalide: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Révoquer le token d'accès
	err := h.authService.RevokeToken(token.(string))
	if err != nil {
		log.Printf("❌ Erreur lors de la révocation du token d'accès: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke access token"})
		return
	}

	// Révoquer également le refresh token
	err = h.authService.RevokeToken(request.RefreshToken)
	if err != nil {
		log.Printf("❌ Erreur lors de la révocation du refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke refresh token"})
		return
	}

	log.Printf("✅ DÉCONNEXION RÉUSSIE: Token révoqué pour l'utilisateur %s", userID)
	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("ERREUR REFRESH: Impossible de lire le corps de la requête: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	log.Printf("REFRESH TOKEN - Corps reçu: %s", string(bodyBytes))

	var request struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("ERREUR REFRESH: Format JSON invalide: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("REFRESH TOKEN - Token à vérifier: %s", request.RefreshToken)

	// Ignorer le token d'accès dans l'en-tête Authorization et utiliser uniquement le refresh token
	response, err := h.authService.RefreshToken(request.RefreshToken)
	if err != nil {
		if err == service.ErrInvalidToken {
			log.Printf("REFRESH ÉCHOUÉ: Token invalide ou expiré")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
			return
		}
		log.Printf("REFRESH ÉCHOUÉ: Erreur interne: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh token"})
		return
	}

	log.Printf("REFRESH RÉUSSI: Nouveau token généré pour l'utilisateur %s", response.User.ID)
	c.JSON(http.StatusOK, response)
}
