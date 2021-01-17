// controller 包写web应用的主逻辑，请尽量让它与model层解耦
// 不同router下的controller也建议分文件写
package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"github.com/KSkun/tqb-backend/util/context"
)

// TODO: 完成controller

func HelloWorldHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World!")
}

// HTTPErrorHandler 替换默认的错误处理，统一成目前使用的格式
func HTTPErrorHandler(err error, c echo.Context) {
	httpError, ok := err.(*echo.HTTPError)
	if ok {
		_ = context.Error(c, httpError.Code, fmt.Sprintf("%s", httpError.Message), err)
		return
	}

	_ = context.Error(c, http.StatusInternalServerError, "Unhandled internal server error", err)
}
