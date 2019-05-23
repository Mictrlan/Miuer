package gin

import (
	"database/sql"

	"github.com/Mictrlan/Miuer/smsservice/model/mysql"
	services "github.com/Mictrlan/Miuer/smsservice/services"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SmsController -
type SmsController struct {
	Db   *sql.DB
	Conf *services.Config
}

// New creaete a new smscontroller
func New(Db *sql.DB, conf *services.Config) *SmsController {
	return &SmsController{
		Db: Db,
		Conf: &services.Config{
			Host:           conf.Host,
			Appcode:        conf.Appcode,
			Digits:         conf.Digits,
			ResendInterval: conf.ResendInterval,
			OnCheck:        conf.OnCheck,
			DB:             conf.DB,
		},
	}
}

// Register register router
func (sc *SmsController) Register(r gin.IRouter) error {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	err := mysql.CreateDatabase(sc.Db)
	if err != nil {
		return err
	}

	err = mysql.CreateTable(sc.Db)
	if err != nil {
		return err
	}

	r.POST("/api/v1/sms/send", sc.send)
	r.POST("/api/v1/sms/check", sc.check)

	return nil
}

// Send -
func (sc *SmsController) send(ctx *gin.Context) {
	var (
		req struct {
			Mobile string `json:"mobile"`
			Sign   string `json:"sign"`
		}
	)

	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = services.Send(sc.Db, req.Mobile, req.Sign, sc.Conf)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})

}

// Check -
func (sc *SmsController) check(ctx *gin.Context) {
	var (
		req struct {
			Code string `json:"code"`
			Sign string `json:"sign"`
		}
	)

	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	var resp struct {
		sign   string
		mobile string
	}

	resp.sign = req.Sign
	resp.mobile, _ = mysql.GetMobileBySign(sc.Db, req.Sign)

	err = services.Check(req.Code, req.Sign, sc.Conf, sc.Db)
	if err != nil {
		sc.Conf.OnCheck.OnVerifyFailed(resp.sign, resp.mobile)
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	sc.Conf.OnCheck.OnVerifySucceed(resp.sign, resp.mobile)
	ctx.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": resp,
	})
}
