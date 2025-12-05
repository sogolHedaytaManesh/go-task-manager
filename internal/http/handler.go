package http

import (
	"context"
	"fmt"
	"net/http"
	"task-manager/internal/config"
	"task-manager/internal/service"
	"task-manager/pkg/logger"
	"task-manager/pkg/monitoring"
	"time"

	"github.com/cockroachdb/errors"
)

// -------------------------------
// Server Constants
// -------------------------------

// DefaultTimeOutForGracefulShutDown defines how long the server will wait
// for ongoing requests to finish during shutdown before forcefully closing.
const DefaultTimeOutForGracefulShutDown = 5 * time.Second

// IdleTimeout is the maximum amount of time to wait for the next request
// when keep-alives are enabled. Helps free resources for idle connections.
const IdleTimeout = 60 * time.Second

// ReadTimeout specifies the maximum duration for reading the entire
// request, including the body. Protects against slow clients.
const ReadTimeout = 15 * time.Second

// WriteTimeout specifies the maximum duration before timing out writes
// of the response. Protects against clients that are slow to read.
const WriteTimeout = 15 * time.Second

// MaxMultipartMemory defines the maximum memory for multipart forms (file uploads).
// Files larger than this limit will be stored on disk to prevent excessive RAM usage.
const MaxMultipartMemory = 8 << 20 // 8 MiB

// -------------------------------
// Handler Struct
// -------------------------------

// Handler contains HTTP server, services, logger, metrics, and version info.
type Handler struct {
	TaskService service.TaskService
	logger      logger.Logger
	HTTPServer  *http.Server

	VersionInfo struct {
		GitCommit     string
		BuildTime     string
		StartTime     time.Time
		ContainerName string
	}

	config      config.Config
	TaskMetrics *monitoring.TaskMetrics
}

// -------------------------------
// Constructor
// -------------------------------

// CreateHandler initializes a new HTTP handler with all dependencies.
func CreateHandler(
	logger logger.Logger,
	config config.Config,
	TaskService service.TaskService,
	TaskMetrics *monitoring.TaskMetrics,
) *Handler {
	return &Handler{
		logger:      logger,
		config:      config,
		TaskService: TaskService,
		TaskMetrics: TaskMetrics,
	}
}

// -------------------------------
// Server Lifecycle Methods
// -------------------------------

// StartBlocking starts the HTTP server and blocks the main goroutine.
// Sets up proper timeouts to protect against slow clients and attacks.
func (h *Handler) StartBlocking(ctx context.Context, defaultPort int) {
	addr := fmt.Sprintf(":%v", defaultPort)

	h.HTTPServer = &http.Server{
		Addr:         addr,
		Handler:      h.SetupRouter(),
		WriteTimeout: WriteTimeout,
		ReadTimeout:  ReadTimeout,
		IdleTimeout:  IdleTimeout,
	}

	h.logger.InfoF("[OK] Starting HTTP REST Server on %s", addr)
	err := h.HTTPServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		h.logger.Error(err.Error())
	}

	// Code reaches here only after HTTP server shuts down
	h.logger.Info("[OK] HTTP REST Server is shutting down!")
}

// Stop gracefully shuts down the HTTP server within DefaultTimeOutForGracefulShutDown.
// Any in-flight requests will be given up to 5 seconds to complete.
func (h *Handler) Stop() {
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), DefaultTimeOutForGracefulShutDown)
	defer cancelTimeout()

	h.HTTPServer.SetKeepAlivesEnabled(false)
	if err := h.HTTPServer.Shutdown(ctxTimeout); err != nil {
		h.logger.Error(err.Error())
	}

	h.logger.Info("[OK] HTTP REST Server graceful shutdown completed")
}
