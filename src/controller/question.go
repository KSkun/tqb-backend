package controller

import (
	"github.com/KSkun/tqb-backend/controller/param"
	"github.com/KSkun/tqb-backend/model"
	"github.com/KSkun/tqb-backend/util/context"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func QuestionGetList(ctx echo.Context) error {
	userIDHex := context.GetUserFromJWT(ctx)
	userID, _ := primitive.ObjectIDFromHex(userIDHex)

	m := model.GetModel()
	defer m.Close()

	user, err := m.GetUser(userID)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user info", err)
	}

	questionList, err := m.GetQuestionList()
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get question list", err)
	}

	ret := make([]param.ObjRspQuestion, 0)
	for _, question := range questionList {
		nextScene := make([]string, 0)
		for _, scene := range question.NextScene {
			nextScene = append(nextScene, scene.Scene.Hex())
		}

		status := 0 // 未解锁
		if question.ID == user.LastQuestion { // 正在作答
			status = 1
		}
		// TODO: 已提交

		ret = append(ret, param.ObjRspQuestion{
			ID:        question.ID.Hex(),
			Title:     question.Title,
			NextScene: nextScene,
			Status:    status,
		})
	}
	return context.Success(ctx, param.RspQuestionGetList{Question: ret})
}

func QuestionGetInfo(ctx echo.Context) error {
	idHex := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return context.Error(ctx, http.StatusBadRequest, "invalid question id", err)
	}

	userIDHex := context.GetUserFromJWT(ctx)
	userID, _ := primitive.ObjectIDFromHex(userIDHex)

	m := model.GetModel()
	defer m.Close()

	// 检查用户是否有权限查看该问题
	user, err := m.GetUser(userID)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user info", err)
	}
	scene, err := m.GetScene(user.LastScene)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user info", err)
	}
	if scene.NextQuestion != id && user.LastQuestion != id { // TODO: 增加已提交问题的判断
		return context.Error(ctx, http.StatusForbidden, "this question is locked", nil)
	}

	question, err := m.GetQuestion(id)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get question info", err)
	}

	subQuestionRet := make([]param.ObjRspSubQuestion, 0)
	for _, subQuestion := range question.SubQuestion {
		subQuestionRet = append(subQuestionRet, param.ObjRspSubQuestion{
			Type:       subQuestion.Type,
			Desc:       subQuestion.Desc,
			Option:     subQuestion.Option,
			FullPoint:  subQuestion.FullPoint,
			PartPoint:  subQuestion.PartPoint,
		})
	}

	nextSceneRet := make([]param.ObjRspNextScene, 0)
	for _, nextScene := range question.NextScene {
		nextSceneRet = append(nextSceneRet, param.ObjRspNextScene{
			Scene:  nextScene.Scene.Hex(),
			Option: nextScene.Option,
		})
	}

	status := 0 // 未解锁
	if user.LastQuestion == id { // 正在作答
		status = 1
	}
	// TODO: 已提交

	return context.Success(ctx, param.RspQuestionGetInfo{
		Title:       question.Title,
		Desc:        question.Desc,
		SubQuestion: subQuestionRet,
		Author:      question.Author,
		Audio:       question.Audio,
		TimeLimit:   question.TimeLimit,
		NextScene:   nextSceneRet,
		Status:      status,
	})
}

func QuestionSetStart(ctx echo.Context) error {
	idHex := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return context.Error(ctx, http.StatusBadRequest, "invalid question id", err)
	}

	userIDHex := context.GetUserFromJWT(ctx)
	userID, _ := primitive.ObjectIDFromHex(userIDHex)

	m := model.GetModel()
	defer m.Close()

	// 检查用户是否有权限开始答题
	user, err := m.GetUser(userID)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user info", err)
	}
	scene, err := m.GetScene(user.LastScene)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user info", err)
	}
	if scene.NextQuestion != id {
		return context.Error(ctx, http.StatusForbidden, "you cannot answer this question", nil)
	}
	if user.LastQuestion == id {
		return context.Error(ctx, http.StatusBadRequest, "you have already started this question", nil)
	}

	err = m.UpdateUser(userID, bson.M{"last_question": id, "start_time": time.Now().Unix()})
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process user info", err)
	}
	return context.Success(ctx, nil)
}
