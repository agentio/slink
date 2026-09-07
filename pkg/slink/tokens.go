package slink

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/agentio/slink/pkg/resolve"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jws"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/mr-tron/base58"
)

func NewJwtID() string {
	bytes := make([]byte, 16)
	_, _ = rand.Read(bytes) // never errors (https://cs.opensource.google/go/go/+/refs/tags/go1.27.1:src/crypto/rand/rand.go;l=47)
	return hex.EncodeToString(bytes)
}

func GenerateAuthToken(keybytes []byte, claims map[string]any, typ string) ([]byte, error) {
	privateJwk, err := jwk.Import(secp256k1.PrivKeyFromBytes(keybytes).ToECDSA())
	if err != nil {
		return nil, err
	}
	publicJwk, err := privateJwk.PublicKey()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	builder := jwt.NewBuilder().
		JwtID(NewJwtID()).
		IssuedAt(now).
		Expiration(now.Add(30 * time.Second))
	for k, v := range claims {
		builder = builder.Claim(k, v)
	}
	token, err := builder.Build()
	if err != nil {
		return nil, err
	}
	var headers jws.Headers = jws.NewHeaders()
	if err = headers.Set("typ", typ); err != nil {
		return nil, err
	}
	if err = headers.Set("alg", "ES256"); err != nil {
		return nil, err
	}
	if err = headers.Set("jwk", publicJwk); err != nil {
		return nil, err
	}
	return jwt.Sign(token, jwt.WithKey(jwa.ES256(), privateJwk, jws.WithProtectedHeaders(headers)))
}

func VerifyAuthHeader(ctx context.Context, authHeader string) (*jwt.Token, error) {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, fmt.Errorf("invalid authorization header (expected \"Bearer\" prefix)")
	}
	return VerifyAuthToken(ctx, strings.TrimPrefix(authHeader, "Bearer "))
}

func VerifyAuthToken(ctx context.Context, accessToken string) (*jwt.Token, error) {
	token, err := jwt.ParseInsecure([]byte(accessToken))
	if err != nil {
		return nil, err
	}
	did, ok := token.Issuer()
	if !ok {
		return nil, fmt.Errorf("issuer field not set (it should be the caller's DID)")
	}
	didDoc, err := resolve.Did(ctx, did)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup user DID: %v", err)
	}
	keyTypes := []string{}
	for _, m := range didDoc.VerificationMethod {
		if m.Type == "Multikey" {
			key, err := parseMultibasePublicKey(m.PublicKeyMultibase)
			if err != nil {
				return nil, fmt.Errorf("failed to get user public key: %v", err)
			}
			token, err = jwt.Parse([]byte(accessToken), jwt.WithKey(jwa.ES256K(), key.ToECDSA()))
			if err != nil {
				return nil, err
			}
			return &token, nil
		} else {
			keyTypes = append(keyTypes, m.Type)
		}
	}
	return nil, fmt.Errorf("none of these key types are supported: %s", strings.Join(keyTypes, ", "))
}

func parseMultibasePublicKey(encoding string) (*secp256k1.PublicKey, error) {
	if len(encoding) < 2 || encoding[0] != 'z' {
		return nil, fmt.Errorf("not a multibase base58btc string")
	}
	data, err := base58.Decode(encoding[1:])
	if err != nil {
		return nil, fmt.Errorf("not a multibase base58btc string")
	}
	if len(data) < 3 {
		return nil, fmt.Errorf("multibase key was too short")
	}
	if data[0] == 0xE7 && data[1] == 0x01 {
		// multicodec secp256k1-pub, code 0xE7, varint bytes: [0xE7, 0x01]
		return secp256k1.ParsePubKey(data[2:])
	} else {
		return nil, fmt.Errorf("unsupported key type (unknown multicodec prefix)")
	}
}
