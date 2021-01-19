package controller

import (
	"github.com/KSkun/tqb-backend/controller/param"
	"github.com/KSkun/tqb-backend/model"
	"github.com/KSkun/tqb-backend/util/context"
	"github.com/labstack/echo/v4"
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
	unlocked, err := m.UserHasUnlockedScene(userID, id)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get user info", err)
	}
	if !unlocked {
		return context.Error(ctx, http.StatusForbidden, "this scene is locked", nil)
	}

	scene, err := m.GetScene(id)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get scene info", err)
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
	if user.LastQuestion != scene.FromQuestion {
		return context.Error(ctx, http.StatusForbidden, "you cannot unlock this scene", nil)
	}

	err = m.SetUserLastScene(userID, id)
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to process user info", err)
	}
	return context.Success(ctx, nil)
}
