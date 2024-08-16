package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const MIN_SECRET_KEY_SIZE = 6

// Using token by using JWT maker
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker secret key
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < MIN_SECRET_KEY_SIZE {
		return nil, fmt.Errorf("Invalid key size: must be equal or greater than %d", MIN_SECRET_KEY_SIZE)
	}

	return &JWTMaker{
		secretKey: secretKey,
	}, nil
}

func (jwtMaker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(jwtMaker.secretKey))

}

func (jwtMaker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		// keyFunc will receive the parsed token and should return the key for validating.
		_, ok := t.Method.(*jwt.SigningMethodHMAC) // ensure the convert type is signing method with HMAC
		if !ok {
			return nil, ErrInvalidToken
		}

		// convert into byte string
		return []byte(jwtMaker.secretKey), nil
	}

	// Parse jwt token using claims
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		// since the error is hidden from the implementation
		verf, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verf.Inner, ErrInvalidToken) {
			return nil, ErrInvalidToken
		}
		return nil, ErrExpired
	}

	// get data from jwtToken Claims
	// Assert that jwtToken claims should be types as Payload.
	// Since above parse as holding the struct type cast as Payload struct

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
