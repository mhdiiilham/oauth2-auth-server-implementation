package jwt

import (
	"fmt"
	"time"

	"github.com/mhdiiilham/oauth2-auth-server-implementation/entity/user"

	"github.com/dgrijalva/jwt-go"
)

// TokenService interface
type TokenService interface {
	Generate(entity *user.User) string
	Verify(jwtToken string) (*jwt.Token, error)
	Validate(jwtToken string) error
	GetIssuer() string
	Extract(jwtToken string) (*jwtCustomCliams, error)
}

type jwtCustomCliams struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type jwtService struct {
	Secret string
	Issuer string
}

// NewJWTService Create New JWT Service
func NewJWTService(secret, appName string) *jwtService {
	return &jwtService{
		Secret: secret,
		Issuer: appName,
	}
}

// GetIssuer of a
// TokenService
func (j *jwtService) GetIssuer() string {
	return j.Issuer
}

// Generate JWT Token
func (j *jwtService) Generate(e *user.User) string {
	claims := &jwtCustomCliams{
		e.Email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 42).Unix(),
			Issuer:    j.Issuer,
			IssuedAt:  time.Now().Unix(),
			Subject:   fmt.Sprintf("%d", e.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, _ := token.SignedString([]byte(j.Secret))
	return t
}

// Verify if Token
// is correct
func (j *jwtService) Verify(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

// Validate to see
// if token expired or not
func (j *jwtService) Validate(tokenString string) error {
	token, err := j.Verify(tokenString)

	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}

	return nil
}

// Extract JWTToken MetaData
func (j *jwtService) Extract(tokenString string) (*jwtCustomCliams, error) {
	token, err := j.Verify(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		email, ok := claims["email"].(string)
		if !ok {
			return nil, err
		}

		return &jwtCustomCliams{
			Email: email,
		}, nil
	}

	return nil, err
}
