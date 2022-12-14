package DB

import (
	"database/sql"

	"webbc/configuration"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func InitializationDB(configuration *configuration.Configuration) *bun.DB {

	sqlDB := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithAddr(configuration.DataBaseHost+":"+configuration.DataBasePort),
		pgdriver.WithUser(configuration.DataBaseUser),
		pgdriver.WithPassword(configuration.DataBasePassword),
		pgdriver.WithDatabase(configuration.DataBaseName),
		pgdriver.WithInsecure(true),
	))

	return bun.NewDB(sqlDB, pgdialect.New())
}
