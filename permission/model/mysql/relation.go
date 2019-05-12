package mysql

import (
	"github.com/Mictrlan/Miuer/admin/model/mysql"
	"database/sql"
	"time"
)

type (
	relationData struct {
		AdminID uint32
		RoleID  uint32
	}
)

const (
	mysqlRelationCreateTable = iota
	mysqlRelationInsert
	mysqlRelationDelete
	mysqlRelationRoleMap
)

var (
	relationSqlString = []string{
		`CREATE TABLE IF NOT EXISTS admin.relation (
			admin_id 	BIGINT UNSIGNED NOT NULL,
			role_id		INT UNSIGNED NOT NULL,
			created_at 	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (admin_id,role_id)
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO admin.relation(admin_id,role_id,created_at) VALUES (?,?,?)`,
		`DELETE FROM admin.relation WHERE admin_id = ? AND role_id = ? LIMIT 1`,
		`SELECT relation.role_id FROM admin.relation, admin.role WHERE relation.admin_id = ? AND role.active = true AND relation.role_id = role.id LOCK IN SHARE MODE`,
	}
)

// CreateRoleTable create role table.
func CreateRelationTable(db *sql.DB) error {
	_, err := db.Exec(relationSqlString[mysqlRelationCreateTable])
	return err
}

// AddRole add a role to admin
func AddRelation(db *sql.DB, aid, rid uint32) error {
	adminIsActive, err := mysql.IsActive(db, aid)
	if err != nil {
		return err
	}
	if !adminIsActive {
		return errAdminInactive
	}

	role, err := GetRoleByID(db, rid)
	if err != nil {
		return err
	}

	if !role.Active {
		return errRoleInactive
	}

	result, err := db.Exec(relationSqlString[mysqlRelationInsert], aid, rid, time.Now())

	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

// RemoveRole remove role from admin.
func RemoveRelation(db *sql.DB, aid, rid uint32) error {
	adminIsActive, err := mysql.IsActive(db, aid)
	if err != nil {
		return err
	}

	if !adminIsActive {
		return errAdminInactive
	}
	_, err = db.Exec(relationSqlString[mysqlRelationDelete], aid, rid)

	return err
}

// AssociatedRoleMap list all the roles of the specified admin and the return form is map.
func AssociatedRoleMap(db *sql.DB, aid uint32) (*map[uint32]bool, error) {
	var (
		roleID uint32
		result = make(map[uint32]bool)
	)

	adminIsActive, err := mysql.IsActive(db, aid)
	if err != nil {
		return nil, err
	}

	if !adminIsActive {
		return nil, errAdminInactive
	}

	rows, err := db.Query(roleSqlString[mysqlRelationRoleMap], aid)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&roleID); err != nil {
			return nil, err
		}
		result[roleID] = true
	}
	return &result, nil
}

// AssociatedRoleList list all the roles of the specified admin and the return form is slice.
func AssociatedRoleList(db *sql.DB, aid uint32) ([]*relationData, error) {
	var (
		roleID uint32
		r      *relationData
		result []*relationData
	)
	adminIsActive, err := mysql.IsActive(db, aid)
	if err != nil {
		return nil, err
	}

	if !adminIsActive {
		return nil, errAdminInactive
	}

	rows, err := db.Query(roleSqlString[mysqlRelationRoleMap], aid)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&roleID); err != nil {
			return nil, err
		}
		r = &relationData{
			AdminID: aid,
			RoleID:  roleID,
		}
		result = append(result, r)
	}
	return result, nil

}
