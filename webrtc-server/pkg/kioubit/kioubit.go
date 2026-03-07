package kioubit

import (
	"crypto/ecdsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

const defaultTimeoutSecs = 60

func VerifyAuthToken(signature, params string, pemPubKey []byte, myDomain string) (*KioubitAuthCallbackParams, error) {
	userData := new(KioubitAuthCallbackParams)
	var err error

	// Read public key
	blockPub, _ := pem.Decode(pemPubKey)
	if blockPub == nil {
		return nil, errors.New("failed to decode PEM block")
	}
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	if err != nil {
		return nil, errors.New("internal server error")
	}
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	// Hash parameters
	hash := sha512.Sum512([]byte(params))

	// Decode base64 signature
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return nil, errors.New("failed to decode signature")
	}

	// Verify signature
	if !ecdsa.VerifyASN1(publicKey, hash[:], signatureBytes) {
		return nil, ErrInvalidSignature
	}

	// Decode parameters
	parameterBytes, err := base64.StdEncoding.DecodeString(params)
	if err != nil {
		return nil, fmt.Errorf("failed decoding verified parameters: %w", err)
	}

	err = json.Unmarshal(parameterBytes, &userData)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshaling verified parameters: %w", err)
	}

	if math.Abs(userData.Time-float64(time.Now().Unix())) > defaultTimeoutSecs {
		return nil, ErrRequestExpired
	}

	if domain := userData.Domain; domain == nil || *domain != myDomain {
		log.New(os.Stderr, "", 0).Printf("Domain mismatch from kioubit's callback param: ours=%q, theirs=%q", myDomain, *domain)
		// return nil, ErrDomainMismatch
	}

	return userData, nil
}

func (params *KioubitAuthCallbackParams) GetNonce() string {
	if params != nil {
		if token := params.UserToken; token != nil {
			return *token
		}
	}
	return ""
}

func NewKioubitAuthCallbackParamsFromHTTPRequest(r *http.Request, kioubitPubkey []byte) (*KioubitAuthCallbackParams, error) {
	if kioubitPubkey == nil {
		return nil, ErrInvalidPubkey
	}

	kioubitAuthZCBParams := r.URL.Query().Get("params")
	if kioubitAuthZCBParams == "" {
		return nil, ErrNoParamsFound
	}

	kioubitAuthZCBSignature := r.URL.Query().Get("signature")
	if kioubitAuthZCBSignature == "" {
		return nil, ErrNoSignatureFound
	}

	return VerifyAuthToken(kioubitAuthZCBSignature, kioubitAuthZCBParams, kioubitPubkey, r.URL.Hostname())
}
