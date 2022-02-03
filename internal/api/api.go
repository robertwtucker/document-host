//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robertwtucker/document-host/internal/config"
	doc "github.com/robertwtucker/document-host/internal/document"
	docrepo "github.com/robertwtucker/document-host/internal/document/repository/mongo"
	dochttp "github.com/robertwtucker/document-host/internal/document/transport/http"
	docuc "github.com/robertwtucker/document-host/internal/document/usecase"
	hc "github.com/robertwtucker/document-host/internal/healthcheck"
	hcrepo "github.com/robertwtucker/document-host/internal/healthcheck/repository/mongo"
	hchttp "github.com/robertwtucker/document-host/internal/healthcheck/transport/http"
	hcuc "github.com/robertwtucker/document-host/internal/healthcheck/usecase"
	"github.com/robertwtucker/document-host/pkg/shortlink/tinyurl"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// App hods the singletons and use cases
type App struct {
	config     *config.Configuration
	db         *mongo.Database
	documentUC doc.UseCase
	healthUC   hc.UseCase
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
func NewApp(cfg *config.Configuration) (*App, error) {
	log.Debug("start: wiring app components")

	// Initialize MongoDB
	db, err := initDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("error initializing db: %v", err)
	}
	log.Debugf("connection to %s db initialized", db.Name())

	// Inject the DB into the repo
	documentRepo := docrepo.NewDocumentRepository(db)

	// Initialize the short link generation service
	shortLinkSvc := tinyurl.NewTinyURLService(cfg.ShortLink.APIKey, cfg.ShortLink.Domain)

	// Inject the DB into the helper
	dbHelper := hcrepo.NewHealthCheckDatabaseHelper(db)
	log.Debug("end: wiring app components")

	return &App{
		config:     cfg,
		db:         db,
		documentUC: docuc.NewDocumentUseCase(documentRepo, shortLinkSvc, cfg),
		healthUC:   hcuc.NewHealthCheckUseCase(dbHelper),
	}, nil
}

// Run does the heavy-lifting for the App
func (a *App) Run() {
	log.Debug("start: configuring server")

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
	hchttp.RegisterHTTPHandlers(e, a.healthUC)

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
	log.Debug("end: configuring server")

	// Start server using goroutine
	go func() {
		log.Debug("starting the server")
		err := e.Start(":" + a.config.Server.Port)
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("shutting down the server: %+v", err)
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
		log.Fatal(err)
	}
}

// initDB sets up the MongoDB client and establishes the DB connection
func initDB(cfg *config.Configuration) (*mongo.Database, error) {
	// User configuration
	var authUser string
	var authSource = "admin"
	if authUser = cfg.DB.User; strings.ToLower(authUser) != "root" {
		authSource = cfg.DB.Name
	}
	credential := options.Credential{
		AuthSource: authSource,
		Username:   authUser,
		Password:   cfg.DB.Password,
	}

	// Client configuration
	uri := fmt.Sprintf("%s://%s:%d", cfg.DB.Prefix, cfg.DB.Host, cfg.DB.Port)
	log.Debugf("creating db client for %s.%s@%s:%d",
		authSource, authUser, cfg.DB.Host, cfg.DB.Port)
	opts := options.Client().ApplyURI(uri).SetAuth(credential)
	client, err := mongo.NewClient(opts)
	if err != nil {
		log.Errorf("error creating db client: %+v", err)
		return nil, err
	}

	// Set a timeout for blocking functions
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(cfg.DB.Timeout)*time.Second)
	defer cancel()

	// Create connection using the timeout context
	if err := client.Connect(ctx); err != nil {
		log.Errorf("error connecting client: %+v", err)
		return nil, err
	}

	log.Debug("validating connection to db (ping)")
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Errorf("error connecting to db: %+v", err)
		return nil, err
	}

	return client.Database(cfg.DB.Name), nil
}
