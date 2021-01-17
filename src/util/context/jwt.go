package context

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/KSkun/tqb-backend/config"
	"github.com/KSkun/tqb-backend/util"
)

// GetAssociationFromJWT 获得JWTClaims
func GetJWTClaims(c echo.Context) *util.JWTClaims {
	return c.Get(config.JWTContextKey).(*jwt.Token).Claims.(*util.JWTClaims)
}

//获得payload中指定字段的值
func GetJWTUserID(c echo.Context) string {
	return getJWTFiled(c, "user_id")
}

func getJWTFiled(c echo.Context, filedName string) string {
	token := c.Get(config.JWTContextKey)
	if token != nil {
		if tokenStr, ok := token.(*jwt.Token).Claims.(jwt.MapClaims)[filedName].(string); ok {
			return tokenStr
		}
		return ""
	}
	return ""
}
