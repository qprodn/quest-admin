package pg

import (
	"database/sql"
	"fmt"
	"quest-admin/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// kratosLogger implements bun.Logging interface
type kratosLogger struct {
	logger *log.Helper
}

func newKratosLogger(logger log.Logger) *kratosLogger {
	return &kratosLogger{
		logger: log.NewHelper(logger),
	}
}

// Printf implements bun.Logging.Printf
func (l *kratosLogger) Printf(format string, v ...any) {
	l.logger.Info(log.LevelInfo, "message", fmt.Sprintf(format, v...))
}

func NewDB(c *conf.Bootstrap, logger log.Logger) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(c.Data.Database.Source)))
	db := bun.NewDB(sqldb, pgdialect.New())
	if err := db.Ping(); err != nil {
		panic(err)
	}
	// 连接池配置
	sqldb.SetMaxIdleConns(int(c.Data.Database.MaxIdleConns))
	sqldb.SetMaxOpenConns(int(c.Data.Database.MaxOpenConns))
	sqldb.SetConnMaxLifetime(time.Duration(c.Data.Database.ConnMaxLifetime) * time.Second)

	// 设置基于 Kratos logger 的 Bun logger
	bun.SetLogger(newKratosLogger(logger))

	return db
}
