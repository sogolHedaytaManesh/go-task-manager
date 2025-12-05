package http

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	_ "net/http/pprof"
	_ "task-manager/docs"
	"task-manager/pkg/monitoring"
	"task-manager/pkg/rest"

	"github.com/gin-gonic/gin"
)

// @title Go My Project API
// @version 1.0
// @description This is the API documentation for Go My Project
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// SetupRouter @host localhost:8080
// @BasePath /
// @schemes http
func (h *Handler) SetupRouter() *gin.Engine {
	// Set Gin to release mode to reduce logging overhead in production
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Set up pprof
	r.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))

	// Programmatically set Swagger host and base path
	//docs.SwaggerInfo.Host = h.config.HostBasePath
	r.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER"))

	// Limit maximum memory for multipart forms (uploads)
	r.MaxMultipartMemory = MaxMultipartMemory

	// Global middlewares
	r.Use(gin.Recovery())   // recover from panics and prevent server crash
	r.Use(CORSMiddleware()) // handle Cross-Origin Resource Sharing
	r.Use(TracingMiddleware())

	// Initialize Prometheus metrics endpoint
	// In a Kubernetes setup, each pod exposes its own /metrics URL,
	// Prometheus server scrapes all pods, and aggregation happens at the Prometheus level.
	_ = monitoring.InitialGinMetrics(
		r,
		h.config.Metrics.Path,
		h.config.Metrics.Port,
		h.config.Metrics.UserName,
		h.config.Metrics.Password,
	)

	// -------------------------------
	// Task CRUD endpoints
	// -------------------------------
	// @tag.name Tasks
	// @tag.description Task management endpoints
	tasks := r.Group("/api/tasks/").Use(TaskMetricsMiddleware(h.TaskMetrics))
	{
		tasks.POST("", h.TaskCreate)
		tasks.GET("", h.TaskList)
		tasks.GET(":id", h.TaskGetByID)
		tasks.PUT(":id", h.TaskUpdate)
		tasks.DELETE(":id", h.TaskDelete)
	}

	// Handle unknown routes
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, rest.NotFound)
	})

	return r
}
