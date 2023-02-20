package DB

import (
	"database/sql"
	"time"

	"webbc/configuration"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func InitializationDB(configuration *configuration.GeneralConfiguration) *bun.DB {

	sqlDB := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithAddr(configuration.DataBaseHost+":"+configuration.DataBasePort),
		pgdriver.WithUser(configuration.DataBaseUser),
		pgdriver.WithPassword(configuration.DataBasePassword),
		pgdriver.WithDatabase(configuration.DataBaseName),
		pgdriver.WithInsecure(true),
		pgdriver.WithReadTimeout(time.Duration(int(configuration.CallTimeoutInSeconds))*time.Second),
	))

	return bun.NewDB(sqlDB, pgdialect.New())
}

func ChangeDB(configuration *configuration.GeneralConfiguration, db *bun.DB) {

	db.DB = sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithAddr(configuration.DataBaseHost+":"+configuration.DataBasePort),
		pgdriver.WithUser(configuration.DataBaseUser),
		pgdriver.WithPassword(configuration.DataBasePassword),
		pgdriver.WithDatabase(configuration.DataBaseName),
		pgdriver.WithInsecure(true),
		pgdriver.WithReadTimeout(time.Duration(int(configuration.CallTimeoutInSeconds))*time.Second),
	))
}
