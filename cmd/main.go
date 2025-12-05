package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"task-manager/internal/config"
	_ "task-manager/internal/config"
	"task-manager/pkg/logger"
	"time"
)

// Set Main Operation
const op = "TaskManager.app"

var (
	GitCommit     = "Development"
	BuildTime     = time.Now().Format(time.RFC1123Z)
	ContainerName string
)

func main() {

	// Load Config File
	var configFile string
	flag.StringVar(&configFile, "c", "", "the environment configuration file of application")
	flag.StringVar(&configFile, "config", "", "the environment configuration file of application")
	flag.Usage = usage
	flag.Parse()

	// Loading the config file
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to load config: %s, error: %v", op, err))
		os.Exit(1)
	}

	// Setup Logger
	loggerInstance := logger.CreateLogger(cfg.Logger)
	loggerInstance.Info("logger configured")
	// Show the loaded config file
	loggerInstance.InfoF("loaded config file: '%s'", configFile)

	// Get the container or host name and log it along with Git commit and build time.
	//
	// Why we do this:
	// 1. ContainerName / hostname: In a distributed environment (e.g., multiple Docker containers
	//    or Kubernetes pods), this allows us to identify which instance is producing the logs.
	// 2. GitCommit: Records the exact version of the code currently running, which is crucial
	//    for debugging, tracing issues, and reproducing bugs in production.
	// 3. BuildTime: Helps verify when this binary was built, providing additional context
	//    for deployment tracking and incident investigations.
	//
	// Logging these details at startup ensures traceability and improves observability,
	// especially in production environments with centralized logging systems like Graylog, ELK, or Loki.
	hostname, err := os.Hostname()
	if err != nil {
		loggerInstance.FatalF("cause:%v,message:%v", err.Error(), op)
	}

	ContainerName = hostname
	loggerInstance.InfoF("hostname acquired :%s", hostname)

	// Commit, BuildTime
	loggerInstance.InfoF("commit number:%s, build time: %s", GitCommit, BuildTime)

	// Create New Server
	server := NewServer(*cfg)

	// Initialize the Server Dependencies
	err = server.Initialize(loggerInstance)
	if err != nil {
		loggerInstance.FatalF("failed to initialize server: %s", err.Error())
	}

	done := make(chan bool, 1)
	quiteSignal := make(chan os.Signal, 1)
	signal.Notify(quiteSignal, syscall.SIGINT, syscall.SIGTERM)

	// Graceful shutdown goroutine
	go server.GracefulShutdown(quiteSignal, done)

	// Setting up the main context
	ctx, cancel := context.WithCancel(context.Background())

	// Start server in blocking mode
	server.Start(ctx)

	// Wait for HTTP Server to be killed gracefully !
	<-done

	// Killing other background jobs !
	cancel()
	loggerInstance.Info("waiting for background jobs to finish their works...")

	// Wait for all other background jobs to finish their works
	server.Wait()

	loggerInstance.Info("TaskManager app shutdown successfully.")

}

func usage() {
	usageStr := `
TaskManager - High Performance Backend Service

Usage:
    TaskManager [options]

Options:
    -c, --config   <file>      Path to YAML configuration file (default: config.yaml)
    -h, --help                 Show this help message
    -v, --version              Show version of the application

Environment Variables:
    DB_TYPE                    Database type to use (postgres or mysql)
    REDIS_ADDR                 Redis address (host:port)
    SERVER_PORT                Server port for HTTP API
    LOG_LEVEL                  Logging level (debug, info, warn, error)

Examples:
    # Run with default config
    TaskManager

    # Run with custom config file
    TaskManager --config ./config.yaml

    # Override environment variable
    SERVER_PORT=9090 TaskManager --config ./config.yaml
`
	fmt.Println(usageStr)
	os.Exit(0)
}
