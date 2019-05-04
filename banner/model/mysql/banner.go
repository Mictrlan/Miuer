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
	mysqlBannerLisitDate
	mysqlBannerInfoById
	mysqlBannerDeleteById
)

type Banner struct {
	BannerId  int
	Name      string
	ImagePath string
	Event     string
	StartDate time.Time
	EndDate   time.Time
}

var (
	errInvalidInsert = errors.New("insert banner:insert affected 0 rows")

	bannerSqlString = []string{
		`CREATE DATABASE IF NOT EXISTS banner`,
		`CREATE TABLE IF NOT EXISTS banner.ads(
			bannerid 	BIGINT  NOT NULL AUTO_INCREMENT,
			name 		VARCHAR(512) UNIQUE DEFAULT ' ',
			imagepath   VARCHAR(512) UNIQUE DEFAULT ' ',
			event 		VARCHAR(512) DEFAULT ' ',
			startdate  	DATETIME NOT NULL,
			enddate		DATETIME NOT NULL,
			PRIMARY KEY(bannerid)
		)ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO banner.ads(name,imagepath,event,startdate,enddate)VAULES(?,?,?,?,?)`,
		`SELECT * FROM banner.ads WHERE UNIX_TIMESTAMP(startdate) <= ? AND  UNIX_TIMESTAMP(enddate) >= ? LOCK IN SHARE MODE `,
		`SELECT * FROM banner.ads WHERE bannerid = ? LOCK IN SHARE MODE `,
		`DELETE FROM banner.ads WHERE bannerid = ? LIMIT 1`,
	}
)

// CreateDB create banner database
func CreateDB(db *sql.DB) error {
	_, err := db.Exec(bannerSqlString[mysqlBannerCreateDatabase])
	return err
}

// CreateTable create banner data table
func CreateTable(db *sql.DB) error {
	_, err := db.Exec(bannerSqlString[mysqlBannerCreateTable])
	return err
}

// InsertBanner return bannerId
func InsertBanner(db *sql.DB, name, imagepath, event *string, startdate, enddate *time.Time) (int, error) {
	result, err := db.Exec(bannerSqlString[mysqlBannerInsert], name, imagepath, event, startdate, enddate)
	if err != nil {
		return 0, err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errInvalidInsert
	}

	bannerId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(bannerId), nil

}

func LisitValidBannerByUnixDate(db *sql.DB, unixtime int64) ([]*Banner, error) {
	var (
		bans []*Banner

		bannerId  int
		name      string
		imagepath string
		eventpath string
		sdate     time.Time
		edate     time.Time
	)

	rows, err := db.Query(bannerSqlString[mysqlBannerLisitDate], unixtime)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&bannerId, &name, &imagepath, &eventpath, &sdate, &edate); err != nil {
			return nil, err
		}

		ban := &Banner{
			BannerId:  bannerId,
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

func InfoById(db *sql.DB, id int) (*Banner, error) {
	var ban Banner

	rows, err := db.Query(bannerSqlString[mysqlBannerInfoById], id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&ban.BannerId, &ban.Name, &ban.ImagePath, &ban.Event, &ban.StartDate, &ban.EndDate); err != nil {
			return nil, err
		}
	}
	return &ban, nil
}

func deleteById(db *sql.DB, id int) error {
	_, err := db.Exec(bannerSqlString[mysqlBannerDeleteById], id)
	return err
}
