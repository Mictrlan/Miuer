package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Order struct {
	ID         uint32
	OrderCode  string    `json:"ordercode"`
	UserID     uint64    `json:"userid"`
	ShipCode   string    `json:"shipcode"`
	AddressID  string    `json:"addressid"`
	TotalPrice uint32    `json:"totalprice"`
	PayWay     uint8     `json:"payway"`
	Promotion  bool      `json:"promotion"`
	Freight    uint32    `json:"freight"`
	Status     uint8     `json:"status"`
	Created    time.Time `json:"created"`
	Closed     time.Time `json:"closed"`
	Updated    time.Time `json:"updated"`
}

type Item struct {
	ProductId uint32 `json:"productid"`
	OrderID   uint32 `json:"orderid"`
	Count     uint32 `json:"count"`
	Price     uint32 `json:"price"`
	Discount  uint32 `json:"discount"`
}

type ItemOrder struct {
	*Order
	Ite []*Item
}

const (
	orderTable = iota
	itemTable
	orderInsert
	itemInsert
	orderIdByOrderCode
	orderByOrderID
	itemsByOrderID
	orderListByUserID
	payByOrderID
	consignByOrderID
	statusByOrderID
)

var (
	orderSqlString = []string{
		`CREATE TABLE IF NOT EXISTS Miuer.%s (
			id INT UNSIGNED NOT NULL AUTO_INCREMENT ,
			orderCode VARCHAR(50) NOT NULL,
			userID BIGINT UNSIGNED NOT NULL,
			shipCode VARCHAR(50) NOT NULL DEFAULT '100000',
			addressID VARCHAR(20) NOT NULL,
			totalPrice INT UNSIGNED NOT NULL,
			payWay TINYINT UNSIGNED DEFAULT '0',
			promotion TINYINT(1) UNSIGNED DEFAULT '0',
			freight INT UNSIGNED NOT NULL,
			status TINYINT UNSIGNED DEFAULT '0',
			created DATETIME DEFAULT NOW(),
			closed DATETIME DEFAULT '8012-12-31 00:00:00',
			updated DATETIME DEFAULT NOW(),
			PRIMARY KEY (id),
			UNIQUE KEY orderCode (orderCode) USING BTREE,
			KEY created (created),
			KEY updated (updated),
			KEY status (status), 
			KEY payWay (payWay)
		)ENGINE=InnoDB AUTO_INCREMENT = 10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='order info'`,
		`CREATE TABLE IF NOT EXISTS Miuer.%s(
			productID INT UNSIGNED NOT NULL,
			orderID VARCHAR(50) NOT NULL,
			count INT UNSIGNED NOT NULL,
			price INT UNSIGNED NOT NULL,
			discount TINYINT UNSIGNED NOT NULL,
			KEY orderID (orderID)
		)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='orderitem info'`,
		`INSERT INTO   Miuer.%s  (orderCode,userID,addressID,totalPrice,promotion,freight,closed) VALUES(?,?,?,?,?,?,?)`,
		`INSERT INTO  Miuer.%s   (productID,orderID,count,price,discount) VALUES(?,?,?,?,?)`,
		`SELECT id FROM   Miuer.%s   WHERE orderCode = ? LOCK IN SHARE MODE`,
		`SELECT * FROM   Miuer.%s  WHERE id = ? LOCK IN SHARE MODE`,
		`SELECT * FROM  Miuer.%s   WHERE orderID = ? LOCK IN SHARE MODE`,
		`SELECT * FROM  Miuer.%s   WHERE userID = ? AND status = ? LOCK IN SHARE MODE`,
		`UPDATE  Miuer.%s  SET payWay = ? , updated = ? , status = 2 WHERE id = ? LIMIT 1 `,
		`UPDATE   Miuer.%s   SET shipCode = ? , updated = ? , status = 3 WHERE id = ? LIMIT 1 `,
		`UPDATE  Miuer.%s  SET status = ? , updated = ? WHERE id = ? LIMIT 1 `,
	}
)

func CreateOrderTable(db *sql.DB, tableName string) error {
	sql := fmt.Sprintf(orderSqlString[orderTable], tableName)

	_, err := db.Exec(sql)
	fmt.Println(err)
	return err
}

func CreateItemTabke(db *sql.DB, tableName string) error {
	sql := fmt.Sprintf(orderSqlString[itemTable], tableName)
	_, err := db.Exec(sql)
	return err
}

func Insert(db *sql.DB, order Order, orderTable, itemTable string, items []Item, closedIntercal int) (uint32, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	order.Closed = order.Created.Add(time.Duration(closedIntercal * int(time.Hour)))

	sql := fmt.Sprintf(orderSqlString[orderInsert], orderTable)
	result, err := db.Exec(sql, order.OrderCode, order.UserID, order.AddressID, order.TotalPrice, order.Promotion, order.Freight, order.Closed)
	if err != nil {
		return 0, err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errors.New("[insert order] : insert order affected 0 rows")
	}

	Id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	order.ID = uint32(Id)

	for _, x := range items {
		sql := fmt.Sprintf(orderSqlString[itemInsert], itemTable)
		result, err := db.Exec(sql, x.ProductId, order.ID, x.Count, x.Price, x.Discount)
		if err != nil {
			return 0, err
		}

		if affected, _ := result.RowsAffected(); affected == 0 {
			return 0, errors.New("insert item: insert affected 0 rows")
		}
	}
	return order.ID, nil
}

