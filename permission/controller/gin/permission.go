package gin

import (
	"github.com/Mictrlan/Miuer/permission/model/mysql"
	"net/http"
	"github.com/gin-gonic/gin"
)

// AddURLPermission - 
func (pc *PermissionController) AddURLPermission(ctx *gin.Context) {
	var (
		url struct {
			URL    string `json:"url"     binding:"required"`
			RoleID uint32 `json:"role_id" binding:"required"`
		}
	)

	err := ctx.ShouldBind(&url)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadGateway})
		return
	}

	err = mysql.AddURLPermission(pc.db, url.RoleID, url.URL)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

// RemoveURLPermission - 
func (pc *PermissionController) RemoveURLPermission(ctx *gin.Context) {
	var (
		url struct {
			URL    string `json:"url"  binding:"required"`
			RoleID uint32 `json:"role_id"   binding:"required"`
		}
	)

	err := ctx.ShouldBind(&url)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadGateway})
		return
	}

	err = mysql.RemoveURLPermission(pc.db, url.RoleID, url.URL)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

// URLPermissions - 
func (pc *PermissionController) URLPermissions(ctx *gin.Context) {
	var (
		url struct {
			URL string `json:"url" binding:"required"`
		}
	)

	err := ctx.ShouldBind(&url)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadGateway})
		return
	}

	result, err := mysql.URLPermissions(pc.db, url.URL)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":        http.StatusOK,
		"rulpermission": result,
	})
}

// Permissions -
func (pc *PermissionController) Permissions(ctx *gin.Context) {
	result, err := mysql.Permissions(pc.db)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":        http.StatusOK,
		"rulpermission": result,
	})
}
