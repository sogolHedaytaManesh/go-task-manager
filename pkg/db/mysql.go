package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type MySQLDB struct {
	Conn *sqlx.DB
}

func NewMySQLDB(cfg Config) (*MySQLDB, error) {
	// DSN format for MySQL: user:password@tcp(host:port)/dbname?parseTime=true
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)

	// sqlx.Connect combines sqlx.Open and Ping: it creates a DB object and
	// immediately tests the connection. Use this for initial bootstrap to
	// ensure the database is reachable before running queries.
	// If you prefer lazy connection (defer actual DB connection until first query),
	// use sqlx.Open instead and call db.Ping() manually when needed.
	dbConn, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Connection Pool Settings:
	// - SetMaxOpenConns: maximum number of open connections to the database.
	//   Queries exceeding this limit will wait until a connection becomes available.
	// - SetMaxIdleConns: maximum number of idle (ready-to-use) connections maintained.
	//   Helps reduce latency for the first queries.
	// - SetConnMaxLifetime: maximum lifetime of a connection before being closed and replaced.
	//   Prevents stale connections and long-lived timeouts.
	// - SetConnMaxIdleTime: maximum amount of time an idle connection can remain before being closed.
	dbConn.SetMaxOpenConns(cfg.MaxOpenConns)
	dbConn.SetMaxIdleConns(cfg.MaxIdleConns)
	dbConn.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	dbConn.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return &MySQLDB{Conn: dbConn}, nil
}

func (m *MySQLDB) Close() error {
	return m.Conn.Close()
}

func (m *MySQLDB) Raw() *sqlx.DB {
	return m.Conn
}

// QueryContext Forward methods to underlying sqlx.DB
func (m *MySQLDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return m.Conn.QueryxContext(ctx, query, args...)
}

func (m *MySQLDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return m.Conn.ExecContext(ctx, query, args...)
}

func (m *MySQLDB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return m.Conn.GetContext(ctx, dest, query, args...)
}

func (m *MySQLDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return m.Conn.SelectContext(ctx, dest, query, args...)
}
