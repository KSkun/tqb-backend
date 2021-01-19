package controller

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/KSkun/tqb-backend/controller/param"
	"github.com/KSkun/tqb-backend/model"
	"github.com/KSkun/tqb-backend/util"
	"github.com/KSkun/tqb-backend/util/context"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func UserGetPublicKey(ctx echo.Context) error {
	req := param.ReqUserGetPublicKey{}
	if err := ctx.Bind(&req); err != nil {
		return context.Error(ctx, http.StatusBadRequest, "bad request", err)
	}
	if err := ctx.Validate(req); err != nil {
		return context.Error(ctx, http.StatusBadRequest, "bad request", err)
	}

	m := model.GetModel()
	defer m.Close()

	key, found, err := m.GetPrivateKey(req.Email)
	if !found {
		key, err = rsa.GenerateKey(rand.Reader, 1024)
		if err != nil {
			return context.Error(ctx, http.StatusInternalServerError, "failed to generate rsa key", err)
		}
		err = m.AddPrivateKey(req.Email, key)
		if err != nil {
			return context.Error(ctx, http.StatusInternalServerError, "failed on model", err)
		}
	}
	publicKey, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to generate public key", err)
	}
	publicKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKey,
	})
	return context.Success(ctx, param.RspUserGetPublicKey{PublicKey: string(publicKeyPem)})
}

func UserGetToken(ctx echo.Context) error {
	req := param.ReqUserGetToken{}
	if err := ctx.Bind(&req); err != nil {
		return context.Error(ctx, http.StatusBadRequest, "bad request", err)
	}
	if err := ctx.Validate(req); err != nil {
		return context.Error(ctx, http.StatusBadRequest, "bad request", err)
	}

	m := model.GetModel()
	defer m.Close()

	user, found, err := m.GetUserByEmail(req.Email)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user info", err)
	}
	if !found {
		return context.Error(ctx, http.StatusBadRequest, "wrong password or user not found", nil)
	}

	key, found, err := m.GetPrivateKey(req.Email)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user key", err)
	}
	if !found {
		return context.Error(ctx, http.StatusBadRequest, "key not found, please generate key first", nil)
	}

	pwDecode, err := base64.StdEncoding.DecodeString(req.Password)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process user info", err)
	}
	pwDecrypt, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, pwDecode, nil)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process user info", err)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), pwDecrypt) != nil {
		return context.Error(ctx, http.StatusBadRequest, "wrong password or user not found", nil)
	}

	token, expire, err := util.GenerateJWTToken(util.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			Subject:   "tuiqunbei",
		},
		User:           user.ID.Hex(),
	})
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to generate token", err)
	}
	return context.Success(ctx, param.RspUserGetToken{
		Token:  token,
		Expire: expire,
	})
}
