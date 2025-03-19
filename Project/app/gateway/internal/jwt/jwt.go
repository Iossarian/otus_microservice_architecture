package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type Payload struct {
	ID   string
	Name string
	jwt.RegisteredClaims
}

type Service struct {
	secret string
}

func NewService(secret string) *Service {
	return &Service{secret: secret}
}

func (s *Service) Parse(tokenString string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "parse jwt token")
	}

	if claims, ok := token.Claims.(*Payload); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
