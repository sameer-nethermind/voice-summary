package middleware

import "database/sql"

type AuthJWT struct {
	DB *sql.DB
}
