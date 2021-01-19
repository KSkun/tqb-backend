package router

import (
	"github.com/KSkun/tqb-backend/controller"
	"github.com/KSkun/tqb-backend/middleware"
	"github.com/labstack/echo/v4"
)

func initSceneGroup(g *echo.Group) {
	g.Use(middleware.JWTMiddleware())

	g.GET("", controller.SceneGetList)
	g.GET("/:id", controller.SceneGetInfo)
	g.POST("/:id/unlock", controller.SceneSetUnlock)
}
