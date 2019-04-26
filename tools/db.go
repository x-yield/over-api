package tools

import (
	"github.com/go-pg/pg"
)

func NewDbConnector() *pg.DB {
	// go-pg has its own concurrent-safe connection pool
	conn := pg.Connect(&pg.Options{
		User:     "username",
		Addr:     "address:port",
		Database: "database name",
		Password: "password",
	})
	return conn
}
