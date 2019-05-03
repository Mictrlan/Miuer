package main

import (
	admin "Miuer/admin/controller/gin"
	"database/sql"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var (
	JWTmw *jwt.GinJWTMiddleware
)

func main() {
	router := gin.Default()

	dbConn, err := sql.Open("mysql", "root:Miufighting.@tcp(127.0.0.1:3306)/Miuer")
	if err != nil {
		panic(err)
	}

	adminCon := admin.New(dbConn)

	adminCon.RegisterRouter(router)

	JWTmw = &jwt.GinJWTMiddleware{

		Realm:   "Template",
		Key:     []byte("hydra"),
		Timeout: 24 * time.Hour,
	}

	adminCon.EmbodyJWTMiddleWare(JWTmw)

	router.POST("/api/v1/admin/login", JWTmw.LoginHandler)

	router.Use(func(c *gin.Context) {
		JWTmw.MiddlewareFunc()(c)
	})

	router.Run(":8080")

}
