package mysql

import (
	"database/sql"
	"errors"
)

type (
	role struct {
		ID       uint32
		Name     string
		Intro    string
		Active   bool
		CreateAt string
	}
)

const (
	mysqlRoleCreateTable = iota
	mysqlRoleInsert
	mysqlRoleModifyByID
	mysqlRoleModifyActiveByID
	mysqlRoleGetList
	mysqlRoleGetById
)

var (
	errInvalidMysql  = errors.New("affected 0 rows")
	errAdminInactive = errors.New("the admin is not activated")
	errRoleInactive  = errors.New("the role is not activated")

	roleSqlString = []string{
		`CREATE TABLE IF NOT EXISTS admin.role (
			id	 INT UNSIGNED NOT NULL AUTO_INCREMENT,
			name	VARCHAR(512) UNIQUE NOT NULL DEFAULT ' ',
			intro	 VARCHAR(512) NOT NULL DEFAULT ' ',
			active	BOOLEAN DEFAULT TRUE,
			create_at	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO admin.role(name,intro,active)VALUES (?,?,?)`,
		`UPDATE admin.role SET name = ?, intro = ? WHERE id = ? LIMIT 1`,
		`UPDATE admin.role SET active = ? WHERE id = ? LIMIT 1`,
		`SELECT * FROM admin.role LOCK IN SHARE MODE`,
		`SELECT * FROM admin.role WHERE id = ? AND active = true LOCK IN SHARE MODE`,
	}
)

// CreateRoleTable create role table
func CreateRoleTable(db *sql.DB) error {
	_, err := db.Exec(roleSqlString[mysqlRoleCreateTable])
	return err
}

// InsertRole insert a new line role information
func InsertRole(db *sql.DB, name, intro string) error {
	result, err := db.Exec(roleSqlString[mysqlRoleInsert], name, intro, true)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

// ModifyRoleById modify role information by id
func ModifyRoleByID(db *sql.DB, id uint32, name, intro string) error {
	_, err := db.Exec(roleSqlString[mysqlRoleModifyByID], name, intro, id)
	return err
}

// ModifyRoleActiveById modify role active by id
func ModifyRoleActiveByID(db *sql.DB, id uint32, active bool) error {
	_, err := db.Exec(roleSqlString[mysqlRoleModifyActiveByID], active, id)
	return err
}

// GetRoleList get all role information
func GetRoleList(db *sql.DB) (*[]*role, error) {
	var (
		id       uint32
		name     string
		intro    string
		active   bool
		createAt string

		roles []*role
	)

	rows, err := db.Query(roleSqlString[mysqlRoleGetList])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&id, &name, &intro, &active, &createAt); err != nil {
			return nil, err
		}

		r := &role{
			ID:       id,
			Name:     name,
			Intro:    intro,
			Active:   active,
			CreateAt: createAt,
		}

		roles = append(roles, r)
	}

	return &roles, nil
}

// GetRoleById get role information by id
func GetRoleByID(db *sql.DB, id uint32) (*role, error) {
	var roler role

	err := db.QueryRow(roleSqlString[mysqlRoleGetById], id).Scan(&roler.ID, &roler.Name, &roler.Intro, &roler.Active, &roler.CreateAt)
	return &roler, err
}
