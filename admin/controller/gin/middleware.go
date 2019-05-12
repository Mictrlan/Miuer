package gin

import (
	"errors"
	"fmt"
	"net/http"

	mysql "github.com/Mictrlan/Miuer/admin/model/mysql"

	ginjwt "github.com/appleboy/gin-jwt"
	gojwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	errUserIdNoExists = errors.New("user id no exists")
)

func (c *controller) EmbodyJWTMiddleWare(authMW *ginjwt.GinJWTMiddleware) func(ctx *gin.Context) (uint32, error) {
	authMW.Authenticator = func(ctx *gin.Context) (interface{}, error) {
		return c.Login(ctx)
	}

	//err
	fmt.Println("111")
	authMW.PayloadFunc = func(data interface{}) ginjwt.MapClaims {
		if v, ok := data.(uint32); ok {
			return ginjwt.MapClaims{
				"identity": uint32(v),
			}
		}
		return ginjwt.MapClaims{}
	}

	authMW.IdentityHandler = func(ctx *gin.Context) interface{} {
		claims := gojwt.MapClaims(ginjwt.ExtractClaims(ctx))
		return claims["identity"]
	}

	return func(ctx *gin.Context) (uint32, error) {
		Id, exists := ctx.Get("identity")
		if !exists {
			return 0, errUserIdNoExists
		}

		// why ?!
		ID := Id.(float64)
		return uint32(ID), nil
	}
}

func (c *controller) CheckIsActive(GetUID func(ctx *gin.Context) (uint32, error)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := GetUID(ctx)
		if err != nil {
			ctx.AbortWithError(http.StatusBadGateway, err)
			return
		}

		active, err := mysql.IsActive(c.db, id)
		if err != nil {
			ctx.AbortWithError(http.StatusConflict, err)
			return
		}

		if !active {
			ctx.AbortWithError(http.StatusFailedDependency, err)
			return
		}

	}
}
