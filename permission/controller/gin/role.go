package gin

import (
	"database/sql"
	"log"
	"net/http"

	"Miuer/permission/model/mysql"

	"github.com/gin-gonic/gin"
)

type PermissionController struct {
	db *sql.DB
}

func New(db *sql.DB) *PermissionController {
	return &PermissionController{
		db: db,
	}
}

func (pc *PermissionController) Register(r gin.IRouter) {
	if r == nil {
		log.Fatal("[RegisterRouter]: server is nil")
	}

	err := mysql.CreateRoleTable(pc.db)
	if err != nil {
		log.Fatal(err)
	}

	err = mysql.CreatePermissionTable(pc.db)
	if err != nil {
		log.Fatal(err)
	}

	err = mysql.CreateRelationTable(pc.db)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/api/v1/permission/addrole", pc.CreateRole)
	r.POST("/api/v1/permission/modifyrole", pc.ModifyRoleByID)
	r.POST("/api/v1/permission/activerole", pc.ModifyRoleActiveByID)
	r.POST("/api/v1/permission/getrole", pc.GetRoleList)
	r.POST("/api/v1/permission/getidrole", pc.GetRoleByID)

	r.POST("/api/v1/permission/addurl", pc.AddURLPermission)
	r.POST("/api/v1/permission/removeurl", pc.RemoveURLPermission)
	r.POST("/api/v1/permission/urlgetrole", pc.URLPermissions)
	r.POST("/api/v1/permission/getpermission", pc.Permissions)

	r.POST("/api/v1/permission/addrelation", pc.AddRelation)
	r.POST("/api/v1/permission/removerelation", pc.RemoveRelation)

}

func (pc *PermissionController) CreateRole(ctx *gin.Context) {

	var (
		role struct {
			Name  string `json:"name"  binding:"required,alphanum,min=2,max=64"`
			Intro string `json:"intro" binding:"required,alphanum,min=2,max=256"`
		}
	)

	err := ctx.ShouldBind(&role)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.InsertRole(pc.db, role.Name, role.Intro)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})

}

func (pc *PermissionController) ModifyRoleByID(ctx *gin.Context) {
	var (
		role struct {
			ID    uint32 `json:"id"    binding:"required"`
			Name  string `json:"name"  binding:"required,alphanum,min=2,max=64"`
			Intro string `json:"intro" binding:"required,alphanum,min=2,max=256"`
		}
	)

	err := ctx.ShouldBind(&role)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyRoleByID(pc.db, role.ID, role.Name, role.Intro)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})

}

func (pc *PermissionController) ModifyRoleActiveByID(ctx *gin.Context) {
	var (
		role struct {
			ID     uint32 `json:"id"    binding:"required"`
			Active bool   `json:"active"`
		}
	)

	err := ctx.ShouldBind(&role)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyRoleActiveByID(pc.db, role.ID, role.Active)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (pc *PermissionController) GetRoleList(ctx *gin.Context) {
	result, err := mysql.GetRoleList(pc.db)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"roles":  result,
	})
}

func (pc *PermissionController) GetRoleByID(ctx *gin.Context) {
	var (
		role struct {
			ID uint32 `json:"id" binding:"required"`
		}
	)

	err := ctx.ShouldBind(&role)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	result, err := mysql.GetRoleByID(pc.db, role.ID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"roles":  result,
	})
}
