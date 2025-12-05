package main

import (
	"context"
	"fmt"
	"github.com/cockroachdb/errors"
	"os"
	"sync"
	"task-manager/internal/config"
	"task-manager/internal/http"
	"task-manager/internal/repository/postgres"
	"task-manager/internal/service"
	"task-manager/pkg/db"
	"task-manager/pkg/logger"
	"task-manager/pkg/monitoring"
)

// Global database connection (could also be encapsulated)
var dbConn db.DB
var err error

// Server represents the main application server with all dependencies
type Server struct {
	sync.WaitGroup
	Config      config.Config // Application configuration
	Logger      logger.Logger // Logger instance
	restHandler *http.Handler // REST API handler
}

// NewServer creates a new Server instance with the provided configuration
func NewServer(cfg config.Config) *Server {
	return &Server{
		Config: cfg,
	}
}

// Initialize sets up the application: DB connection, repositories, services, metrics, and HTTP handler
func (s *Server) Initialize(logger logger.Logger) error {
	// Initialize primary DB connection depending on DBType (Postgres / MySQL)
	if s.Config.DBType == "mysql" {
		dbConn, err = db.NewMySQLDB(s.Config.DB.Postgres)
		if err != nil {
			return errors.Wrap(err, "[NOK] failed to initialize MySQL database")
		}
	} else {
		dbConn, err = db.NewPostgresDB(s.Config.DB.Postgres)
		if err != nil {
			return errors.Wrap(err, "[NOK] failed to initialize Postgres database")
		}
	}
	logger.Info("[OK] database connection established")

	// Initialize Prometheus metrics manager
	metricsManager := monitoring.NewMetricsManager()
	logger.Info("[OK] metrics manager initialized")

	// Initialize Task-related metrics
	taskMetrics := monitoring.InitTaskMetrics(metricsManager)

	// Create repositories
	taskRepository := postgres.NewTaskRepository(dbConn)

	// Create services and inject dependencies (repositories + metrics)
	TaskService := service.NewTaskService(taskRepository, taskMetrics)

	s.Logger = logger

	// Initialize REST handler with services, logger, metrics, and config
	s.restHandler = http.CreateHandler(
		s.Logger,
		s.Config,
		TaskService,
		taskMetrics,
	)

	return nil
}

// Start runs the HTTP server in blocking mode
func (s *Server) Start(ctx context.Context) {
	fmt.Println("Starting server with config:", s.Config)
	s.restHandler.StartBlocking(ctx, s.Config.Port)
}

// GracefulShutdown listens for OS signals and performs a clean shutdown of the server
func (s *Server) GracefulShutdown(quitSignal <-chan os.Signal, done chan<- bool) {
	// Wait for OS signal (SIGINT/SIGTERM)
	<-quitSignal

	// Stop the REST HTTP server gracefully
	s.restHandler.Stop()

	// Signal that shutdown is complete
	close(done)
}
