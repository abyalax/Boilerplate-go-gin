package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/abyalax/Boilerplate-go-gin/src/conf/logger"
	middlewares "github.com/abyalax/Boilerplate-go-gin/src/middleware"
	"github.com/abyalax/Boilerplate-go-gin/src/modules/auth"
	users "github.com/abyalax/Boilerplate-go-gin/src/modules/users"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

// App represents the application with all dependencies
type App struct {
	server *http.Server
	logger *zap.Logger
	db     *pgxpool.Pool
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

	userModule := users.NewUserModule(db, logger)
	authModule := auth.NewAuthModule(db, logger)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(middlewares.LoggingMiddleware(logger))
	router.Use(middlewares.RecoveryMiddleware(logger))
	router.Use(middlewares.ErrorHandler(logger))

	v1 := router.Group("/api/v1")

	userModule.RegisterRoutes(v1, logger)
	authModule.RegisterRoutes(v1, logger)

	v1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	v1.GET("/ready", func(c *gin.Context) {
		if err := db.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

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

func initDatabase(logger *zap.Logger, dbURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	logger.Info("database connected successfully")

	return pool, nil
}
