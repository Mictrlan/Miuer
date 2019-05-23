package mysql

import (
	"database/sql"
	"errors"
	"time"
)

const (
	mysqlFileCreateTable = iota
	mysqlFileInsert
	mysqlFileGetPathByMD5
)

// ErrNoRows -
var (
	ErrNoRows = errors.New("there is no such data in database")

	UploadSQLString = []string{
		`CREATE TABLE IF NOT EXISTS upload.files (
			user_id     INTEGER UNSIGNED NOT NULL,
			md5         VARCHAR(512) DEFAULT ' ',
			path        VARCHAR(512) DEFAULT ' ',
			created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (md5)
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO upload.files(user_id,md5,path,created_at) VALUES (?,?,?,?)`,
		`SELECT path FROM upload.files WHERE md5 = ? LOCK IN SHARE MODE`,
	}
)

// CreateTable create files table
func CreateTable(db *sql.DB) error {
	_, err := db.Exec(UploadSQLString[mysqlFileCreateTable])
	return err
}

// Insert  add file info to table
func Insert(db *sql.DB, userID uint32, path, md5 string) error {
	result, err := db.Exec(UploadSQLString[mysqlFileInsert], userID, path, md5, time.Now())
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return err
	}

	return nil
}

// QueryPathByMD5 query path by md5
func QueryPathByMD5(db *sql.DB, md5 string) (string, error) {
	var (
		path string
	)

	err := db.QueryRow(UploadSQLString[mysqlFileGetPathByMD5], md5).Scan(&path)
	if err != nil {
		if err == sql.ErrNoRows {
			return path, ErrNoRows
		}
		return path, err
	}

	return path, nil
}
