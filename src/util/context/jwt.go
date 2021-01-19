package context

import (
	"github.com/KSkun/tqb-backend/config"
	"github.com/KSkun/tqb-backend/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// GetAssociationFromJWT 获得JWTClaims
func GetJWTClaims(c echo.Context) *util.JWTClaims {
	return c.Get(config.JWTContextKey).(*jwt.Token).Claims.(*util.JWTClaims)
}

//获得payload中指定字段的值
func GetUserFromJWT(c echo.Context) string {
	return getJWTField(c, "user")
}

func getJWTField(c echo.Context, fieldName string) string {
	token := c.Get(config.JWTContextKey)
	if token != nil {
		if tokenStr, ok := token.(*jwt.Token).Claims.(jwt.MapClaims)[fieldName].(string); ok {
			return tokenStr
		}
		return ""
	}
	return ""
}
