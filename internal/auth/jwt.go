package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	secretKey string
	Aud       string
	Iss       string
	Exp       int64
}

func validateJWTConfig(secretKey, aud, iss string, exp int64) {

	if secretKey == "" {
		panic("secret key cannot be empty")
	}

	if aud == "" {
		panic("audience cannot be empty")
	}

	if iss == "" {
		panic("issuer cannot be empty")
	}

	if exp <= 0 {
		panic("expiration time must be greater than zero")
	}
}
func NewJWTAuthenticator(secretKey, aud, iss string, exp int64) *JWTAuthenticator {

	validateJWTConfig(secretKey, aud, iss, exp)

	return &JWTAuthenticator{secretKey: secretKey, Aud: aud, Iss: iss, Exp: exp}
}

func (j *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secretKey))
}

// ValidateToken validates the given JWT token string and returns the parsed token if valid.
func (j *JWTAuthenticator) ValidateToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(j.secretKey), nil
	},
		jwt.WithAudience(j.Aud),
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}

// GetUserIDFromClaims extracts the user ID from the "sub" claim in the JWT token claims.
func (j *JWTAuthenticator) GetUserIDFromClaims(token *jwt.Token) (int64, error) {

	// extract user ID from token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	sub, ok := claims["sub"].(float64) // JWT numeric claims are float64
	if !ok {
		return 0, errors.New("invalid token claims: missing 'sub'")
	}

	userID := int64(sub)

	return userID, nil
}
