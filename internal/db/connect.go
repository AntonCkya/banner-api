package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Connect struct {
	Conn *sql.DB
}
