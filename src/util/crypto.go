package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
)

func RSADecryptFromString(toDecode string, key *rsa.PrivateKey) ([]byte, error) {
	bytes, err := base64.StdEncoding.DecodeString(toDecode)
	if err != nil {
		return nil, err
	}
	result, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, bytes, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}
