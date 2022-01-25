//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/robertwtucker/document-host/internal/config"
	"github.com/robertwtucker/document-host/internal/document"
	docrepo "github.com/robertwtucker/document-host/internal/document/repository/mongo"
	dochttp "github.com/robertwtucker/document-host/internal/document/transport/http"
	"github.com/robertwtucker/document-host/internal/document/usecase"
	health "github.com/robertwtucker/document-host/internal/healthcheck/transport/http"
	"github.com/robertwtucker/document-host/pkg/log"
	"github.com/robertwtucker/document-host/pkg/shortlink/tinyurl"
)

// App hods the singletons and use cases
type App struct {
	config     *config.Configuration
	logger     log.Logger
	db         *mongo.Database
	documentUC document.UseCase
}

// CustomValidator provides a validator implementation for echo
type CustomValidator struct {
	validator *validator.Validate
}

// Validate checks to see if the object satisfies validation annotations
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// NewApp initializes and returns an App instance
func NewApp(cfg *config.Configuration, logger log.Logger) (*App, error) {
	logger.Debug("start: wiring App components")

	// Initialize MongoDB
	db, err := initDB(cfg, logger)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error initializing db: %v", err))
	}
	logger.Debug("db connection initialized")

	// Inject the DB into the repo
	documentRepo := docrepo.NewDocumentRepository(db)

	// Initialize the short link generation service
	shortLinkSvc := tinyurl.NewTinyURLService(cfg.ShortLink.APIKey, cfg.ShortLink.Domain)
	logger.Debug("end: wiring App components")

	return &App{
		logger:     logger,
		config:     cfg,
		db:         db,
		documentUC: usecase.NewDocumentUseCase(documentRepo, shortLinkSvc),
	}, nil
}

// Run does the heavy-lifting for the App
func (a *App) Run() {
	a.logger.Debug("start: configure server")

	// Echo setup
	e := echo.New()
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: time.Duration(a.config.Server.Timeout) * time.Second,
	}))
	e.Validator = &CustomValidator{validator: validator.New()}

	// Register HTTP endpoints
	health.RegisterHTTPHandlers(e)

	// Register API endpoints
	// Pattern:
	// --------
	// func RegisterHTTPHandlers(e *echo.Echo, uc resource.UseCase) {
	//   h := NewHandler(uc)
	//	 r := e.Group("/v1")
	//   r.POST("/resource", h.Create)
	//   ...
	// }
	dochttp.RegisterHTTPHandlers(e, a.documentUC)
	a.logger.Debug("end: configuring server")

	// Start server using goroutine
	go func() {
		a.logger.Debug("starting server")
		err := e.Start(":" + a.config.Server.Port)
		if err != nil && err != http.ErrServerClosed {
			e.Logger.Fatalf("shutting down the server: %+v", err)
		}
	}()

	// Create channel for interrupt signals
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint
	// Interrupt received, try to shut down gracefully
	ctx, cancel := context.WithTimeout(
		context.Background(), time.Duration(a.config.Server.Timeout)*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

// initDB sets up the MongoDB client and establishes the DB connection
func initDB(cfg *config.Configuration, logger log.Logger) (*mongo.Database, error) {
	credential := options.Credential{
		AuthSource: cfg.DB.Name,
		Username:   cfg.DB.User,
		Password:   cfg.DB.Password,
	}
	uri := fmt.Sprintf("%s://%s:%s", cfg.DB.Prefix, cfg.DB.Host, cfg.DB.Port)
	client, err := mongo.NewClient(options.Client().ApplyURI(uri).SetAuth(credential))
	logger.Debugf("creating connection to %s", uri)
	if err != nil {
		logger.Errorf("error creating db client: %+v", err)
		return nil, err
	}
	defer client.Disconnect(context.Background())

	// Set a timeout for blocking functions
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(cfg.DB.Timeout)*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		logger.Errorf("error setting timeout context: %+v", err)
		return nil, err
	}

	logger.Debug("pinging db server")
	err = client.Ping(context.Background(), nil)
	if err != nil {
		logger.Errorf("error connecting to db: %+v", err)
		return nil, err
	}

	return client.Database(cfg.DB.Name), nil
}
