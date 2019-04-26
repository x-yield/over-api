package tools

import (
	"context"

	"github.com/go-pg/pg"

	"github.com/x-yield/over-api/internal/config"
)

func NewDbConnector() *pg.DB {
	// go-pg has its own concurrent-safe connection pool
	conn := pg.Connect(&pg.Options{
		User:     config.GetValue(context.Background(), config.DbUser).String(),
		Addr:     config.GetValue(context.Background(), config.DbAddr).String(),
		Database: config.GetValue(context.Background(), config.DbName).String(),
		Password: config.GetValue(context.Background(), config.DbPass).String(),
	})
	return conn
}
