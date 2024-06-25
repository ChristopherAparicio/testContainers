package user

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func NewConn(connStr string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := testConnection(conn); err != nil {
		return nil, err
	}

	return conn, err
}

func NewConnWithDbConfig(dbConfig DBConfig) (*sql.DB, error) {
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

	return conn, err
}

func testConnection(conn *sql.DB) error {
	_, err := conn.Query("SELECT 1")
	return err
}
