package router

import (
	"github.com/KSkun/tqb-backend/controller"
	"github.com/labstack/echo/v4"
)

func initUserRouter(g *echo.Group) {
	g.Add(echo.GET, "/public_key", controller.UserGetPublicKey)
}
