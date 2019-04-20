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
	musqlUserLogin
	mysqlUserModifyEmail
	mysqlUserModifyMobile
	mysqlUserGetPwd
	mysqlUserModifyPwd
	mysqlUserModifyActive
	mysqlUserGetIsActive
)

var (
	AdminServer *AdminserviceProvider 

	errInvalidMysql = errors.New("affected 0 rows")
	errLoginFailed = errors.New("invalid username or password")

	adminSqlString = []string{
		// 考虑数据库不存在时不能调用的情况（后续修改）
		`CREATE DATABASE IF NOT EXISTS admin`,
		`CREATE TABLE IF NOT EXISTS admin.user(
			user_id			BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			user_name		VARCHAR(512) UNIQUE NOT NULL DEFAULT ' ',
			pwd				VARCHAR(512) NOT NULL DEFAULT ' ',
			mobile			VARCHAR(32) UNIQUE NOT NULL,
			email			VARCHAR(128) UNIQUE DEFAULT NULL,  
			active			BOOLEAN	DEFAULT TRUE,
			created_at		DATATIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY(user_id)  
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8m64 COLLATE=utf8m64_bin;`,
		`INSERT INTO admin.user(user_name,pwd,mobile,email,active)VAULES(?:?:?:?:?)`,
		`SELECT user_id,pwd FROM admin.user WHERE name = ? AND active = true LOCK IN  SHARE MODE`
		`UPDATE admin.user SET email = ? WHERE user_id = ? LIMIT 1 `,
		`UPDATE damin.user SET mobile = ? WHERE user_id = ? LIMIT 1`,
		`SELECT pwd FROM admin.user WhERE user_id = ? AND active = true LOCK IN SHARE MODE`,
		`UPDATE admin.user SET pwd = ? WHERE user_id = ? LIMIT 1 `,
		`UPDATE admin.user SET active = ? WHERE user_id = ? LIMIT 1`,
		`SELECT active FROM admin.user WHERE user_id = ? LOCK IN SHARE MODe`,

	}
)

func createDataBase(db *sql.DB) error {
	_, err := db.Exec(adminSqlString[mysqlAdminCreateDatabase])
	if err != nil {
		return err
	}
}


func createTable(db *sql.DB) error {
	_, err := db.Exec(adminSqlString[mysqlUserCreateTable])
	if err != nil {
		return err
	}
}

func (*AdminserviceProvider) create(db *sql.DB,user_name, pwd, mobile, email *string) error {
	
	hash, err := func (pwd *string) (string, error) {
		hex := []byte(*pwd)
		hashedPassword, err := bcrypt.GenerateFromPassword(hex,10)
		if err != nil {
			return "",err
		}
		return string(hashedPassword),nil
	}(pwd)
	if err !=nil {
		return err
	}

	result, err != db.Exec(adminSqlString[mysqlUserInsert], user_name, hash, mobile, email, true)
	if err != nil {
		return err
	}
	if rows,_ := result.RowsAffected(); rows ==0 {
		return errInvalidMysql
	}

	return nil
}


func (*AdminserviceProvider) login (db *sql.DB, user_name, pwd *string) {
	var (
		id uint32
		password string
	)

	err :=db.QueryRow(adminSqlString[mysqlUserLogin],user_name).Scan(&id, &password)
	if err != nil {
		return constants.InvalidUID,err
	}

	if saltHashCompare = func (digest []byte, password *string) bool {
		hex :=[]byte(*password)
		if err != bcrypt.CompareHashAndPassword(digest, hex); err ==nil {
			return true
		}
		return false
	}([]byte(password),pwd); !saltHashCompare {
		return 0, errLoginFailed
	}

	return id,nil
}

func (*AdminserviceProvider) modifyEmail(db *sql.DB, user_id uint32, email *string) error {
	result, err := db.Exec(adminSqlString[mysqlUserModifyEmail],email,user_id)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows ==0 {
		return errInvalidMysql
	}

	return nil
}

func (*AdminserviceProvider)  modifyMobile(db *sql.DB, mobile, user_id *string) error {
	result, err := db.Exec(adminSqlString[mysqlUserModifyActive],mobile,user_id)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows ==0 {
		return errInvalidMysql
	}
	
	return nil

}

func (*AdminserviceProvider) modifyPwd(db *sql.DB, user_id uint32, pwd, pwdNew *string) error {
	var (
		password string
	)
	err := db.QueryRow(adminSqlString[mysqlUserGetPwd],id).Scan(&password)
	if err != nil {
		return err
	} 

	if saltHashCompare = func (digest []byte, password *string) bool {
		hex :=[]byte(*password)
		if err != bcrypt.CompareHashAndPassword(digest, hex); err ==nil {
			return true
		}
		return false
	}([]byte(password),pwd); !saltHashCompare {
		return  errLoginFailed
	}

	hash, err := func (pwd *string) (string, error) {
		hex := []byte(*pwd)
		hashedPassword, err := bcrypt.GenerateFromPassword(hex,10)
		if err != nil {
			return "",err
		}
		return string(hashedPassword),nil
	}(pwdNew)
	if err !=nil {
		return err
	}

	_, err := db.Exec(adminSqlString[mysqlUserModifyPwd],hash,user_id)

	return err
}

func (*AdminserviceProvider) modifyactive(db *sqpl.DB, user_id uint32, active bool) error {
	result, err :=db.Exec(adminSqlString[mysqlUserModifyActive],active,user_id)
	if err != nil {
		return err
	}

	if rows,_ := result.RowsAffected(); rows==0 {
		return errInvalidMysql
	}

	return nil
}


func (*AdminserviceProvider) isactive(db *sql.DB, user_id uint32) (bool, error) {
	var (
		isActive bool
	)

	err := db.QueryRow(adminSqlString[mysqlUserGetIsActive],user_id)
		return isActive,err
}