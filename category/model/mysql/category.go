package mysql

import (
	"database/sql"
	"errors"
	"fmt"
)

// Category -
type Category struct {
	CategoryID uint
	ParentID   uint
	Name       string
	Status     int8
	CreateTime string
}

const (
	mysqlCategoryCreateDatabase = iota
	mysqlCategoryCreateTable
	mysqlCategoryInsert
	mysqlCategoryChangeStatus
	mysqlCategoryChangeName
	mysqlCategoryListChirdByParentID
)

var (
	errInvaildInsert         = errors.New("insert comment: insert affected 0 rows")
	errInvalidChangeCategory = errors.New("change status: affected 0 rows")

	categorySQLString = []string{
		`CREATE DATABASE IF NOT EXISTS %s`,
		`CREATE TABLE IF NOT EXISTS %s.%s (
			categoryId 			INT(11) NOT NULL AUTO_INCREMENT COMMENT '类别id',
				parentId 		INT(11) DEFAULT NULL  COMMENT '父类别id',
				name 			VARCHAR(50) DEFAULT NULL COMMENT '类别名称',
				status 			TINYINT(1) DEFAULT '1' COMMENT '状态1-在售，2-废弃',
				createTime 		DATETIME DEFAULT current_timestamp COMMENT '创建时间',
				PRIMARY KEY (categoryId),INDEX(parentId)
				)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4`,
		`INSERT INTO %s.%s (parentId,name) VALUES (?,?)`,
		`UPDATE %s.%s SET status = ? WHERE categoryId = ? LIMIT 1`,
		`UPDATE %s.%s SET name = ? WHERE categoryId = ? LIMIT 1`,
		`SELECT * FROM %s.%s WHERE parentId = ?`,
	}
)

// CreateDB -
func CreateDB(db *sql.DB, dBName string) error {
	sql := fmt.Sprintf(categorySQLString[mysqlCategoryCreateDatabase], dBName)

	_, err := db.Exec(sql)
	return err
}

// CreateTable -
func CreateTable(db *sql.DB, dBName, tableName string) error {
	sql := fmt.Sprintf(categorySQLString[mysqlCategoryCreateTable], dBName, tableName)

	_, err := db.Exec(sql)
	return err
}

// InsertCategory - 
func InsertCategory(db *sql.DB, dBName, tableName string, parentID uint, name string) (uint, error) {
	sql := fmt.Sprintf(categorySQLString[mysqlCategoryInsert], dBName, tableName)

	result, err := db.Exec(sql, parentID, name)
	if err != nil {
		return 0, err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errInvaildInsert
	}

	categoryID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(categoryID), nil
}

// ChangeCategoryStatus -
func ChangeCategoryStatus(db *sql.DB, dBName, tableName string, status int8, category uint) error {
	sql := fmt.Sprintf(categorySQLString[mysqlCategoryChangeStatus], dBName, tableName)

	result, err := db.Exec(sql, status, category)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errInvalidChangeCategory
	}

	return nil
}

// ChangeCategoryName - 
func ChangeCategoryName(db *sql.DB, dBName, tableName string, name string, category uint) error {
	sql := fmt.Sprintf(categorySQLString[mysqlCategoryChangeName], dBName, tableName)

	result, err := db.Exec(sql, name, category)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errInvalidChangeCategory
	}

	return nil
}

// LisitChirldrenByParentID - 
func LisitChirldrenByParentID(db *sql.DB, dBName, tableName string, parentID uint) ([]*Category, error) {
	var (
		categoryID uint
		name       string
		status     int8
		creatTime  string

		categorys []*Category
	)

	sql := fmt.Sprintf(categorySQLString[mysqlCategoryListChirdByParentID], dBName, tableName)
	rows, err := db.Query(sql, parentID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&categoryID, &parentID, &name, &status, &creatTime); err != nil {
			return nil, err
		}

		cgy := &Category{
			CategoryID: categoryID,
			ParentID:   parentID,
			Name:       name,
			Status:     status,
			CreateTime: creatTime,
		}

		categorys = append(categorys, cgy)
	}

	return categorys, nil
}
