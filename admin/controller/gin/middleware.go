package gin

import (
	"errors"
	"net/http"

	mysql "Miuer/admin/model/mysql"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

var (
	ErrUserIdNoExist = errors.New("user Id no exist!")
)

func (c *controller) EmbodyJWTMiddleWare(JWT *jwt.GinJWTMiddleware) func(ctx *gin.Context) (uint32, error) {

	// get identity at first logging
	JWT.Authenticator = func(ctx *gin.Context) (interface{}, error) {
		return c.Login(ctx)
	}

	// add identity to MapClaims
	JWT.PayloadFunc = func(data interface{}) jwt.MapClaims {
		return jwt.MapClaims{
			"userId": data,
		}
	}

	// 	add identity to Context
	JWT.IdentityHandler = func(ctx *gin.Context) interface{} {
		claims := jwt.ExtractClaims(ctx)
		return claims["userId"]
	}

	return func(ctx *gin.Context) (uint32, error) {
		id, ok := ctx.Get("userId")
		if ok != true {
			return 0, ErrUserIdNoExist
		}

		return id.(uint32), nil
	}
}

func (c *controller) CheckIsActive(ctx *gin.Context) bool {
	id, ok := ctx.Get("userId")
	if ok != true {
		panic(ErrUserIdNoExist)
	}

	isactive, err := mysql.IsActive(c.SQLStore(), id.(uint32))
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return false
	}

	return isactive

}
