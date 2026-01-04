package pg

import (
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"quest-admin/internal/conf"
	"time"
)

func NewDB(c *conf.Bootstrap) *bun.DB {

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(c.Data.Database.Source)))
	db := bun.NewDB(sqldb, pgdialect.New())

	// 连接池配置
	sqldb.SetMaxIdleConns(int(c.Data.Database.MaxIdleConns))
	sqldb.SetMaxOpenConns(int(c.Data.Database.MaxOpenConns))
	sqldb.SetConnMaxLifetime(time.Duration(c.Data.Database.ConnMaxLifetime))

	return db
}
