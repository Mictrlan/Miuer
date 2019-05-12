package gin

import (
	"github.com/Mictrlan/Miuer/permission/model/mysql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (pc *PermissionController) AddRelation(ctx *gin.Context) {
	var (
		relation struct {
			AdminID uint32 `json:"admin_id" binding:"required"`
			RoleID  uint32 `json:"role_id" binding:"required"`
		}
	)

	err := ctx.ShouldBind(&relation)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.AddRelation(pc.db, relation.AdminID, relation.RoleID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (pc *PermissionController) RemoveRelation(ctx *gin.Context) {
	var (
		relation struct {
			AdminID uint32 `json:"admin_id" binding:"required"`
			RoleID  uint32 `json:"role_id" binding:"required"`
		}
	)

	err := ctx.ShouldBind(&relation)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.RemoveRelation(pc.db, relation.AdminID, relation.RoleID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
