package controller

import (
	"github.com/KSkun/tqb-backend/model"
	"github.com/KSkun/tqb-backend/util/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func addBlankSubmission(m model.Model, userID primitive.ObjectID, questionID primitive.ObjectID) error {
	submission := model.Submission{
		User:     userID,
		Question: questionID,
		Time:     time.Now().Unix(),
		File:     nil,
		Option:   nil,
		Point:    0.0,
	}
	_, err := m.AddSubmission(submission)
	if err != nil {
		return err
	}
	err = m.AddUserFinishedQuestion(userID, questionID)
	if err != nil {
		return err
	}
	return nil
}

func timedOutWorker(userID primitive.ObjectID, questionID primitive.ObjectID) {
	m := model.GetModel()
	question, err := m.GetQuestion(questionID)
	if err != nil {
		log.Logger.Error(err)
		return
	}
	m.Close()

	time.Sleep(time.Second * time.Duration(question.TimeLimit + 10))

	m = model.GetModel()
	finished, err := m.UserHasFinishedQuestion(userID, questionID)
	if err != nil {
		log.Logger.Error(err)
		return
	}
	if finished {
		return
	}

	// 超时问题复位
	user, err := m.GetUser(userID)
	if err != nil {
		log.Logger.Error(err)
		return
	}
	lastQuestion, err := m.GetQuestion(user.LLastQuestion)
	if err != nil {
		log.Logger.Error(err)
		return
	}
	// 检查上一问题是否还有其他路径可走
	flagHasOtherOption := false
	for _, scene := range lastQuestion.NextScene {
		unlocked, err := m.UserHasUnlockedScene(userID, scene.Scene)
		if err != nil {
			log.Logger.Error(err)
			return
		}
		if !unlocked {
			flagHasOtherOption = true
			break
		}
	}
	if !flagHasOtherOption { // 如果没有其他路径可走，则本题判为零分
		err = addBlankSubmission(m, userID, questionID)
		if err != nil {
			log.Logger.Error(err)
			return
		}
		return
	}

	// 如果存在其他路径，则后退一步
	err = m.SetUserBackQuestion(userID)
	if err != nil {
		log.Logger.Error(err)
		return
	}
}
