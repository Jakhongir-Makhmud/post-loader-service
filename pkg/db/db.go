package db

import (
	"fmt"
	"post-service/pkg/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDB(cfg config.Config) *sqlx.DB {

	dns := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		cfg.GetString("app.database.username"),
		cfg.GetString("app.database.password"),
		cfg.GetString("app.database.dbname"),
		cfg.GetString("app.database.host"),
		cfg.GetString("app.database.port"),
		cfg.GetString("app.database.sslmode"),
	)

	db, err := sqlx.Open("postgres", dns)
	if err != nil {
		panic(err)
	}

	return db
}
