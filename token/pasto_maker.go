package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PastoMaker struct {
	pasto        *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d character", chacha20poly1305.KeySize)
	}

	// since two method as declare with pointer receiver struct -> should be return as address
	maker := &PastoMaker{
		pasto:        paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (p *PastoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	// create new payload
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	// assign payload into encrypt
	return p.pasto.Encrypt(p.symmetricKey, payload, nil)
}

func (p *PastoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	// decrypt token by using paseto
	err := p.pasto.Decrypt(token, p.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil

}
