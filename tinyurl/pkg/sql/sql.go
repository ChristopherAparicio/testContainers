package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// NewDefaultDbConfig: Returns a default DBConfig
// Mostly used for testing
func NewDefaultDbConfig() DBConfig {
	return DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Database: "tinyurl",
	}
}

// TODO : Add more methods
type Querier interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func NewConn(dbConfig DBConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Database,
	)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := testConnection(conn); err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(10)
	conn.SetConnMaxIdleTime(10 * time.Second)

	return conn, err
}

func testConnection(conn *sql.DB) error {
	_, err := conn.Query("SELECT 1")
	return err
}
