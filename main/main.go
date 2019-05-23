package main

import (
	"time"

	admin "github.com/Mictrlan/Miuer/admin/controller/gin"
	banner "github.com/Mictrlan/Miuer/banner/controller/gin"
	category "github.com/Mictrlan/Miuer/category/controller/gin"
	order "github.com/Mictrlan/Miuer/order/controller/gin"
	permission "github.com/Mictrlan/Miuer/permission/controller/gin"
	smsservice "github.com/Mictrlan/Miuer/smsservice/controller/gin"
	services "github.com/Mictrlan/Miuer/smsservice/services"
	upload "github.com/Mictrlan/Miuer/upload/controller/gin"

	"database/sql"

	ginjwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// JWWTmw -
var (
	JWTmw *ginjwt.GinJWTMiddleware
)

type funcv struct{}

var v funcv

func (v funcv) OnVerifySucceed(a, b string) {}
func (v funcv) OnVerifyFailed(a, b string)  {}

func main() {
	router := gin.Default()

	dbConn, err := sql.Open("mysql", "root:Miufighting.@tcp(127.0.0.1:3306)/Miuer")
	if err != nil {
		panic(err)
	}

	adminCon := admin.New(dbConn)

	authMiddleware := &ginjwt.GinJWTMiddleware{
		Realm:            "Template",
		Key:              []byte("hydra"),
		Timeout:          24 * time.Hour,
		TimeFunc:         time.Now,
		SigningAlgorithm: "HS256",
		TokenLookup:      "header:Authorization",
	}

	GetUID := adminCon.ExtendJWTMiddleWare(authMiddleware)
	router.POST("/api/v1/admin/login", authMiddleware.LoginHandler)

	router.Use(func(ctx *gin.Context) {
		authMiddleware.MiddlewareFunc()(ctx)
	})

	bannerCon := banner.New(dbConn)
	bannerCon.Register(router)

	categoryCon := category.New(dbConn, "category", "cate")
	categoryCon.Register(router)

	orderCon := order.New(dbConn, "order", "item")
	orderCon.Register(router)

	permissionCon := permission.New(dbConn)
	router.Use(permission.CheckPermission(permissionCon, GetUID))
	permissionCon.Register(router)

	sm := &services.Config{
		Host:           "https://fesms.market.alicloudapi.com/sms/",
		Appcode:        "6f37345cad574f408bff3ede627f7014",
		Digits:         6,
		ResendInterval: 60,
		OnCheck:        v,
		DB:             dbConn,
	}

	smsserviceCon := smsservice.New(dbConn, sm)
	smsserviceCon.Register(router)

	uploadCon := upload.New(dbConn, "http://127.0.0.1:9573", GetUID)
	uploadCon.Register(router)

	router.Use(adminCon.CheckIsActive(GetUID))
	adminCon.RegisterRouter(router)

	router.Run(":8080")

}
