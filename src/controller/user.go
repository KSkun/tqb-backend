package controller

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/KSkun/tqb-backend/constant"
	"github.com/KSkun/tqb-backend/controller/param"
	"github.com/KSkun/tqb-backend/model"
	"github.com/KSkun/tqb-backend/util"
	"github.com/KSkun/tqb-backend/util/context"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to generate rsa key", err)
	}
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
	publicKey := x509.MarshalPKCS1PublicKey(&key.PublicKey)
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
	if !user.IsEmailVerified {
		return context.Error(ctx, http.StatusBadRequest, "please verify your email first", nil)
	}

	key, found, err := m.GetPrivateKey(req.Email)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user info", err)
	}
	if !found {
		return context.Error(ctx, http.StatusBadRequest, "key not found, please generate key first", nil)
	}

	pwDecrypt, err := util.RSADecryptFromString(req.Password, key)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process user info", err)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), pwDecrypt) != nil {
		return context.Error(ctx, http.StatusBadRequest, "wrong password or user not found", nil)
	}

	token, expire, err := util.GenerateJWTToken(util.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
			Subject:  "tuiqunbei",
		},
		User: user.ID.Hex(),
	})
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to generate token", err)
	}
	return context.Success(ctx, param.RspUserGetToken{
		Token:  token,
		Expire: expire,
	})
}

func UserAddUser(ctx echo.Context) error {
	req := param.ReqUserAddUser{}
	if err := ctx.Bind(&req); err != nil {
		return context.Error(ctx, http.StatusBadRequest, "bad request", err)
	}
	if err := ctx.Validate(req); err != nil {
		return context.Error(ctx, http.StatusBadRequest, "bad request", err)
	}

	m := model.GetModel()
	defer m.Close()

	// 检查邮箱是否已使用
	_, found, err := m.GetUserByEmail(req.Email)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user info", err)
	}
	if !found {
		return context.Error(ctx, http.StatusForbidden, "this email has been occupied", nil)
	}

	key, found, err := model.GetModel().GetPrivateKey(req.Email)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user info", err)
	}
	if !found {
		return context.Error(ctx, http.StatusBadRequest, "key not found, please generate key first", nil)
	}
	pwDecrypt, err := util.RSADecryptFromString(req.Password, key)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process user info", err)
	}
	pwEncrypt, err := bcrypt.GenerateFromPassword(pwDecrypt, bcrypt.DefaultCost)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process user info", err)
	}

	user := model.User{
		Username:         req.Username,
		Password:         string(pwEncrypt),
		Email:            req.Email,
		IsEmailVerified:  false,
		UnlockedScene:    make([]primitive.ObjectID, 0),
		FinishedQuestion: make([]primitive.ObjectID, 0),
	}
	id, err := m.AddUser(user)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process user info", err)
	}
	return context.Success(ctx, param.RspUserAddUser{ID: id.Hex()})
}

func UserSendVerifyMail(ctx echo.Context) error {
	req := param.ReqUserSendVerifyMail{}
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
		return context.Success(ctx, nil) // 如果用户不存在则不发送邮件，也不反馈
	}
	_, found, err = m.GetVerifyIDByEmail(req.Email)
	if found {
		return context.Success(ctx, nil) // 如果冷却时间未到则不发送邮件，也不反馈
	}

	verifyID, err := uuid.NewUUID()
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to generate verify id", err)
	}
	err = m.AddVerifyID(req.Email, verifyID.String())
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to generate verify id", err)
	}

	mailContent := fmt.Sprintf(constant.TextVerifyMailContent, user.Username, verifyID.String())
	util.SendMail(req.Email, constant.TextVerifyMailTitle, mailContent)
	return context.Success(ctx, nil)
}

func UserVerifyEmail(ctx echo.Context) error {
	req := param.ReqUserVerifyEmail{}
	if err := ctx.Bind(&req); err != nil {
		return context.Error(ctx, http.StatusBadRequest, "bad request", err)
	}
	if err := ctx.Validate(req); err != nil {
		return context.Error(ctx, http.StatusBadRequest, "bad request", err)
	}

	m := model.GetModel()
	defer m.Close()

	email, found, err := m.GetVerifyID(req.VerifyID)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user info", err)
	}
	if !found {
		return context.Error(ctx, http.StatusBadRequest, "invalid verify id", nil)
	}

	user, found, err := m.GetUserByEmail(email)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user info", err)
	}
	if !found {
		return context.Error(ctx, http.StatusBadRequest, "invalid verify id", nil)
	}
	if user.IsEmailVerified {
		return context.Error(ctx, http.StatusBadRequest, "your email has been verified", nil)
	}

	err = m.UpdateUser(user.ID, bson.M{"is_email_verified": true})
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process user info", err)
	}
	return context.Success(ctx, nil)
}

