package handlers

import (
	"log"
	"net/http"

	"github.com/amirtalbi/examen_go/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	// Vérifier si l'ID de l'utilisateur est présent dans le contexte
	userID, exists := c.Get("userID")
	if !exists {
		// Cela ne devrait jamais arriver si le middleware d'authentification fonctionne correctement
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User ID not found in context"})
		return
	}

	// Vérifier si l'ID de l'utilisateur est une chaîne valide
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Récupérer les informations de l'utilisateur
	user, err := h.userService.GetUserByID(userIDStr)
	if err != nil {
		// Journaliser l'erreur
		log.Printf("❌ Erreur lors de la récupération de l'utilisateur ID %s: %v", userIDStr, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Renvoyer les informations de l'utilisateur
	log.Printf("✅ Profil récupéré avec succès pour l'utilisateur ID: %s", userIDStr)
	c.JSON(http.StatusOK, user)
}
