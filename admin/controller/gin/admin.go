package gin

import (
	mysql "github.com/Mictrlan/Miuer/admin/model/mysql"
	"database/sql"
	"errors"
	"log"

	"net/http"

	"github.com/gin-gonic/gin"
)

// Controller - 
type Controller struct {
	db *sql.DB
}

// New - 
func New(db *sql.DB) *Controller {
	return &Controller{
		db: db,
	}
}

var (
	errPwdRepeat   = errors.New("the new password can't be the same as the old password")
	errPwdDisagree = errors.New("the new password and confirming password disagree")
)

// RegisterRouter - 
func (c *Controller) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[RegisterRouter]: server is nil")
	}

	err := mysql.CreateDataBase(c.db)
	if err != nil {
		log.Fatal(err)
	}

	err = mysql.CreateTable(c.db)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/api/v1/admin/create", c.create)
	r.POST("/api/v1/admin/modifyEmail", c.modifyEmail)
	r.POST("/api/v1/admin/modifyMobile", c.modifyMobile)
	r.POST("/api/v1/admin/modifyPwd", c.modifyPwd)

}

// Create create staff information
func (c *Controller) create(ctx *gin.Context) {
	var admin struct {
		Name   string `json:"name"  binding:"required,alphanum,min=2,max=30"`
		Pwd    string `json:"pwd"       binding:"printascii,min=6,max=30"`
		Mobile string `json:"mobile"    binding:"required,numeric,len=11"`
		Email  string `json:"email"     binding:"required,email"`
	}

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.Create(c.db, &admin.Name, &admin.Pwd, &admin.Mobile, &admin.Email)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated})

}

// Login user login
func (c *Controller) Login(ctx *gin.Context) (uint32, error) {
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

	ID, err := mysql.Login(c.db, &admin.Name, &admin.Pwd)
	if err != nil {
		return 0, err
	}

	return ID, nil

}

// Email modify email
func (c *Controller) modifyEmail(ctx *gin.Context) {
	var (
		admin struct {
			ID    uint32 `json:"id" binding:"required"`
			Email string `json:"email" binding:"required,email"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		ctx.Error(err)
		return
	}

	err = mysql.ModifyEmail(c.db, &admin.ID, &admin.Email)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})

}

func (c *Controller) modifyMobile(ctx *gin.Context) {
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

	err = mysql.ModifyMobile(c.db, &admin.ID, &admin.Mobile)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) modifyPwd(ctx *gin.Context) {
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

	err = mysql.ModifyPwd(c.db, &admin.ID, &admin.Pwd, &admin.NewPwd)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) modifyActive(ctx *gin.Context) {
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

	err = mysql.ModifyActive(c.db, &admin.ID, admin.Active)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) isactive(ctx *gin.Context) {
	var (
		admin struct {
			ID uint32 `json:"id" binding:"required"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	_, err = mysql.IsActive(c.db, admin.ID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
