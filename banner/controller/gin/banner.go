package gin

import (
	mysql "Miuer/banner/model/mysql"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type BannerController struct {
	db *sql.DB
}

func New(db *sql.DB) *BannerController {
	return &BannerController{
		db: db,
	}
}

func (bc *BannerController) Register(r gin.IRouter) error {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	if err := mysql.CreateDB(bc.db); err != nil {
		return err
	}

	if err := mysql.CreateTable(bc.db); err != nil {
		return err
	}

	r.POST("/api/v1/banner/create", bc.insert)
	r.POST("/api/v1/banner/delete", bc.deleteById)
	r.POST("/api/v1/banner/info/id", bc.infoById)
	r.POST("/api/v1/banner/list/date", bc.lisitValidBannerByUnixDate)

	return nil
}

func (bc *BannerController) insert(ctx *gin.Context) {
	var (
		banner struct {
			Name      string    `json:"name"`
			ImagePath string    `json:"imagepath"`
			Event     string    `json:"event"`
			StartDate time.Time `json:"startDate"`
			EndDate   time.Time `json:"endDate"`
		}
	)

	err := ctx.ShouldBind(&banner)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	id, err := mysql.InsertBanner(bc.db, &banner.Name, &banner.ImagePath, &banner.Event, &banner.StartDate, &banner.EndDate)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   id,
	})
}

func (bc *BannerController) lisitValidBannerByUnixDate(ctx *gin.Context) {
	var (
		banner struct {
			Unixtime int64 `json:"unixtime"`
		}
	)

	err := ctx.ShouldBind(&banner)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	banners, err := mysql.LisitValidBannerByUnixDate(bc.db, banner.Unixtime)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   banners,
	})
}

func (bc *BannerController) infoById(ctx *gin.Context) {
	var (
		banner struct {
			Id int `json:"id"`
		}
	)

	err := ctx.ShouldBind(&banner)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	ban, err := mysql.InfoById(bc.db, banner.Id)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   ban,
	})
}

func (bc *BannerController) deleteById(ctx *gin.Context) {
	var (
		banner struct {
			Id int `json:"id"`
		}
	)

	err := ctx.ShouldBind(&banner)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"statsu": http.StatusBadRequest})
		return
	}

	err = mysql.DeleteById(bc.db, banner.Id)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})

}
