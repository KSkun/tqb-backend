package router

import (
	"github.com/KSkun/tqb-backend/controller"
	"github.com/KSkun/tqb-backend/middleware"
	"github.com/labstack/echo/v4"
)

func initQuestionGroup(g *echo.Group) {
	g.Use(middleware.JWTMiddleware())

	g.GET("", controller.QuestionGetList)
	g.GET("/:id", controller.QuestionGetInfo)
	g.POST("/:id/start", controller.QuestionSetStart)
}
