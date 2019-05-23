package gin

import (
	"errors"
	"net/http"

	"github.com/Mictrlan/Miuer/permission/model/mysql"
	"github.com/gin-gonic/gin"
)

var (
	// UniqueURL -
	UniqueURL = "/api/v1/permission/addurl"

	errPermission = errors.New("Admin permission is wrong")
)

// CheckPermission -
func CheckPermission(pc *PermissionController, GetUID func(Context *gin.Context) (uint32, error)) func(c *gin.Context) {
	return func(ctx *gin.Context) {

		var trafficability = false
		IURL := ctx.Request.URL.Path

		uid, err := GetUID(ctx)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		ulrRole, err := mysql.URLPermissions(pc.db, UniqueURL)
		if err != nil {
			ctx.AbortWithError(http.StatusFailedDependency, err)
			return
		}

		ulrRoleID, err := mysql.URLPermissions(pc.db, IURL)
		if err != nil {
			ctx.AbortWithError(http.StatusForbidden, err)
			return
		}

		roleByAdmin, err := mysql.AssociatedRoleMap(pc.db, uid)
		if err != nil {
			ctx.AbortWithError(http.StatusConflict, err)
			return
		}

		allRole, err := mysql.GetAllRoleMap(pc.db)

		if err != nil {
			ctx.AbortWithError(http.StatusRequestedRangeNotSatisfiable, err)
			return
		}

		if len(ulrRole) == 0 || len(allRole) == 0 {
			return
		}

		for key := range ulrRoleID {
			if roleByAdmin[key] == true {
				return
			}
		}

		if !trafficability {
			ctx.AbortWithError(http.StatusFailedDependency, nil)
		}

	}

}
