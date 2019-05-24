package gin

import (
	"database/sql"
	"errors"
	"log"

	"net/http"

	"github.com/Mictrlan/Miuer/admin/model/mysql"

	"github.com/gin-gonic/gin"
)

// AdminController -
type AdminController struct {
	db *sql.DB
}

// New create new AdminController
func New(db *sql.DB) *AdminController {
	return &AdminController{
		db: db,
	}
}

var (
	errServerNotExists = errors.New("[RegisterRouter]: server is nil")
	errPwdRepeat       = errors.New("the new password can't be the same as the old password")
	errPwdDisagree     = errors.New("the new password and confirming password disagree")
)

// RegisterRouter register admin router
func (ac *AdminController) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal(errServerNotExists)
	}

	err := mysql.CreateDataBase(ac.db)
	if err != nil {
		log.Fatal(err)
	}

	err = mysql.CreateTable(ac.db)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/api/v1/admin/create", ac.create)
	r.POST("/api/v1/admin/modifyEmail", ac.modifyEmail)
	r.POST("/api/v1/admin/modifyMobile", ac.modifyMobile)
	r.POST("/api/v1/admin/modifyPwd", ac.modifyPwd)
	r.POST("/api/v1/admin/modifyActive", ac.modifyActive)

}

// Create create staff information
func (ac *AdminController) create(ctx *gin.Context) {
	var admin struct {
		Name   string `json:"name"      binding:"required,alphanum,min=2,max=30"`
		Pwd    string `json:"pwd"       binding:"printascii,min=6,max=30"`
		Mobile string `json:"mobile"    binding:"required,numeric,len=11"`
		Email  string `json:"email"     binding:"required,email"`

		PwdConf string
	}

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	if admin.Pwd != admin.PwdConf {
		ctx.Error(errPwdDisagree)
		ctx.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict})
		return
	}

	err = mysql.Create(ac.db, admin.Name, admin.Pwd, admin.Mobile, admin.Email)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated})
}

// Login user login
func (ac *AdminController) Login(ctx *gin.Context) (uint32, error) {
	var (
		admin struct {
			Name string `json:"name" binding:"required,alphanum,min=2,max=30"`
			Pwd  string `json:"pwd" binding:"printascii,min=6,max=30"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		return 0, err
	}

	ID, err := mysql.Login(ac.db, admin.Name, admin.Pwd)
	if err != nil {
		return 0, err
	}

	return ID, nil
}

// Email modify email
func (ac *AdminController) modifyEmail(ctx *gin.Context) {
	var (
		admin struct {
			ID    uint32 `json:"id"     binding:"required"`
			Email string `json:"email"  binding:"required,email"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		ctx.Error(err)
		return
	}

	err = mysql.ModifyEmail(ac.db, admin.ID, admin.Email)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (ac *AdminController) modifyMobile(ctx *gin.Context) {
	var (
		admin struct {
			ID     uint32 `json:"id" binding:"required"`
			Mobile string `json:"mobile"    binding:"required,numeric,len=11"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyMobile(ac.db, &admin.ID, &admin.Mobile)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (ac *AdminController) modifyPwd(ctx *gin.Context) {
	var (
		admin struct {
			ID      uint32 `json:"id"           binding:"required"`
			Pwd     string `json:"pwd"          binding:"printascii,min=6,max=30"`
			NewPwd  string `json:"newpwd"       binding:"printascii,min=6,max=30"`
			Confirm string `json:"confirm"      binding:"printascii,min=6,max=30"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	if admin.NewPwd == admin.Pwd {
		ctx.Error(errPwdRepeat)
		ctx.JSON(http.StatusExpectationFailed, gin.H{"status": http.StatusExpectationFailed})
		return
	}

	if admin.NewPwd != admin.Confirm {
		ctx.Error(errPwdDisagree)
		ctx.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict})
		return
	}

	err = mysql.ModifyPwd(ac.db, admin.ID, admin.Pwd, admin.NewPwd)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (ac *AdminController) modifyActive(ctx *gin.Context) {
	var (
		admin struct {
			ID     uint32 `json:"id" binding:"required"`
			Active bool   `json:"active"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.ModifyActive(ac.db, &admin.ID, admin.Active)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
