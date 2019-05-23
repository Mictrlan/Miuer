package mysql

import (
	"database/sql"
	"errors"
)

// Message -
type Message struct {
	Mobile string `db:"mobile"`
	Date   int64  `db:"date"`
	Code   string `db:"code"`
	Sign   string `db:"sign"`
}

const (
	mysqlSmsCreateDatabase = iota
	mysqlSmsCreateTable
	mysqlSmsAddMessage
	mysqlSmsGetMobile
	mysqlSmsGetDate
	mysqlSmsGetCode
	mysqlSmsDeleteMessage
)

var smsSQLString = []string{
	`CREATE DATABASE IF NOT EXISTS  SMS`,
	`CREATE TABLE IF NOT EXISTS SMS.msg(
		mobile      VARCHAR(32) UNIQUE NOT NULL,
		date        INT(11) DEFAULT 0,
		code        VARCHAR(32) ,
		sign        VARCHAR(32) UNIQUE NOT NULL
	)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
	`INSERT INTO SMS.msg(mobile,date,code,sign) VALUES (?,?,?,?)`,
	`SELECT mobile FROM SMS.msg WHERE sign = ? LOCK IN SHARE MODE`,
	`SELECT date FROM SMS.msg WHERE sign = ? LOCK IN SHARE MODE`,
	`SELECT code FROM SMS.msg WHERE sign = ? LOCK IN SHARE MODE`,
	`DELETE FROM SMS.msg WHERE sign = ? LIMIT 1`,
}

// CreateDatabase -
func CreateDatabase(db *sql.DB) error {
	_, err := db.Exec(smsSQLString[mysqlSmsCreateDatabase])
	return err
}

// CreateTable -
func CreateTable(db *sql.DB) error {
	_, err := db.Exec(smsSQLString[mysqlSmsCreateTable])
	return err
}

// AddSmsMessage add  new sms
func AddSmsMessage(db *sql.DB, mobile string, date int64, code string, sign string) error {
	result, err := db.Exec(smsSQLString[mysqlSmsAddMessage], mobile, date, code, sign)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errors.New("errInvalidAdd")
	}

	return nil
}

// GetMobileBySign return user's mobile  by sign
func GetMobileBySign(db *sql.DB, sign string) (string, error) {
	var mobile string

	err := db.QueryRow(smsSQLString[mysqlSmsGetMobile], sign).Scan(&mobile)
	if err != nil {
		return "0", err
	}

	return mobile, nil
}

// GetDateBySign return sms date(uinxtime) by sign
func GetDateBySign(db *sql.DB, sign string) (int64, error) {
	var unixtime int64

	err := db.QueryRow(smsSQLString[mysqlSmsGetDate], sign).Scan(&unixtime)
	if err != nil {
		return 0, err
	}

	return unixtime, nil
}

// GetCodeBySign return sms code by sign
func GetCodeBySign(db *sql.DB, sign string) (string, error) {
	var code string

	err := db.QueryRow(smsSQLString[mysqlSmsGetCode], sign).Scan(&code)
	if err != nil {
		return "0", err
	}

	return code, nil
}

// DeleteSmsMessage remove sms by sign
func DeleteSmsMessage(db *sql.DB, sign string) error {
	_, err := db.Exec(smsSQLString[mysqlSmsDeleteMessage], sign)
	return err
}

// GetMessageBySign return msg
func GetMessageBySign(db *sql.DB, sign string) *Message {
	var msg Message

	msg.Mobile, _ = GetMobileBySign(db, sign)
	msg.Date, _ = GetDateBySign(db, sign)
	msg.Code, _ = GetCodeBySign(db, sign)
	msg.Sign = sign

	return &msg
}
