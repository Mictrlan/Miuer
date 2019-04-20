package gin

import (
	"database/sql"
	"net/http"


)

type controller struct {
	db *sql.DB
}

func New(db *sql.DB) *controller {
	return &controller{
		db : db
	}
}



