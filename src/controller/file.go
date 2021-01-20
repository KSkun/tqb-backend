package controller

import (
	"github.com/KSkun/tqb-backend/config"
	"github.com/KSkun/tqb-backend/controller/param"
	"github.com/KSkun/tqb-backend/model"
	"github.com/KSkun/tqb-backend/util/context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func FileUpload(ctx echo.Context) error {
	userIDHex := context.GetUserFromJWT(ctx)
	userID, _ := primitive.ObjectIDFromHex(userIDHex)

	m := model.GetModel()
	defer m.Close()

	filenameUUID, err := uuid.NewUUID()
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process file", err)
	}
	filename := filenameUUID.String() + ".pdf"

	file, err := ctx.FormFile("file")
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process file", err)
	}
	if !strings.HasSuffix(file.Filename, ".pdf") {
		return context.Error(ctx, http.StatusBadRequest, "unsupported file type", err)
	}
	srcFile, err := file.Open()
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process file", err)
	}
	defer srcFile.Close()
	if _, err := os.Stat(config.C.App.UploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(config.C.App.UploadDir, os.ModePerm)
		if err != nil {
			return context.Error(ctx, http.StatusInternalServerError, "failed to process file", err)
		}
	}
	dstFile, err := os.Create(config.C.App.UploadDir + "/" + filename)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process file", err)
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process file", err)
	}

	dbFile := model.File{
		Filename: filename,
		User:     userID,
		Time:     time.Now().Unix(),
	}
	fileID, err := m.AddFile(dbFile)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process file", err)
	}
	return context.Success(ctx, param.RspFileUpload{ID: fileID.Hex()})
}

func FileGet(ctx echo.Context) error {
	fileIDHex := ctx.Param("id")
	fileID, err := primitive.ObjectIDFromHex(fileIDHex)
	if err != nil {
		return context.Error(ctx, http.StatusBadRequest, "invalid file id", err)
	}

	userIDHex := context.GetUserFromJWT(ctx)
	userID, _ := primitive.ObjectIDFromHex(userIDHex)

	m := model.GetModel()
	defer m.Close()

	file, err := m.GetFile(fileID)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process file", err)
	}
	if file.User != userID {
		return context.Error(ctx, http.StatusForbidden, "you cannot get this file", nil)
	}

	return ctx.File(config.C.App.UploadDir + "/" + file.Filename)
}