func OrderIdByOrderCode(db *sql.DB, ostore, ordercode string) (uint32, error) {
	var orderid uint32

	sql := fmt.Sprintf(orderSqlString[orderIdByOrderCode], ostore)

	rows, err := db.Query(sql, ordercode)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&orderid); err != nil {
			return 0, err
		}
	}

	return orderid, nil
}
func ListOrderByUserId(db *sql.DB, ostore, istore string, userid uint64, status uint8) ([]*ItemOrder, error) {
	var ItOs []*ItemOrder

	sql1 := fmt.Sprintf(orderSqlString[orderListByUserID], ostore)
	sql2 := fmt.Sprintf(orderSqlString[itemsByOrderID], istore)

	rows, err := db.Query(sql1, userid, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ito ItemOrder
		var od Order
		if err := rows.Scan(&od.ID, &od.OrderCode, &od.UserID, &od.ShipCode, &od.AddressID, &od.TotalPrice, &od.PayWay, &od.Promotion, &od.Freight, &od.Status, &od.Created, &od.Closed, &od.Updated); err != nil {
			return nil, err
		}

		ito.Order = &od
		ito.Ite, err = ListItemByOrderId(db, sql2, ito.ID)
		if err != nil {
			return nil, err
		}
		ItOs = append(ItOs, &ito)
	}
	return ItOs, nil
}

func OrderInfoByorderKey(db *sql.DB, ostore, istore string, orderid uint32) (*ItemOrder, error) {

	sql1 := fmt.Sprintf(orderSqlString[orderByOrderID], ostore)
	sql2 := fmt.Sprintf(orderSqlString[itemsByOrderID], istore)

	order, err := SelectByOrderKey(db, sql1, sql2, orderid)
	if err != nil {
		return nil, err
	}

	return order, nil

}

func UpdatePayByOrderkey(tx *sql.Tx, ostore string, orderid uint32, payway uint8, time time.Time) (uint32, error) {
	sql := fmt.Sprintf(orderSqlString[payByOrderID], ostore)
	result, err := tx.Exec(sql, payway, time, orderid)
	if err != nil {
		return 0, err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errors.New("[change error] ; not update pay infomation for order module ")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}

func UpdateShipByOrderKey(tx *sql.Tx, ostore string, orderid uint32, shipcode string, time time.Time) (uint32, error) {
	sql := fmt.Sprintf(orderSqlString[consignByOrderID], ostore)
	result, err := tx.Exec(sql, shipcode, time, orderid)
	if err != nil {
		return 0, err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errors.New("[change error] : not update ship infomation for order module ")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}

func UpdateStatusByOrderKey(tx *sql.Tx, ostore string, orderid uint32, status uint8, time time.Time) (uint32, error) {
	sql := fmt.Sprintf(orderSqlString[statusByOrderID], ostore)
	result, err := tx.Exec(sql, status, time, orderid)
	if err != nil {
		return 0, err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errors.New("[change error] : not update status  for order module ")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint32(id), nil
}

// 可能会有错误
func CheckPromotion(tx *sql.Tx, db *sql.DB, ostore, istore string, orderid uint32) ([]*Item, error) {
	sql1 := fmt.Sprintf(orderSqlString[orderByOrderID], ostore)
	sql2 := fmt.Sprintf(orderSqlString[itemsByOrderID], istore)

	order, err := SelectByOrderKey(db, sql1, sql2, orderid)
	if err != nil {
		return nil, err
	}
	if order.Promotion {
		return order.Ite, nil
	}

	return nil, nil
}

func SelectByOrderKey(db *sql.DB, query, queryitem string, orderid uint32) (*ItemOrder, error) {
	var (
		ito ItemOrder
		od  Order
	)
	rows, err := db.Query(query, orderid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&od.OrderCode, &od.UserID, &od.ShipCode, &od.AddressID, &od.TotalPrice, &od.PayWay, &od.Promotion, &od.Freight, &od.Status, &od.Created, od.Closed, &od.Updated); err != nil {
			return nil, err
		}
	}
	ito.Order = &od
	ito.Ite, err = ListItemByOrderId(db, queryitem, orderid)
	if err != nil {
		return nil, err
	}

	return &ito, nil
}

func ListItemByOrderId(db *sql.DB, query string, orderid uint32) ([]*Item, error) {
	var (
		ProductId uint32
		OrderID   uint32
		Count     uint32
		Price     uint32
		Discount  uint32

		items []*Item
	)

	rows, err := db.Query(query, orderid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&ProductId, &OrderID, &Count, &Price, &Discount); err != nil {
			return nil, err
		}

		item := &Item{
			ProductId: ProductId,
			OrderID:   OrderID,
			Count:     Count,
			Price:     Price,
			Discount:  Discount,
		}
		items = append(items, item)
	}

	return items, nil
}
