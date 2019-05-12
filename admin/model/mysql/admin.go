package mysql

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	mysqlAdminCreateDatabase = iota
	mysqlUserCreateTable
	mysqlUserInsert
	mysqlUserLogin
	mysqlUserModifyEmail
	mysqlUserModifyMobile
	mysqlUserGetPwd
	mysqlUserModifyPwd
	mysqlUserModifyActive
	mysqlUserGetIsActive
)

var (
	errInvalidMysql = errors.New("affected 0 rows")
	errLoginFailed  = errors.New("invalid name or password")

	adminSQLString = []string{
		// 考虑数据库不存在时不能调用的情况（后续修改）
		`CREATE DATABASE IF NOT EXISTS admin`,
		`CREATE TABLE IF NOT EXISTS admin.user(
			id				BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			name			VARCHAR(512) UNIQUE NOT NULL DEFAULT ' ',
			pwd				VARCHAR(512) NOT NULL DEFAULT ' ',
			mobile			VARCHAR(32) UNIQUE NOT NULL,
			email			VARCHAR(128) UNIQUE DEFAULT NULL,  
			active			BOOLEAN	DEFAULT TRUE,
			created_at		DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY(id)  
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO admin.user(name,pwd,mobile,email)VALUES(?,?,?,?)`,
		`SELECT id,pwd FROM admin.user WHERE name = ? AND active = true LOCK IN  SHARE MODE`,
		`UPDATE admin.user SET email = ? WHERE id = ? LIMIT 1 `,
		`UPDATE admin.user SET mobile = ? WHERE id = ? LIMIT 1`,
		`SELECT pwd FROM admin.user WHERE id = ? AND active = true LOCK IN SHARE MODE`,
		`UPDATE admin.user SET pwd = ? WHERE id = ? LIMIT 1 `,
		`UPDATE admin.user SET active = ? WHERE id = ? LIMIT 1`,
		`SELECT active FROM admin.user WHERE id = ? LOCK IN SHARE MODE`,
	}
)

// CreateDataBase - 
func CreateDataBase(db *sql.DB) error {
	_, err := db.Exec(adminSQLString[mysqlAdminCreateDatabase])
	return err
}

// CreateTable - 
func CreateTable(db *sql.DB) error {
	_, err := db.Exec(adminSQLString[mysqlUserCreateTable])
	return err
}

// Create - 
func Create(db *sql.DB, name, pwd, mobile, email *string) error {

	hash, err := SaltHashGenerate(pwd)
	if err != nil {
		return err
	}

	result, err := db.Exec(adminSQLString[mysqlUserInsert], name, hash, mobile, email)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

// Login -
func Login(db *sql.DB, name, pwd *string) (uint32, error) {
	var (
		id       uint32
		password string
	)

	err := db.QueryRow(adminSQLString[mysqlUserLogin], name).Scan(&id, &password)
	if err != nil {
		return 0, err
	}

	if !SaltHashCompare([]byte(password), pwd) {
		return 0, errLoginFailed
	}

	return id, nil
}

// ModifyEmail - 
func ModifyEmail(db *sql.DB, id *uint32, email *string) error {
	result, err := db.Exec(adminSQLString[mysqlUserModifyEmail], email, id)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

// ModifyMobile - 
func ModifyMobile(db *sql.DB, id *uint32, mobile *string) error {
	result, err := db.Exec(adminSQLString[mysqlUserModifyMobile], mobile, id)
	if err != nil {

		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {

		return errInvalidMysql
	}

	return nil
}

// ModifyPwd -
func ModifyPwd(db *sql.DB, id *uint32, pwd, pwdNew *string) error {
	var (
		password string
	)

	err := db.QueryRow(adminSQLString[mysqlUserGetPwd], id).Scan(&password)
	if err != nil {
		return err
	}

	if !SaltHashCompare([]byte(password), pwd) {
		return errLoginFailed
	}

	hash, err := SaltHashGenerate(pwdNew)
	if err != nil {
		return err
	}

	_, err = db.Exec(adminSQLString[mysqlUserModifyPwd], hash, id)
	return err
}

// ModifyActive - 
func ModifyActive(db *sql.DB, id *uint32, active bool) error {
	result, err := db.Exec(adminSQLString[mysqlUserModifyActive], active, id)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

// IsActive - 
func IsActive(db *sql.DB, id uint32) (bool, error) {
	var (
		isActive bool
	)

	err := db.QueryRow(adminSQLString[mysqlUserGetIsActive], id).Scan(&isActive)
	return isActive, err
}

// SaltHashGenerate - 
func SaltHashGenerate(password *string) (string, error) {
	hex := []byte(*password)

	hashedPassword, err := bcrypt.GenerateFromPassword(hex, 10)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// SaltHashCompare -
func SaltHashCompare(digest []byte, password *string) bool {
	hex := []byte(*password)
	
	if err := bcrypt.CompareHashAndPassword(digest, hex); err == nil {
		return true
	}

	return false
}
