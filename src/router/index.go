package router

import (
	"github.com/labstack/echo/v4"
)

// InitRouter 初始化所有路由，可以每个路由分函数分文件写，方便之后维护
func InitRouter(g *echo.Group) {
	initIndexRouter(g)

	userGroup := g.Group("/user")
	initUserRouter(userGroup)

	subjectGroup := g.Group("/subject")
	initSubjectGroup(subjectGroup)

	sceneGroup := g.Group("/scene")
	initSceneGroup(sceneGroup)

	questionGroup := g.Group("/question")
	initQuestionGroup(questionGroup)
}

func initIndexRouter(g *echo.Group) {

}
