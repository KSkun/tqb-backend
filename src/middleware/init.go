// middleware包放置所需要的中间件，比如jwt和参数验证器
package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// InitBeforeStart 会在Web服务启动之前对echo实例进行一些初始化操作
func InitBeforeStart(e *echo.Echo) error {
	// 使用cors
	e.Use(middleware.CORS())
	return nil
}
