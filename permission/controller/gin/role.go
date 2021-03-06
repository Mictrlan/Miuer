package gin

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Mictrlan/Miuer/permission/model/mysql"

	"github.com/gin-gonic/gin"
)

// PermissionController -
type PermissionController struct {
	db *sql.DB
}

// New -
func New(db *sql.DB) *PermissionController {
	return &PermissionController{
		db: db,
	}
}

// Register -
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

	r.POST("/api/v1/permission/addrole", pc.createRole)
	r.POST("/api/v1/permission/modifyrole", pc.modifyRoleByID)
	r.POST("/api/v1/permission/activerole", pc.modifyRoleActiveByID)
	r.POST("/api/v1/permission/getrole", pc.getRoleList)
	r.POST("/api/v1/permission/getidrole", pc.getRoleByID)

	r.POST("/api/v1/permission/addurl", pc.addURLPermission)
	r.POST("/api/v1/permission/removeurl", pc.removeURLPermission)
	r.POST("/api/v1/permission/urlgetrole", pc.URLPermissions)
	r.POST("/api/v1/permission/getpermission", pc.permissions)

	r.POST("/api/v1/permission/addrelation", pc.addRelation)
	r.POST("/api/v1/permission/removerelation", pc.removeRelation)

}

func (pc *PermissionController) createRole(ctx *gin.Context) {

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

func (pc *PermissionController) modifyRoleByID(ctx *gin.Context) {
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

func (pc *PermissionController) modifyRoleActiveByID(ctx *gin.Context) {
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

func (pc *PermissionController) getRoleList(ctx *gin.Context) {
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

func (pc *PermissionController) getRoleByID(ctx *gin.Context) {
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
