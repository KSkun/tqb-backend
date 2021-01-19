package controller

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/KSkun/tqb-backend/model"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestEncryption(t *testing.T) {
	m := model.GetModel()
	defer m.Close()

	key, _, _ := m.GetPrivateKey("ks@ksmeow.moe")
	enc, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader, &key.PublicKey, []byte("123"), nil)
	fmt.Println(base64.StdEncoding.EncodeToString(enc))

	hash, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
	fmt.Println(string(hash))
}
