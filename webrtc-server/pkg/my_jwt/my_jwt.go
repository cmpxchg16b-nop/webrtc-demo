package my_jwt

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
)

type JWTManager interface {
	Issue(ctx context.Context, userid string) (string, error)
	Validate(ctx context.Context, token string) (bool, error)
}

const defaultSecretLengthBytes int = 32

type SimpleJWTManager struct {
	secret []byte
}

func NewSimpleJWTManager(secret []byte) {
	secMng := &SimpleJWTManager{secret: secret}
	if secret == nil {
		secMng.secret = make([]byte, defaultSecretLengthBytes)
		if _, err := rand.Read(secMng.secret); err != nil {
			log.Fatal(errors.New("failed to initialize simple jwt manager"))
		}
	}
	return secMng
}
