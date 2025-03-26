package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

func GenerateToken(userID string, secret string, expiryHours int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * time.Duration(expiryHours)).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GetTokenClaims(tokenString string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func ValidateToken(tokenString string, secret string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", errors.New("invalid claim: user_id")
		}
		return userID, nil
	}

	return "", errors.New("invalid token")
}

func GenerateRefreshToken(userID string, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
		"iat":     time.Now().Unix(),
		"type":    "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateRefreshToken(tokenString string, secret string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "refresh" {
			return "", errors.New("invalid token type")
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", errors.New("invalid claim: user_id")
		}
		return userID, nil
	}

	return "", errors.New("invalid token")
}

// GenerateResetToken génère un JWT pour la réinitialisation de mot de passe
// avec un identifiant unique (uid) pour le token
func GenerateResetToken(email string, secret string, expiryHours int) (string, string, error) {
	// Générer un identifiant unique pour ce token de réinitialisation
	tokenUID := uuid.New().String()

	claims := jwt.MapClaims{
		"email": email,
		"uid":   tokenUID,
		"exp":   time.Now().Add(time.Hour * time.Duration(expiryHours)).Unix(),
		"iat":   time.Now().Unix(),
		"type":  "reset",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	return tokenString, tokenUID, nil
}

// ValidateResetToken valide un JWT de réinitialisation de mot de passe
// et retourne l'email associé si le token est valide
func ValidateResetToken(tokenString string, secret string) (string, string, error) {
	claims, err := GetTokenClaims(tokenString, secret)
	if err != nil {
		return "", "", err
	}

	// Vérifier que c'est bien un token de type reset
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "reset" {
		return "", "", errors.New("invalid token type")
	}

	// Récupérer l'email et l'uid
	email, ok := claims["email"].(string)
	if !ok {
		return "", "", errors.New("invalid token claims: missing email")
	}

	uid, ok := claims["uid"].(string)
	if !ok {
		return "", "", errors.New("invalid token claims: missing uid")
	}

	return email, uid, nil
}
