package token

import (
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
	"time"
)

type PasetoMaker struct {
	paseto *paseto.V2
	key    []byte
}

func NewPasetoMaker(key []byte) (Maker, error) {
	if len(key) != chacha20poly1305.KeySize {
		return nil, ErrSecretLen
	}
	return &PasetoMaker{
		paseto: paseto.NewV2(),
		key:    key,
	}, nil
}

func (p *PasetoMaker) CreateToken(username string, expireDate time.Duration) (string, error) {
	paload, err := NewPayload(username, expireDate)
	if err != nil {
		return "", nil
	}
	return p.paseto.Encrypt(p.key, paload, nil)
}

func (p *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	paload := &Payload{}
	err := p.paseto.Decrypt(token, p.key, paload, nil)
	if err != nil {
		return nil, err
	}
	return paload, nil
}
