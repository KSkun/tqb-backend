package router

import (
	"github.com/KSkun/tqb-backend/controller"
	"github.com/KSkun/tqb-backend/middleware"
	"github.com/labstack/echo/v4"
)

func initSceneGroup(g *echo.Group) {
	g.Use(middleware.JWTMiddleware())

	g.Add(echo.GET, "", controller.SceneGetList)
	g.Add(echo.GET, "/:id", controller.SceneGetInfo)
	g.Add(echo.POST, "/:id/unlock", controller.SceneSetUnlock)
}
