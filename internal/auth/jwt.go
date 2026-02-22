package auth

import "github.com/golang-jwt/jwt/v5"

type JWTAuthenticator struct {
	secretKey string
	Aud       string
	Iss       string
	Exp       int64
}

func validateJWTConfig(secretKey, aud, iss string) {

	if secretKey == "" {
		panic("secret key cannot be empty")
	}

	if aud == "" {
		panic("audience cannot be empty")
	}

	if iss == "" {
		panic("issuer cannot be empty")
	}
}
func NewJWTAuthenticator(secretKey, aud, iss string) *JWTAuthenticator {

	validateJWTConfig(secretKey, aud, iss)

	if aud == "" {
		panic("audience cannot be empty")
	}

	if iss == "" {
		panic("issuer cannot be empty")
	}
	return &JWTAuthenticator{secretKey: secretKey, Aud: aud, Iss: iss}
}

func (j *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTAuthenticator) ValidateToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(j.secretKey), nil
	})
}
