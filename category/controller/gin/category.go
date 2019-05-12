package gin

import (
	"database/sql"
	"log"
	"net/http"

	mysql "github.com/Mictrlan/Miuer/category/model/mysql"

	"github.com/gin-gonic/gin"
)

type CateController struct {
	db        *sql.DB
	dBName    string
	tableName string
}

func New(db *sql.DB, dB string, table string) *CateController {
	return &CateController{
		db:        db,
		dBName:    dB,
		tableName: table,
	}
}

func (cc *CateController) Register(r gin.IRouter) error {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	// 考虑只需要注册时需要自定义 database 和 table  放在 main 还是 controller

	err := mysql.CreateDB(cc.db, cc.dBName)
	if err != nil {
		return err
	}

	err = mysql.CreateTable(cc.db, cc.dBName, cc.tableName)
	if err != nil {
		return err
	}

	r.POST("/api/v1/category/create", cc.insert)
	r.POST("/api/v1/category/modify/status", cc.changeCategoryStatus)
	r.POST("/api/v1/category/modify/name", cc.changeCategoryName)
	r.POST("/api/v1/category/children", cc.lisitChirldrenByParentId)

	return nil
}

func (cc *CateController) createDB() error {
	return mysql.CreateDB(cc.db, cc.dBName)
}

func (cc *CateController) createTable() error {
	return mysql.CreateTable(cc.db, cc.dBName, cc.tableName)
}

func (cc *CateController) insert(ctx *gin.Context) {
	var (
		category struct {
			ParentId uint   `json:"parentId"`
			Name     string `json:"name"`
		}
	)

	err := ctx.ShouldBind(&category)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	id, err := mysql.InsertCategory(cc.db, cc.dBName, cc.tableName, category.ParentId, category.Name)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"Id":     id,
	})

}

func (cc *CateController) changeCategoryStatus(ctx *gin.Context) {
	var (
		category struct {
			CategoryId uint `json:"categoryId"`
			Status     int8 `json:"status"`
		}
	)

	err := ctx.ShouldBind(&category)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ChangeCategoryStatus(cc.db, cc.dBName, cc.tableName, category.Status, category.CategoryId)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (cc *CateController) changeCategoryName(ctx *gin.Context) {
	var (
		category struct {
			CategoryId uint   `json:"categoryId"`
			Name       string `json:"name"`
		}
	)

	err := ctx.ShouldBind(&category)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ChangeCategoryName(cc.db, cc.dBName, cc.tableName, category.Name, category.CategoryId)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (cc *CateController) lisitChirldrenByParentId(ctx *gin.Context) {
	var (
		category struct {
			ParentId uint `json:"parentId"`
		}
	)

	err := ctx.ShouldBind(&category)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	categorys, err := mysql.LisitChirldrenByParentId(cc.db, cc.dBName, cc.tableName, category.ParentId)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    http.StatusOK,
		"categorys": categorys,
	})

}
