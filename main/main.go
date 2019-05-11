package main

import (
	admin "Miuer/admin/controller/gin"
	banner "Miuer/banner/controller/gin"
	category "Miuer/category/controller/gin"
	order "Miuer/order/controller/gin"
	permission "Miuer/permission/controller/gin"
	"database/sql"

	ginjwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var (
	JWTmw *ginjwt.GinJWTMiddleware
)

func main() {
	router := gin.Default()

	dbConn, err := sql.Open("mysql", "root:Miufighting.@tcp(127.0.0.1:3306)/Miuer")
	if err != nil {
		panic(err)
	}

	adminCon := admin.New(dbConn)
	adminCon.RegisterRouter(router)

	/*
		authMiddleware := &ginjwt.GinJWTMiddleware{
			Realm:            "Template",
			Key:              []byte("hydra"),
			Timeout:          24 * time.Hour,
			SigningAlgorithm: "HS256",
		}

			// getuid
			GetUID := adminCon.EmbodyJWTMiddleWare(authMiddleware)

			router.POST("/api/v1/admin/login", authMiddleware.LoginHandler)
			router.Use(func(ctx *gin.Context) {
				authMiddleware.MiddlewareFunc()
			})
			router.Use(adminCon.CheckIsActive(GetUID))
	*/
	bannerCon := banner.New(dbConn)
	bannerCon.Register(router)

	categoryCon := category.New(dbConn, "category", "cate")
	categoryCon.Register(router)

	orderCon := order.New(dbConn, "order", "item")
	orderCon.Register(router)

	permissionCon := permission.New(dbConn)
	permissionCon.Register(router)

	router.Run(":8080")

}
