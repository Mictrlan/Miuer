package gin

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Mictrlan/Miuer/banner/model/mysql"

	"github.com/gin-gonic/gin"
)

var errServerNotExists = errors.New("[RegisterRouter]: server is nil")

// BannerController -
type BannerController struct {
	db *sql.DB
}

// New create new BannerController
func New(db *sql.DB) *BannerController {
	return &BannerController{
		db: db,
	}
}

// Register register banner router
func (bc *BannerController) Register(r gin.IRouter) {
	if r == nil {
		log.Fatal(errServerNotExists)
	}

	if err := mysql.CreateDB(bc.db); err != nil {
		log.Fatal(err)
	}

	if err := mysql.CreateTable(bc.db); err != nil {
		log.Fatal(err)
	}

	r.POST("/api/v1/banner/create", bc.insert)
	r.POST("/api/v1/banner/delete", bc.deleteByID)
	r.POST("/api/v1/banner/info/id", bc.infoByID)
	r.POST("/api/v1/banner/list/date", bc.lisitValidBannerByUnixDate)

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

func (bc *BannerController) infoByID(ctx *gin.Context) {
	var (
		banner struct {
			ID int `json:"id"`
		}
	)

	err := ctx.ShouldBind(&banner)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	ban, err := mysql.InfoByID(bc.db, banner.ID)
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

func (bc *BannerController) deleteByID(ctx *gin.Context) {
	var (
		banner struct {
			ID int `json:"id"`
		}
	)

	err := ctx.ShouldBind(&banner)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"statsu": http.StatusBadRequest})
		return
	}

	err = mysql.DeleteByID(bc.db, banner.ID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
