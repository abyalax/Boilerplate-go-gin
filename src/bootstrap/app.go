package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/abyalax/Boilerplate-go-gin/src/config/logger"
	middlewares "github.com/abyalax/Boilerplate-go-gin/src/middleware"
	"github.com/abyalax/Boilerplate-go-gin/src/modules/users"
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

// NewApp init the application
func NewApp(dbURL string, port int) (*App, error) {
	// Initialize logger
	logger := logger.GetLogger()

	// Init database
	db, err := initDatabase(logger, dbURL)
	if err != nil {
		logger.Fatal("failed to initialize database", zap.Error(err))
	}

	userQueries := users.New(db)                     // Initialize repository layer
	userService := users.NewUserService(userQueries) // Initialize service layer
	userHandler := users.NewUserHandler(             // Initialize HTTP handlers
		userService,
		logger,
	)

	// Init Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

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
	a.logger.Info("Application running on http://localhost:4000")
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
