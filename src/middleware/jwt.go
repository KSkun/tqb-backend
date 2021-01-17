package middleware

import (
	"github.com/KSkun/tqb-backend/config"
	"github.com/KSkun/tqb-backend/util"
	. "github.com/KSkun/tqb-backend/util/log"

	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CustomJWTConfig custom jwt config
func CustomJWTConfig(skipperPaths []string, authScheme string) middleware.JWTConfig {
	if authScheme == "" {
		authScheme = "Bearer"
	}
	return middleware.JWTConfig{
		SigningKey:  []byte(config.C.JWT.Secret),
		TokenLookup: "header:" + echo.HeaderAuthorization,
		AuthScheme:  authScheme,
		Claims:      &util.JWTClaims{},
		ContextKey:  config.JWTContextKey,
		Skipper:     CustomSkipper(skipperPaths),
	}
}

// CustomSkipper 自定义的JWT路径Skipper
// 使用正则表达式语法，配置中的每一个条目会进行前缀匹配
func CustomSkipper(skipperPaths []string) func(c echo.Context) bool {
	for i := range skipperPaths {
		skipperPaths[i] = "(?:^" + skipperPaths[i] + ")" // 前缀匹配
	}

	regexStr := strings.Join(skipperPaths, "|")
	r := regexp.MustCompile(regexStr)
	if config.C.Debug {
		Logger.Printf("Complied skipper string: %s", regexStr)
	}

	return func(c echo.Context) bool {
		if c.Request().Method == http.MethodOptions {
			return true
		}

		return r.MatchString(c.Request().URL.Path)
	}
}
