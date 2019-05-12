// need to be further perfected

package gin

import (
	mysql "github.com/Mictrlan/Miuer/order/model/mysql"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// OrderController -
type OrderController struct {
	db             *sql.DB
	orderTable     string
	itemTable      string
	closedIntercal int
}


// New - 
func New(db *sql.DB, orderTable, itemTable string) *OrderController {
	return &OrderController{
		db:         db,
		orderTable: orderTable,
		itemTable:  itemTable,
	}
}

// Register register router and create tables
func (odc *OrderController) Register(r gin.IRouter) error {

	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	err := mysql.CreateOrderTable(odc.db, odc.orderTable)
	if err != nil {
		return err
	}

	err = mysql.CreateItemTabke(odc.db, odc.itemTable)
	if err != nil {
		return err
	}

	r.POST("/api/v1/order/create", odc.insert)
	r.POST("/api/v1/order/info", odc.orderInfoByOrderID)
	r.POST("/api/v1/order/user", odc.lisitOrderByUserIDAndStatus)
	r.POST("/api/v1/order/id", odc.orderIDByOrderCode)

	return nil
}

func (odc *OrderController) insert(ctx *gin.Context) {
	var (
		req struct {
			UserID     uint64 `json:"userid"`
			AddressID  string `json:"addressid"`
			TotalPrice uint32 `json:"totalprice"`
			Promotion  string `json:"promotion"`
			Freight    uint32 `json:"freight"`

			Items []mysql.Item `json:"items"`
		}
		
		rep struct {
			ordercode string
			orderid   uint32
		}
	)

	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
		return
	}

	promotion, err := strconv.ParseBool(req.Promotion)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	times := time.Now()

	rep.ordercode = strconv.Itoa(times.Year()) + strconv.Itoa(int(times.Month())) + strconv.Itoa(times.Day()) + strconv.Itoa(times.Hour()) + strconv.Itoa(times.Minute()) + strconv.Itoa(times.Second()) + strconv.Itoa(int(req.UserID))
	order := mysql.Order{
		OrderCode:  rep.ordercode,
		UserID:     req.UserID,
		AddressID:  req.AddressID,
		TotalPrice: req.TotalPrice,
		Promotion:  promotion,
		Freight:    req.Freight,
		Created:    times,
	}

	rep.orderid, err = mysql.Insert(odc.db, order, odc.orderTable, odc.itemTable, req.Items, odc.closedIntercal)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    http.StatusOK,
		"orderid":   rep.orderid,
		"ordercode": rep.ordercode,
	})
}

func (odc *OrderController) orderIDByOrderCode(ctx *gin.Context) {
	var req struct {
		Ordercode string `json:"ordercode"`
	}

	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	id, err := mysql.OrderIDByOrderCode(odc.db, odc.orderTable, req.Ordercode)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"id":     id,
	})
}

func (odc *OrderController) orderInfoByOrderID(ctx *gin.Context) {
	var req struct {
		OrderID uint32 `json:"orderid"`
	}

	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	rep, err := mysql.OrderInfoByorderKey(odc.db, odc.orderTable, odc.itemTable, req.OrderID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"order":  rep.Order,
		"ite":    rep.Ite,
	})
}

func (odc *OrderController) lisitOrderByUserIDAndStatus(ctx *gin.Context) {
	var req struct {
		Userid uint64 `json:"userid"`
		Status uint8  `json:"status"`
	}

	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	orders, err := mysql.ListOrderByUserID(odc.db, odc.orderTable, odc.itemTable, req.Userid, req.Status)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": http.StatusPreconditionFailed})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"orders": orders,
	})

}
