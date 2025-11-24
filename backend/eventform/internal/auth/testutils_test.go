package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type testJWKS struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	server     *httptest.Server
	kid        string
}

func (tj *testJWKS) buildJWKS() map[string]any {
	n := base64.RawURLEncoding.EncodeToString(tj.publicKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(tj.publicKey.E)).Bytes())

	return map[string]any{
		"keys": []map[string]any{
			{
				"kid": tj.kid,
				"kty": "RSA",
				"alg": "RS256",
				"use": "sig",
				"n":   n,
				"e":   e,
			},
		},
	}
}

func (tj *testJWKS) createToken(sub, email string, expiry time.Time) (string, error) {
	claims := jwt.MapClaims{
		"sub":   sub,
		"email": email,
		"exp":   expiry.Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = tj.kid

	return token.SignedString(tj.privateKey)
}

func newTestJWKS(t testing.TB) *testJWKS {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key. Failed with error %v", err)
	}

	tj := &testJWKS{
		privateKey: privateKey, // pragma: allowlist secret - dynamically generated test key
		publicKey:  &privateKey.PublicKey,
		kid:        "test-key-id",
	}

	tj.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwks := tj.buildJWKS()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(jwks); err != nil {
			log.Fatalf("Error when encoding JWT: %v", err)
		}
	}))

	return tj
}

func (tj *testJWKS) close() {
	tj.server.Close()
}
