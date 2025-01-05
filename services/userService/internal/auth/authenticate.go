package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"
	"younified-backend/contracts/user/model"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrInvalidToken is returned when the token is invalid
	ErrInvalidToken = errors.New("invalid token")

	// ErrExpiredToken is returned when the token has expired
	ErrExpiredToken = errors.New("token has expired")

	// ErrInvalidCredentials is returned when login credentials are incorrect
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// TokenClaim struct for JWT claims
type TokenClaim struct {
	Username string             `json:"username"`
	UserID   primitive.ObjectID `json:"user_id"`
	UnionID  primitive.ObjectID `json:"union_id"`
	jwt.StandardClaims
}

// HashPassword creates a secure password hash using bcrypt with additional salt from UnionID
func HashPassword(password string, unionID string) (string, error) {
	// Create a salt by combining password with UnionID
	salt := createSalt(password, unionID)

	// Use bcrypt to hash the salted password
	// The work factor (cost) is typically between 12-14
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(salt), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// VerifyPassword checks if the provided password matches the stored hash
func VerifyPassword(hashedPassword string, inputPassword string, unionID string) bool {
	// Create the same salt used during hashing
	salt := createSalt(inputPassword, unionID)

	// Compare the input (salted) with the stored hash
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(salt))
	return err == nil
}

// createSalt generates a unique salt by combining password and unionID
func createSalt(password string, unionID string) string {
	// Create a SHA-256 hash of the combined password and unionID
	hash := sha256.New()
	hash.Write([]byte(password + unionID))
	return hex.EncodeToString(hash.Sum(nil))
}

// GenerateJWTToken creates a new JWT token
func GenerateJWTToken(username string, userID, unionID model.ObjectID, expirationHours int) (res string, err error) {
	// Get JWT secret from environment variable
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		err = fmt.Errorf("could not create session please ask admin to check")
		return
	}

	// Set token expiration
	expirationTime := time.Now().Add(time.Duration(expirationHours) * time.Hour)

	// Create the JWT claims
	claims := TokenClaim{
		Username: username,
		UserID:   userID,
		UnionID:  unionID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "User-Service",
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token
	res, _ = token.SignedString([]byte(jwtSecret))
	return
}

// ValidateJWTToken validates and parses a JWT token
func ValidateJWTToken(tokenString string) (*TokenClaim, error) {
	// Get JWT secret from environment variable
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaim{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		return []byte(jwtSecret), nil
	})

	// Check for parsing errors
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrExpiredToken
			}
		}
		return nil, ErrInvalidToken
	}

	// Extract and type assert claims
	if claims, ok := token.Claims.(*TokenClaim); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// RefreshJWTToken creates a new token with extended expiration
func RefreshJWTToken(tokenString string, expirationHours int) (string, error) {
	// Validate the existing token
	claims, err := ValidateJWTToken(tokenString)
	if err != nil {
		return "", err
	}

	// Generate a new token with the same user information
	return GenerateJWTToken(
		claims.Username,
		claims.UserID,
		claims.UnionID,
		expirationHours,
	)
}

// GeneratePasswordResetToken creates a short-lived token for password reset
func GeneratePasswordResetToken(username string, userID, unionID model.ObjectID) (string, error) {
	// Shorter expiration for password reset tokens (e.g., 1 hour)
	return GenerateJWTToken(username, userID, unionID, 1)
}

// IsPasswordCompromised performs basic password strength checks
func IsPasswordCompromised(password string) bool {
	// Basic password strength checks
	// You can expand these rules as needed
	return len(password) < 8 ||
		len(password) > 64 ||
		!containsUppercase(password) ||
		!containsLowercase(password) ||
		!containsDigit(password)
}

// Helper functions for password strength
func containsUppercase(s string) bool {
	for _, char := range s {
		if char >= 'A' && char <= 'Z' {
			return true
		}
	}
	return false
}

func containsLowercase(s string) bool {
	for _, char := range s {
		if char >= 'a' && char <= 'z' {
			return true
		}
	}
	return false
}

func containsDigit(s string) bool {
	for _, char := range s {
		if char >= '0' && char <= '9' {
			return true
		}
	}
	return false
}
