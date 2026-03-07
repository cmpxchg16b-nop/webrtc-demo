package kioubit

import (
	"errors"
)

var ErrRequestExpired = errors.New("the request has expired")
var ErrDomainMismatch = errors.New("this request is for a different domain")
var ErrInvalidPubkey = errors.New("invalid pubkey")
var ErrNoParamsFound = errors.New("no params found")
var ErrNoSignatureFound = errors.New("no signature found")
var ErrInvalidSignature = errors.New("invalid signature")

type KioubitAuthCallbackParams struct {
	ASN          string   `json:"asn"`
	Time         float64  `json:"time"`
	Allowed4     []string `json:"allowed4,omitempty"`
	Allowed6     []string `json:"allowed6,omitempty"`
	Mnt          []string `json:"mnt,omitempty"`
	EffectiveMnt string   `json:"effective_mnt"`
	Domain       *string  `json:"domain,omitempty"`
	AuthType     *string  `json:"authtype,omitempty"`
	UserToken    *string  `json:"user_token,omitempty"`
}
