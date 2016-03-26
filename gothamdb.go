package gotham_admin

import (
	"database/sql"
)

type GothamDB struct {
	DB *sql.DB
}