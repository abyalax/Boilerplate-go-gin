package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	middlewares "github.com/abyalax/Boilerplate-go-gin/internal/middleware"
	user_handlers "github.com/abyalax/Boilerplate-go-gin/internal/users/handlers"
	user_repositories "github.com/abyalax/Boilerplate-go-gin/internal/users/repositories/users"
	user_services "github.com/abyalax/Boilerplate-go-gin/internal/users/services"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

// App represents the application with all dependencies
type App struct {
	server *http.Server
	logger *zap.Logger
	db     *sql.DB
}

// NewApp initializes the application
func NewApp(dbURL string, port int) (*App, error) {
	// Initialize logger
	logger, err := initLogger()
	if err != nil {
		return nil, err
	}

	// Initialize database
	db, err := initDatabase(logger, dbURL)
	if err != nil {
		logger.Fatal("failed to initialize database", zap.Error(err))
	}

	// Initialize service layer
	userQueries := user_repositories.New(db)
	userService := user_services.NewUserService(userQueries)

	// Initialize HTTP handlers
	userHandler := user_handlers.NewUserHandler(
		userService,
		logger,
	)

	// Initialize Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Add middlewares
	router.Use(middlewares.LoggingMiddleware(logger))
	router.Use(middlewares.RecoveryMiddleware(logger))
	router.Use(middlewares.ErrorHandler(logger))

	// Register routes
	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "healthy"})
		})

		// Ready check
		v1.GET("/ready", func(c *gin.Context) {
			if err := db.PingContext(c.Request.Context()); err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "ready"})
		})
	}

	// Initialize HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &App{
		server: server,
		logger: logger,
		db:     db,
	}, nil
}

// Start starts the application
func (a *App) Start() error {
	a.logger.Info("starting server", zap.String("address", a.server.Addr))
	return a.server.ListenAndServe()
}

// Stop gracefully shuts down the application
func (a *App) Stop(ctx context.Context) error {
	a.logger.Info("shutting down server")

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("failed to shutdown server", zap.Error(err))
		return err
	}

	a.db.Close()
	return nil
}

// initLogger initializes the zap logger
func initLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return logger, nil
}

// initDatabase initializes the database connection using sql.DB with pgx driver
func initDatabase(logger *zap.Logger, dbURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test connection
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	logger.Info("database connected successfully")
	return db, nil
}
