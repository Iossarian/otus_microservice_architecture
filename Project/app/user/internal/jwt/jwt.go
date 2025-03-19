package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type CustomClaims struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	jwt.RegisteredClaims
}

type Service struct {
	secret string
}

func NewService(secret string) *Service {
	return &Service{secret: secret}
}

func (s *Service) Token(id, name string) (string, error) {
	claims := &CustomClaims{
		id,
		name,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", errors.Wrap(err, "failed to sign token")
	}

	return t, nil
}