func UserChangePassword(ctx echo.Context) error {
	req := param.ReqUserChangePassword{}
	if err := ctx.Bind(&req); err != nil {
		return context.Error(ctx, http.StatusBadRequest, "bad request", err)
	}
	if err := ctx.Validate(req); err != nil {
		return context.Error(ctx, http.StatusBadRequest, "bad request", err)
	}

	m := model.GetModel()
	defer m.Close()

	email, found, err := m.GetVerifyID(req.VerifyID)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user info", err)
	}
	if !found {
		return context.Error(ctx, http.StatusBadRequest, "invalid verify id", nil)
	}

	user, found, err := m.GetUserByEmail(email)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user info", err)
	}
	if !found {
		return context.Error(ctx, http.StatusBadRequest, "invalid verify id", nil)
	}

	key, found, err := model.GetModel().GetPrivateKey(email)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user info", err)
	}
	if !found {
		return context.Error(ctx, http.StatusBadRequest, "key not found, please generate key first", nil)
	}
	pwDecrypt, err := util.RSADecryptFromString(req.Password, key)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process user info", err)
	}
	pwEncrypt, err := bcrypt.GenerateFromPassword(pwDecrypt, bcrypt.DefaultCost)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process user info", err)
	}

	err = m.UpdateUser(user.ID, bson.M{"password": string(pwEncrypt)})
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process user info", err)
	}
	return context.Success(ctx, nil)
}

func UserGetInfo(ctx echo.Context) error {
	idHex := context.GetUserFromJWT(ctx)
	id, _ := primitive.ObjectIDFromHex(idHex)

	m := model.GetModel()
	defer m.Close()

	user, err := m.GetUser(id)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user info", err)
	}

	unlockedScene := make([]string, 0)
	for _, id := range user.UnlockedScene {
		unlockedScene = append(unlockedScene, id.Hex())
	}

	finishedQuestion := make([]string, 0)
	for _, id := range user.FinishedQuestion {
		finishedQuestion = append(finishedQuestion, id.Hex())
	}

	return context.Success(ctx, param.RspUserGetInfo{
		ID:               user.ID.Hex(),
		Username:         user.Username,
		Email:            user.Email,
		IsEmailVerified:  user.IsEmailVerified,
		LastQuestion:     user.LastQuestion.Hex(),
		LastScene:        user.LastScene.Hex(),
		StartTime:        user.StartTime,
		UnlockedScene:    unlockedScene,
		FinishedQuestion: finishedQuestion,
	})
}

func UserGetUnlockedScene(ctx echo.Context) error {
	idHex := context.GetUserFromJWT(ctx)
	id, _ := primitive.ObjectIDFromHex(idHex)

	m := model.GetModel()
	defer m.Close()

	user, err := m.GetUser(id)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user info", err)
	}

	sceneRet := make([]param.ObjRspUserScene, 0)
	for _, id := range user.UnlockedScene {
		scene, err := m.GetScene(id)
		if err != nil {
			return context.Error(ctx, http.StatusInternalServerError, "failed to get scene info", err)
		}

		sceneRet = append(sceneRet, param.ObjRspUserScene{
			ID:           scene.ID.Hex(),
			Title:        scene.Title,
			Text:         scene.Text,
			FromQuestion: scene.FromQuestion.Hex(),
			NextQuestion: scene.NextQuestion.Hex(),
		})
	}

	return context.Success(ctx, param.RspUserGetUnlockedScene{
		Scene: sceneRet,
	})
}

func UserGetSubmission(ctx echo.Context) error {
	idHex := context.GetUserFromJWT(ctx)
	id, _ := primitive.ObjectIDFromHex(idHex)

	m := model.GetModel()
	defer m.Close()

	submissionList, err := m.GetSubmissionByUser(id)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get submission info", err)
	}

	submissionListRet := make([]param.ObjRspSubmission, 0)
	for _, submission := range submissionList {
		question, err := m.GetQuestion(submission.Question)
		if err != nil {
			return context.Error(ctx, http.StatusInternalServerError, "failed to get submission info", err)
		}
		questionRet := param.ObjRspSubmissionQuestion{
			ID:    question.ID.Hex(),
			Title: question.Title,
		}

		fileRet := make([]string, 0)
		for _, file := range submission.File {
			fileRet = append(fileRet, file.Hex())
		}

		submissionListRet = append(submissionListRet, param.ObjRspSubmission{
			ID:         submission.ID.Hex(),
			Time:       submission.Time,
			Question:   questionRet,
			File:       fileRet,
			Option:     submission.Option,
			Point:      submission.Point,
			AnswerTime: submission.AnswerTime,
		})
	}
	return context.Success(ctx, param.RspUserGetSubmission{Submission: submissionListRet})
}

func UserRefreshStatus(ctx echo.Context) error {
	idHex := context.GetUserFromJWT(ctx)
	id, _ := primitive.ObjectIDFromHex(idHex)

	m := model.GetModel()
	defer m.Close()

	err := m.UpdateUser(id, bson.M{
		"last_scene":      model.NullID,
		"last_question":   model.NullID,
		"l_last_question": model.NullID,
	})
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process user info", err)
	}
	return context.Success(ctx, nil)
}
