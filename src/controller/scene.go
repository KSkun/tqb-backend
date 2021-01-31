package controller

import (
	"github.com/KSkun/tqb-backend/controller/param"
	"github.com/KSkun/tqb-backend/model"
	"github.com/KSkun/tqb-backend/util/context"
	"github.com/KSkun/tqb-backend/util/log"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func SceneGetList(ctx echo.Context) error {
	m := model.GetModel()
	defer m.Close()

	sceneList, err := m.GetSceneList()
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get scene list", err)
	}

	ret := make([]param.ObjRspScene, 0)
	for _, scene := range sceneList {
		ret = append(ret, param.ObjRspScene{
			ID:           scene.ID.Hex(),
			FromQuestion: scene.FromQuestion.Hex(),
			NextQuestion: scene.NextQuestion.Hex(),
			Title:        scene.Title,
		})
	}
	return context.Success(ctx, param.RspSceneGetList{Scene: ret})
}

func SceneGetInfo(ctx echo.Context) error {
	idHex := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return context.Error(ctx, http.StatusBadRequest, "invalid scene id", err)
	}

	userIDHex := context.GetUserFromJWT(ctx)
	userID, _ := primitive.ObjectIDFromHex(userIDHex)

	m := model.GetModel()
	defer m.Close()

	// 检查用户是否有权限查看该剧情
	scene, err := m.GetScene(id)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get scene info", err)
	}
	unlocked, err := m.UserHasUnlockedScene(userID, id)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user info", err)
	}
	allUnlocked, err := m.UserIsAllUnlocked(userID)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user info", err)
	}
	if !unlocked && scene.FromQuestion != model.NullID && !allUnlocked {
		return context.Error(ctx, http.StatusForbidden, "this scene is locked", nil)
	}

	return context.Success(ctx, param.RspSceneGetInfo{
		Title:        scene.Title,
		Text:         scene.Text,
		FromQuestion: scene.FromQuestion.Hex(),
		NextQuestion: scene.NextQuestion.Hex(),
	})
}

func SceneSetUnlock(ctx echo.Context) error {
	idHex := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return context.Error(ctx, http.StatusBadRequest, "invalid scene id", err)
	}

	userIDHex := context.GetUserFromJWT(ctx)
	userID, _ := primitive.ObjectIDFromHex(userIDHex)

	m := model.GetModel()
	defer m.Close()

	scene, err := m.GetScene(id)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get scene info", err)
	}

	// 检查用户是否有权限解锁该剧情
	user, err := m.GetUser(userID)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user info", err)
	}
	if user.LastQuestion != scene.FromQuestion && scene.FromQuestion != model.NullID { // 排除入口剧情
		return context.Error(ctx, http.StatusForbidden, "you cannot unlock this scene", nil)
	}
	finished, err := m.UserHasFinishedQuestion(userID, scene.NextQuestion)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user info", err)
	}
	if finished {
		return context.Success(ctx, nil)
	}

	err = m.SetUserLastScene(userID, id)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process user info", err)
	}

	// 如果走到最后一个剧情，则清空用户状态
	if scene.NextQuestion == model.NullID {
		err = m.UpdateUser(userID, bson.M{
			"last_scene":      model.NullID,
			"last_question":   model.NullID,
			"l_last_question": model.NullID,
		})
		if err != nil {
			log.Logger.Error(err)
		}

		err = m.IncUserCompleteCount(userID)
		if err != nil {
			log.Logger.Error(err)
		}
	}
	return context.Success(ctx, nil)
}
