package router

import (
	"github.com/KSkun/tqb-backend/controller"
	"github.com/KSkun/tqb-backend/middleware"
	"github.com/labstack/echo/v4"
)

func initUserRouter(g *echo.Group) {
	g.Add(echo.GET, "/public_key", controller.UserGetPublicKey)
	g.Add(echo.GET, "/token", controller.UserGetToken)
	g.Add(echo.POST, "", controller.UserAddUser)
	g.Add(echo.GET, "/email_verify", controller.UserSendVerifyMail)
	g.Add(echo.POST, "/email_verify", controller.UserVerifyEmail)
	g.Add(echo.PUT, "/password", controller.UserChangePassword)
	g.Add(echo.GET, "", controller.UserGetInfo, middleware.JWTMiddleware())
}
