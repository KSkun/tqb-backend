package controller

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/KSkun/tqb-backend/model"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestEncryption(t *testing.T) {
	m := model.GetModel()
	defer m.Close()

	key, _, _ := m.GetPrivateKey("ks@ksmeow.moe")
	fmt.Println(string(pem.EncodeToMemory(&pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(key),
	})))
	enc, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader, &key.PublicKey, []byte("1234"), nil)
	fmt.Println(base64.StdEncoding.EncodeToString(enc))

	hash, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
	fmt.Println(string(hash))
}
