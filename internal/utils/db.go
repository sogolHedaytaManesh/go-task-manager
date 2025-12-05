package utils

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"sync"
	"task-manager/pkg/db"
	"testing"
	"time"
)

var (
	initialDBOnce sync.Once
	dbTest        db.DB
)

func TruncateTables(t *testing.T) {
	dbTest = CreateTestDatabaseConnection()
	tables := []string{"tasks"}

	for _, tbl := range tables {
		_, err := dbTest.ExecContext(context.Background(),
			fmt.Sprintf(`TRUNCATE %s RESTART IDENTITY CASCADE`, tbl),
		)
		require.NoError(t, err)
	}
}

func CreateTestDatabaseConnection() db.DB {
	initialDBOnce.Do(func() {
		var err error

		cfg := LoadTestDBConfig()

		dbTest, err = db.NewPostgresDB(cfg)

		if err != nil {
			panic(err)
		}

	})

	return dbTest
}

func LoadTestDBConfig() db.Config {
	return db.Config{
		Host:     "0.0.0.0",
		Port:     6433,
		User:     "test_user",
		Password: "test_pass",
		Name:     "test_db",
		SSLMode:  "disable",

		MaxOpenConns:    3,
		MaxIdleConns:    3,
		ConnMaxLifetime: time.Minute,
		ConnMaxIdleTime: time.Minute,
	}
}
