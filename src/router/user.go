package router

import (
	"github.com/KSkun/tqb-backend/config"
	"github.com/KSkun/tqb-backend/controller"
	"github.com/KSkun/tqb-backend/middleware"
	"github.com/labstack/echo/v4"
)

func initUserRouter(g *echo.Group) {
	g.GET("/public_key", controller.UserGetPublicKey)
	g.GET("/token", controller.UserGetToken)
	g.POST("", controller.UserAddUser)
	g.GET("/email_verify", controller.UserSendVerifyMail)
	g.POST("/email_verify", controller.UserVerifyEmail)
	g.PUT("/password", controller.UserChangePassword)
	g.GET("", controller.UserGetInfo, middleware.JWTMiddleware())
	g.GET("/unlocked_scene", controller.UserGetUnlockedScene, middleware.JWTMiddleware())
	g.GET("/submission", controller.UserGetSubmission, middleware.JWTMiddleware())

	if config.C.Debug {
		g.PUT("/refresh", controller.UserRefreshStatus, middleware.JWTMiddleware())
	}
}
