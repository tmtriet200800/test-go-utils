package pkgPostgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// ConnectionConfig provides values for gRPC connection configuration
type ConnectionConfig struct {
	Host            string
	Port            int
	User            string
	Pass            string
	Database        string
	ConnMaxLifetime time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
}

// NewConnection provides new mysql connection
func NewConnection(ctx context.Context, cfg ConnectionConfig) (db *sqlx.DB) {
	logger := logrus.New()
	
	db, err := sqlx.Connect(
		"postgres", 
		fmt.Sprintf("dbname=%v user=%v password=%v host=%v port=%v sslmode=disable", cfg.Database, cfg.User, cfg.Pass, cfg.Host, cfg.Port),
	)
	if err != nil {
		logger.Fatal(ctx, "[POSTGRES|Connection] %v", err)
		os.Exit(1)
	}

	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	return db
}
