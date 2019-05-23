package mysql

import (
	"database/sql"
	"errors"
	"time"
)

const (
	mysqlBannerCreateDatabase = iota
	mysqlBannerCreateTable
	mysqlBannerInsert
	mysqlBannerListDate
	mysqlBannerInfoByID
	mysqlBannerDeleteByID
)

// Banner -
type Banner struct {
	BannerID  int
	Name      string
	ImagePath string
	Event     string
	StartDate string
	EndDate   string
}

var (
	errInvalidInsert = errors.New("insert banner:insert affected 0 rows")

	bannerSQLString = []string{
		`CREATE DATABASE IF NOT EXISTS banner`,
		`CREATE TABLE IF NOT EXISTS banner.ads(
			bannerid        BIGINT  NOT NULL AUTO_INCREMENT,
			name            VARCHAR(512) UNIQUE DEFAULT ' ',
			imagepath       VARCHAR(512) UNIQUE NOT NULL,
			event           VARCHAR(512) DEFAULT ' ',
			startdate       DATETIME NOT NULL,
			enddate         DATETIME NOT NULL,
			PRIMARY KEY(bannerid)
		)ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO banner.ads(name,imagepath,event,startdate,enddate)VALUES(?,?,?,?,?)`,
		`SELECT * FROM banner.ads WHERE UNIX_TIMESTAMP(startdate) <= ? AND  UNIX_TIMESTAMP(enddate) >= ? LOCK IN SHARE MODE `,
		`SELECT * FROM banner.ads WHERE bannerid = ? LOCK IN SHARE MODE `,
		`DELETE FROM banner.ads WHERE bannerid = ? LIMIT 1`,
	}
)

// CreateDB create banner database
func CreateDB(db *sql.DB) error {
	_, err := db.Exec(bannerSQLString[mysqlBannerCreateDatabase])
	return err
}

// CreateTable create banner data table
func CreateTable(db *sql.DB) error {
	_, err := db.Exec(bannerSQLString[mysqlBannerCreateTable])
	return err
}

// InsertBanner add banner information and return bannerId
func InsertBanner(db *sql.DB, name, imagepath, event *string, startdate, enddate *time.Time) (uint32, error) {
	result, err := db.Exec(bannerSQLString[mysqlBannerInsert], name, imagepath, event, startdate, enddate)
	if err != nil {
		return 0, err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errInvalidInsert
	}

	bannerID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint32(bannerID), nil
}

// LisitValidBannerByUnixDate query banner info  Within the specified time
func LisitValidBannerByUnixDate(db *sql.DB, unixtime int64) ([]*Banner, error) {
	var (
		bans []*Banner

		bannerID  int
		name      string
		imagepath string
		eventpath string
		sdate     string
		edate     string
	)

	rows, err := db.Query(bannerSQLString[mysqlBannerListDate], unixtime, unixtime)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&bannerID, &name, &imagepath, &eventpath, &sdate, &edate); err != nil {
			return nil, err
		}

		ban := &Banner{
			BannerID:  bannerID,
			Name:      name,
			ImagePath: imagepath,
			Event:     eventpath,
			StartDate: sdate,
			EndDate:   edate,
		}

		bans = append(bans, ban)
	}

	return bans, nil
}

// InfoByID query banner by bannerid
func InfoByID(db *sql.DB, id int) (*Banner, error) {
	var ban Banner

	rows, err := db.Query(bannerSQLString[mysqlBannerInfoByID], id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&ban.BannerID, &ban.Name, &ban.ImagePath, &ban.Event, &ban.StartDate, &ban.EndDate); err != nil {
			return nil, err
		}
	}

	return &ban, nil
}

// DeleteByID delete banner by id
func DeleteByID(db *sql.DB, id int) error {
	_, err := db.Exec(bannerSQLString[mysqlBannerDeleteByID], id)
	return err
}
