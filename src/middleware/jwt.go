package middleware

import (
	"github.com/KSkun/tqb-backend/config"
	"github.com/KSkun/tqb-backend/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CustomJWTConfig custom jwt config
func CustomJWTConfig(authScheme string) middleware.JWTConfig {
	if authScheme == "" {
		authScheme = "Bearer"
	}
	return middleware.JWTConfig{
		SigningKey:  []byte(config.C.JWT.Secret),
		TokenLookup: "header:" + echo.HeaderAuthorization,
		AuthScheme:  authScheme,
		Claims:      &util.JWTClaims{},
		ContextKey:  config.JWTContextKey,
	}
}

func JWTMiddleware() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(CustomJWTConfig("Bearer"))
}
