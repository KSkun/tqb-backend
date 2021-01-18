package controller

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/KSkun/tqb-backend/controller/param"
	"github.com/KSkun/tqb-backend/model"
	"github.com/KSkun/tqb-backend/util/context"
	"github.com/labstack/echo/v4"
	"net/http"
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
			return context.Error(ctx, http.StatusInternalServerError,
				"failed to generate rsa key", err)
		}
		err = m.AddPrivateKey(req.Email, key)
		if err != nil {
			return context.Error(ctx, http.StatusInternalServerError,
				"failed on model", err)
		}
	}
	publicKey, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError,
			"failed to generate public key", err)
	}
	publicKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:    "RSA PUBLIC KEY",
		Bytes:   publicKey,
	})
	return context.Success(ctx, param.RspUserGetPublicKey{PublicKey: string(publicKeyPem)})
}
