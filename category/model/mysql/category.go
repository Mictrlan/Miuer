package mysql

import (
	"database/sql"
	"errors"
	"fmt"
)

type Category struct {
	CategoryId uint
	ParentId   uint
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
	mysqlCategoryListChirdByParentId
)

var (
	errInvaildInsert         = errors.New("insert comment: insert affected 0 rows")
	errInvalidChangeCategory = errors.New("change status: affected 0 rows")

	categorySqlString = []string{
		`CREATE DATABASE IF NOT EXISTS %s`,
		`CREATE TABLE IF NOT EXISTS %s.%s (
			categoryId INT(11) NOT NULL AUTO_INCREMENT COMMENT '类别id',
				parentId INT(11) DEFAULT NULL  COMMENT '父类别id',
				name VARCHAR(50) DEFAULT NULL COMMENT '类别名称',
				status TINYINT(1) DEFAULT '1' COMMENT '状态1-在售，2-废弃',
				createTime DATETIME DEFAULT current_timestamp COMMENT '创建时间',
				PRIMARY KEY (categoryId),INDEX(parentId)
				)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4`,
		`INSERT INTO %s.%s (parentId,name) VALUES (?,?)`,
		`UPDATE %s.%s SET status = ? WHERE categoryId = ? LIMIT 1`,
		`UPDATE %s.%s SET name = ? WHERE categoryId = ? LIMIT 1`,
		`SELECT * FROM %s.%s WHERE parentId = ?`,
	}
)

func CreateDB(db *sql.DB, dBName string) error {
	sql := fmt.Sprintf(categorySqlString[mysqlCategoryCreateDatabase], dBName)
	_, err := db.Exec(sql)
	return err
}

func CreateTable(db *sql.DB, dBName, tableName string) error {
	sql := fmt.Sprintf(categorySqlString[mysqlCategoryCreateTable], dBName, tableName)
	_, err := db.Exec(sql)
	return err
}

func InsertCategory(db *sql.DB, dBName, tableName string, parentId uint, name string) (uint, error) {
	sql := fmt.Sprintf(categorySqlString[mysqlCategoryInsert], dBName, tableName)
	result, err := db.Exec(sql, parentId, name)
	if err != nil {
		return 0, err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errInvaildInsert
	}

	categoryId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(categoryId), nil
}

func ChangeCategoryStatus(db *sql.DB, dBName, tableName string, status int8, category uint) error {
	sql := fmt.Sprintf(categorySqlString[mysqlCategoryChangeStatus], dBName, tableName)
	result, err := db.Exec(sql, status, category)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errInvalidChangeCategory
	}

	return nil
}

func ChangeCategoryName(db *sql.DB, dBName, tableName string, name string, category uint) error {
	sql := fmt.Sprintf(categorySqlString[mysqlCategoryChangeName], dBName, tableName)
	result, err := db.Exec(sql, name, category)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errInvalidChangeCategory
	}

	return nil
}

func LisitChirldrenByParentId(db *sql.DB, dBName, tableName string, parentId uint) ([]*Category, error) {
	var (
		categoryId uint
		name       string
		status     int8
		creatTime  string

		categorys []*Category
	)

	sql := fmt.Sprintf(categorySqlString[mysqlCategoryListChirdByParentId], dBName, tableName)
	rows, err := db.Query(sql, parentId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&categoryId, &parentId, &name, &status, &creatTime); err != nil {
			return nil, err
		}

		cgy := &Category{
			CategoryId: categoryId,
			ParentId:   parentId,
			Name:       name,
			Status:     status,
			CreateTime: creatTime,
		}

		categorys = append(categorys, cgy)
	}

	return categorys, nil

}
