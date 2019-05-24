package gin

import (
	"errors"
	"net/http"
	"time"

	mysql "github.com/Mictrlan/Miuer/admin/model/mysql"

	ginjwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

var (
	errUserIDNoExists = errors.New("user id no exists")
)

// ExtendJWTMiddleWare improve the middleware and return a function that get uid after successful execution
func (c *AdminController) ExtendJWTMiddleWare(authMW *ginjwt.GinJWTMiddleware) func(ctx *gin.Context) (uint32, error) {
	authMW.Authenticator = func(ctx *gin.Context) (interface{}, error) {
		return c.Login(ctx)
	}

	authMW.PayloadFunc = func(data interface{}) ginjwt.MapClaims {
		if v, ok := data.(uint32); ok {
			return ginjwt.MapClaims{
				"identity": v,
			}
		}

		return ginjwt.MapClaims{}
	}

	/*
		authMW.IdentityHandler = func(ctx *gin.Context) interface{} {
			claims := ginjwt.ExtractClaims(ctx)

			return claims["identity"]
		}
	*/
	authMW.IdentityHandler = func(ctx *gin.Context) interface{} {
		claims := ginjwt.ExtractClaims(ctx)
		return claims
	}

	authMW.LoginResponse = func(c *gin.Context, code int, token string, expire time.Time) {
		c.JSON(http.StatusOK, gin.H{
			"code":   http.StatusOK,
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		})
	}

	return func(ctx *gin.Context) (uint32, error) {
		ID, exists := ctx.Get("identity")
		if !exists {
			return 0, errUserIDNoExists
		}

		IDNew := ID.(float64)
		return uint32(IDNew), nil
	}
}

// CheckIsActive is a middlerware that check user active
func (c *AdminController) CheckIsActive(GetUID func(ctx *gin.Context) (uint32, error)) func(ctx *gin.Context) {
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
