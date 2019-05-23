package gin

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/Mictrlan/Miuer/category/model/mysql"

	"github.com/gin-gonic/gin"
)

var errServerNotExists = errors.New("[RegisterRouter]: server is nil")

// CateController -
type CateController struct {
	db        *sql.DB
	dBName    string
	tableName string
}

// New create new CateController
func New(db *sql.DB, dB string, table string) *CateController {
	return &CateController{
		db:        db,
		dBName:    dB,
		tableName: table,
	}
}

// Register register category router
func (cc *CateController) Register(r gin.IRouter) {
	if r == nil {
		log.Fatal(errServerNotExists)
	}

	err := mysql.CreateDB(cc.db, cc.dBName)
	if err != nil {
		log.Fatal(err)
	}

	err = mysql.CreateTable(cc.db, cc.dBName, cc.tableName)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/api/v1/category/create", cc.insert)
	r.POST("/api/v1/category/modify/status", cc.changeCategoryStatus)
	r.POST("/api/v1/category/modify/name", cc.changeCategoryName)
	r.POST("/api/v1/category/children", cc.lisitChirldrenByParentID)

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
			ParentID uint   `json:"parentId"`
			Name     string `json:"name"`
		}
	)

	err := ctx.ShouldBind(&category)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	id, err := mysql.InsertCategory(cc.db, cc.dBName, cc.tableName, category.ParentID, category.Name)
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
			CategoryID uint `json:"categoryId"`
			Status     int8 `json:"status"`
		}
	)

	err := ctx.ShouldBind(&category)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ChangeCategoryStatus(cc.db, cc.dBName, cc.tableName, category.Status, category.CategoryID)
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
			CategoryID uint   `json:"categoryId"`
			Name       string `json:"name"`
		}
	)

	err := ctx.ShouldBind(&category)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ChangeCategoryName(cc.db, cc.dBName, cc.tableName, category.Name, category.CategoryID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (cc *CateController) lisitChirldrenByParentID(ctx *gin.Context) {
	var (
		category struct {
			ParentID uint `json:"parentId"`
		}
	)

	err := ctx.ShouldBind(&category)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	categorys, err := mysql.LisitChirldrenByParentID(cc.db, cc.dBName, cc.tableName, category.ParentID)
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
