package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	Conn *sqlx.DB
}

func NewPostgresDB(cfg Config) (*PostgresDB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	// sqlx.Connect combines sqlx.Open and Ping: it creates a DB object and
	// immediately tests the connection. Use this for initial bootstrap to
	// ensure the database is reachable before running queries.
	// If you prefer lazy connection (defer actual DB connection until first query),
	// use sqlx.Open instead and call db.Ping() manually when needed.
	// dbConn, err := sqlx.Open("postgres", m.cfg.DB.dsn)
	//err = dbConn.Ping()
	//if err != nil {
	//		return err
	//	}
	dbConn, err := sqlx.Connect("postgres", dsn)
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

	return &PostgresDB{Conn: dbConn}, nil
}

func (p *PostgresDB) Close() error {
	return p.Conn.Close()
}

func (p *PostgresDB) Raw() *sqlx.DB {
	return p.Conn
}

// QueryContext Forward methods to underlying sqlx.DB
func (p *PostgresDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return p.Conn.QueryxContext(ctx, query, args...)
}

func (p *PostgresDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return p.Conn.ExecContext(ctx, query, args...)
}

func (p *PostgresDB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return p.Conn.GetContext(ctx, dest, query, args...)
}

func (p *PostgresDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return p.Conn.SelectContext(ctx, dest, query, args...)
}
