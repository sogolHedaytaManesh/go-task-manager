package db

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

// Config DBConfig holds database configuration and connection pool settings
type Config struct {
	Driver          string        `yaml:"driver" envconfig:"DB_DRIVER"`
	Host            string        `yaml:"host" envconfig:"DB_HOST"`
	Port            int           `yaml:"port" envconfig:"DB_PORT"`
	User            string        `yaml:"user" envconfig:"DB_USER"`
	Password        string        `yaml:"password" envconfig:"DB_PASSWORD"`
	Name            string        `yaml:"name" envconfig:"DB_NAME"`
	SSLMode         string        `yaml:"sslmode" envconfig:"DB_SSLMODE"` // Postgres only
	MaxOpenConns    int           `yaml:"max_open_conns" envconfig:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns    int           `yaml:"max_idle_conns" envconfig:"DB_MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" envconfig:"DB_CONN_MAX_LIFETIME"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time" envconfig:"DB_CONN_MAX_IDLE_TIME"`
}

// Configs DBConfigs holds multiple database configs (Postgres & MySQL)
type Configs struct {
	Postgres Config `yaml:"postgres" envconfig:"POSTGRES"`
	MySQL    Config `yaml:"mysql" envconfig:"MYSQL"`
}

type Manager struct {
	Postgres DB
	MySQL    DB
}

// DB interface represents a generic database abstraction layer.
// It allows repositories and services to interact with any SQL database
// (Postgres, MySQL, etc.) without depending on a specific implementation.
//
// Methods:
//   - QueryContext: Execute a query that returns multiple rows (SELECT).
//   - ExecContext: Execute a query that does NOT return rows (INSERT, UPDATE, DELETE).
//   - GetContext: Execute a query that returns a single row mapped to dest.
//   - SelectContext: Execute a query that returns multiple rows mapped to dest.
//   - Close: Close the database connection gracefully.
//   - Raw: Access the underlying *sqlx.DB instance for advanced usage.
//
// Example usage:
//
//	func NewUserRepository(db DB) *UserRepo {
//	    return &UserRepo{db: db}
//	}
//
//	// SELECT example
//	var users []User
//	db.SelectContext(ctx, &users, "SELECT * FROM users WHERE active = ?", true)
//
//	// INSERT example
//	res, err := db.ExecContext(ctx, "INSERT INTO users(name) VALUES(?)", "Alice")
//	lastID, _ := res.LastInsertId()
//	affectedRows, _ := res.RowsAffected()
type DB interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Close() error
	Raw() *sqlx.DB
}

func NewDB(cfg Configs) (*Manager, error) {
	pg, err := NewPostgresDB(cfg.Postgres)
	if err != nil {
		return nil, err
	}

	my, err := NewMySQLDB(cfg.MySQL)
	if err != nil {
		pg.Close()
		return nil, err
	}

	return &Manager{
		Postgres: pg,
		MySQL:    my,
	}, nil
}

func (m *Manager) Close() {
	if m.Postgres != nil {
		m.Postgres.Close()
	}
	if m.MySQL != nil {
		m.MySQL.Close()
	}
}
