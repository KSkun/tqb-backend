package controller

import (
	"github.com/KSkun/tqb-backend/controller/param"
	"github.com/KSkun/tqb-backend/model"
	"github.com/KSkun/tqb-backend/util/context"
	"github.com/labstack/echo/v4"
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
	finishedQuestion := make(map[primitive.ObjectID]bool, 0)
	for _, question := range user.FinishedQuestion {
		finishedQuestion[question] = true
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

		status := param.StatusUnlock          // 未解锁
		if question.ID == user.LastQuestion { // 正在作答
			status = param.StatusAnswering
		}
		if finishedQuestion[question.ID] { // 已提交
			status = param.StatusFinish
		}

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
	finished, err := m.UserHasFinishedQuestion(userID, id)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user info", err)
	}
	if user.CompleteCount < 2 {
		if user.LastScene == model.NullID {
			return context.Error(ctx, http.StatusForbidden, "you have to select a subject first", nil)
		}
		scene, err := m.GetScene(user.LastScene)
		if err != nil {
			return context.Error(ctx, http.StatusInternalServerError, "failed to get scene info", err)
		}
		if scene.NextQuestion != id && user.LastQuestion != id && !finished {
			return context.Error(ctx, http.StatusForbidden, "this question is locked", nil)
		}
	}

	question, err := m.GetQuestion(id)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get question info", err)
	}

	subQuestionRet := make([]param.ObjRspSubQuestion, 0)
	for _, subQuestion := range question.SubQuestion {
		subQuestionRet = append(subQuestionRet, param.ObjRspSubQuestion{
			Type:      subQuestion.Type,
			Desc:      subQuestion.Desc,
			Option:    subQuestion.Option,
			FullPoint: subQuestion.FullPoint,
			PartPoint: subQuestion.PartPoint,
		})
	}

	nextSceneRet := make([]param.ObjRspNextScene, 0)
	for _, nextScene := range question.NextScene {
		nextSceneRet = append(nextSceneRet, param.ObjRspNextScene{
			Scene:  nextScene.Scene.Hex(),
			Option: nextScene.Option,
		})
	}

	status := param.StatusUnlock // 未解锁
	if user.LastQuestion == id { // 正在作答
		status = param.StatusAnswering
	}
	if finished { // 已提交
		status = param.StatusFinish
	}

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
	finished, err := m.UserHasFinishedQuestion(userID, id)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user info", err)
	}
	if finished {
		return context.Success(ctx, nil)
	}

	err = m.SetUserLastQuestion(userID, id)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process user info", err)
	}
	go timedOutWorker(userID, id)
	return context.Success(ctx, nil)
}

// 选择题判断得分
func getPoint(option [][]int, question model.Question) float64 {
	result := 0.0
	for i, sub := range question.SubQuestion {
		trueOption := make(map[int]bool, 0)
		for _, op := range sub.TrueOption {
			trueOption[op] = true
		}
		counter := 0
		for _, op := range option[i] {
			if !trueOption[op] { // 错选
				break
			}
			counter++
		}
		if counter == len(sub.TrueOption) { // 全对
			result += sub.FullPoint
		} else { // 缺项
			result += sub.PartPoint
		}
	}
	return result
}

// 判断一个问题是否全是选择题
func isAllChoice(question model.Question) bool {
	for _, sub := range question.SubQuestion {
		if sub.Type != model.TypeChoice {
			return false
		}
	}
	return true
}

// 判断提交是否合法
func validateSubmission(req param.ReqQuestionAddSubmission, question model.Question) bool {
	subCnt := len(question.SubQuestion)
	if len(req.File) != subCnt || len(req.Option) != subCnt { // 与子题数不符合
		return false
	}
	for i := 0; i < subCnt; i++ {
		if question.SubQuestion[i].Type == model.TypeChoice {
			if len(req.Option[i]) == 0 { // 选择题空答
				return false
			}
		}
		if question.SubQuestion[i].Type == model.TypePDF {
			if len(req.File[i]) == 0 { // 上传 PDF 空答
				return false
			}
		}
	}
	return true
}

// 判断问题是否已超时
func isQuestionTimedOut(user model.User, question model.Question) bool {
	return time.Now().After(time.Unix(user.StartTime, 0).
		Add(time.Second * time.Duration(question.TimeLimit+10)))
}

func QuestionAddSubmission(ctx echo.Context) error {
	req := param.ReqQuestionAddSubmission{}
	if err := ctx.Bind(&req); err != nil {
		return context.Error(ctx, http.StatusBadRequest, "bad request", err)
	}
	if err := ctx.Validate(req); err != nil {
		return context.Error(ctx, http.StatusBadRequest, "bad request", err)
	}
	fileID := make([]primitive.ObjectID, 0)
	var err error
	if len(req.File) != 0 {
		for _, file := range req.File {
			if len(file) == 0 { // 该位置非文件
				fileID = append(fileID, primitive.ObjectID{})
				continue
			}
			fileIDReq, err := primitive.ObjectIDFromHex(file)
			if err != nil {
				return context.Error(ctx, http.StatusBadRequest, "invalid file id", err)
			}
			fileID = append(fileID, fileIDReq)
		}
	}

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
	question, err := m.GetQuestion(id)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get question info", err)
	}
	if user.LastQuestion != id {
		return context.Error(ctx, http.StatusForbidden, "you are not answering this question", nil)
	}
	if question.TimeLimit != 0 && isQuestionTimedOut(user, question) { // 超时
		return context.Error(ctx, http.StatusForbidden, "the time limit is exceeded", nil)
	}
	hasFinished, err := m.UserHasFinishedQuestion(userID, id)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process user info", err)
	}
	if hasFinished {
		return context.Error(ctx, http.StatusForbidden, "you have finished this question", nil)
	}

	// 检查答案是否合法
	if !validateSubmission(req, question) {
		return context.Error(ctx, http.StatusBadRequest, "invalid submission", nil)
	}

	// 检查文件是否都是用户自己的
	for _, file := range fileID {
		if file == model.NullID {
			continue
		}
		fileObj, err := m.GetFile(file)
		if err != nil {
			return context.Error(ctx, http.StatusInternalServerError, "failed to get file info", err)
		}
		if fileObj.User != userID {
			return context.Error(ctx, http.StatusBadRequest, "you cannot use file "+file.Hex(), nil)
		}
	}

	submission := model.Submission{
		User:       userID,
		Question:   id,
		Time:       time.Now().Unix(),
		File:       fileID,
		Option:     req.Option,
		Point:      model.PointUnknown,
		AnswerTime: int(time.Now().Sub(time.Unix(user.StartTime, 0)) / time.Second),
		IsTimeOut:  false,
	}
	if isAllChoice(question) { // 选择题自动批改
		submission.Point = getPoint(submission.Option, question)
	}
	submissionID, err := m.AddSubmission(submission)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process submission", err)
	}
	err = m.AddUserFinishedQuestion(userID, id)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process user info", err)
	}
	return context.Success(ctx, param.RspQuestionAddSubmission{ID: submissionID.Hex()})
}
