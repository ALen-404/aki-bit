package postgresql

import (
	"btc_order/internal/config"
	"github.com/go-pg/pg/v10"
)

type Database struct {
	*pg.DB
}

func NewPostgresqlDriver(config *config.PgDB) *Database {
	db := pg.Connect(&pg.Options{
		Addr:     config.Addr,
		User:     config.User,
		Password: config.Password,
		Database: config.Database,
	})
	return &Database{
		db,
	}
}
