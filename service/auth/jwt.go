package auth

import (
	"log"
	"time"

	"github.com/sudarakas/edata/config"
	"github.com/sudarakas/edata/types"

	jwt "github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(user types.User) (string, error) {
	expirationTime := time.Now().Add(6 * time.Hour)
	claims := jwt.MapClaims{
		"email": user.Email,
		"id":    user.ID,
		"iat":   time.Now().Unix(),     // Issued at time
		"exp":   expirationTime.Unix(), // Expiration time
	}

	// Create new JWT token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	secret := config.Envs.JWTSECRET
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Printf("Error signing JWT: %v", err)
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	// Parse the token with the secret key
	secret := config.Envs.JWTSECRET
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})

	// Handle parse error
	if err != nil {
		log.Printf("Error parsing JWT: %v", err)
		return nil, err
	}

	// Validate claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract claims and validate expiration time
		expirationTime := claims["exp"]

		// Check if the expiration time exists and is a valid timestamp
		if expirationTime != nil {
			expTime, ok := expirationTime.(float64) // exp is a float64 (Unix timestamp)
			if !ok {
				log.Printf("Invalid expiration time format")
			} else {
				// Convert exp to time and log
				expiration := time.Unix(int64(expTime), 0)
				log.Printf("Expiration Time: %v", expiration)

				// Check if the token has expired
				if time.Now().After(expiration) {
					log.Printf("JWT has expired")
					return nil, jwt.ErrTokenExpired
				}
			}
		} else {
			log.Printf("No expiration time in JWT")
			return nil, jwt.ErrInvalidKey
		}

		// Return the claims
		return claims, nil
	}

	// Handle invalid token
	log.Printf("Invalid JWT Token: %v", err)
	return nil, jwt.ErrInvalidKey
}
