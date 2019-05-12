// 数据表不同表下函数通过将函数变为方法来调用

package mysql

import (
	"database/sql"
)

const (
	mysqlPermissionCreateTable = iota
	mysqlPermissionInstert
	mysqlPermissionDelete
	mysqlPermissonGetRole
	mysqlPermissonGetAll
)

// Permission -
type (
	Permission struct {
		URL       string
		RoleID    uint32
		CreatedAt string
	}
)

var (
	permissionSQLString = []string{
		`CREATE TABLE IF NOT EXISTS admin.permission (
			url				VARCHAR(512) NOT NULL DEFAULT ' ',
			role_id			MEDIUMINT UNSIGNED NOT NULL,
			created_at 		DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (url,role_id)
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO admin.permission(url,role_id) VALUES (?,?)`,
		`DELETE FROM admin.permission WHERE role_id = ? AND url = ? LIMIT 1`,
		`SELECT permission.role_id FROM admin.permission, admin.role WHERE permission.url = ? AND role.active = true AND permission.role_id = role.id LOCK IN SHARE MODE`, // 同时满足全部条件才算成功
		`SELECT * FROM admin.permission LOCK IN SHARE MODE`,
	}
)

// CreatePermissionTable create permission table.
func CreatePermissionTable(db *sql.DB) error {
	_, err := db.Exec(permissionSQLString[mysqlPermissionCreateTable])
	return err
}

// AddURLPermission create an associated record of the specified URL and role.
// 通过 roleid 建立与 url 的联系
func AddURLPermission(db *sql.DB, rid uint32, url string) error {
	role, err := GetRoleByID(db, rid)
	if err != nil {
		return err
	}

	if !role.Active {
		return errRoleInactive
	}

	_, err = db.Exec(permissionSQLString[mysqlPermissionInstert], url, rid)
	return err
}

// RemoveURLPermission remove the associated records of the specified URL and role.
func RemoveURLPermission(db *sql.DB, rid uint32, url string) error {
	role, err := GetRoleByID(db, rid)
	if err != nil {
		return err
	}

	if !role.Active {
		return errRoleInactive
	}

	_, err = db.Exec(permissionSQLString[mysqlPermissionDelete], rid, url)
	return err
}

// URLPermissions lists all the roles of the specified URL.
func URLPermissions(db *sql.DB, url string) (map[uint32]bool, error) {
	var (
		roleID uint32
		result = make(map[uint32]bool)
	)

	rows, err := db.Query(permissionSQLString[mysqlPermissonGetRole], url)
	if err != nil {
		return nil, err
	}
	
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&roleID); err != nil {
			return nil, err
		}
		result[roleID] = true
	}
	
	return result, nil
}

// Permissions lists all the roles.
func Permissions(db *sql.DB) (*[]*Permission, error) {
	var (
		roleID    uint32
		url       string
		createdAt string

		result []*Permission
	)

	rows, err := db.Query(permissionSQLString[mysqlPermissonGetAll])
	if err != nil {
		return nil, err
	}
	
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&url, &roleID, &createdAt); err != nil {
			return nil, err
		}
	
		data := &Permission{
			URL:       url,
			RoleID:    roleID,
			CreatedAt: createdAt,
		}
	
		result = append(result, data)
	}

	return &result, nil
}

