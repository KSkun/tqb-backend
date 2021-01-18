package context

import (
	"github.com/KSkun/tqb-backend/config"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Response 返回值
type Response struct {
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
	ErrHint string      `json:"hint,omitempty"`
}

// Success 成功
func Success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Data:    data,
		Error:   "",
		Success: true,
	})
}

// Error 错误
func Error(c echo.Context, status int, data string, err error) error {
	ret := Response{
		Data:    nil,
		Error:   data,
		Success: false,
	}

	if config.C.Debug {
		if err != nil {
			ret.Error = err.Error()
		}
	}

	return c.JSON(status, ret)
}
