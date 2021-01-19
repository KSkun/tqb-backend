package router

import (
	"github.com/KSkun/tqb-backend/controller"
	"github.com/KSkun/tqb-backend/middleware"
	"github.com/labstack/echo/v4"
)

func initSubjectGroup(g *echo.Group) {
	g.Use(middleware.JWTMiddleware())

	g.Add(echo.GET, "", controller.SubjectGetList)
}
