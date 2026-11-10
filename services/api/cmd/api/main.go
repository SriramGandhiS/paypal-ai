package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/paybridge/paybridge-ai/services/api/internal/api"
	"github.com/paybridge/paybridge-ai/services/api/internal/auth"
	"github.com/paybridge/paybridge-ai/services/api/internal/cases"
	"github.com/paybridge/paybridge-ai/services/api/internal/checkout"
	"github.com/paybridge/paybridge-ai/services/api/internal/config"
	"github.com/paybridge/paybridge-ai/services/api/internal/evidence"
	"github.com/paybridge/paybridge-ai/services/api/internal/httpx"
	"github.com/paybridge/paybridge-ai/services/api/internal/paypal"
	"github.com/paybridge/paybridge-ai/services/api/internal/reconciliation"
	"github.com/paybridge/paybridge-ai/services/api/internal/store"
	"github.com/paybridge/paybridge-ai/services/api/internal/webhook"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("starting PayBridge AI API service...")

	// 1. Load validated configuration (fail-fast, redacted credentials)
	cfg, err := config.Load()
	if err != nil {
		logger.Error("configuration initialization failed", "error", err.Error())
		os.Exit(1)
	}

	logger.Info("configuration loaded successfully",
		"paypal_env", cfg.PayPalEnv,
		"paypal_api_base", cfg.PayPalAPIBase,
	)

	// 2. Initialize database connection pool
	storeCfg := store.DefaultStoreConfig(cfg.DatabaseURL)
	dbStore, err := store.NewStore(storeCfg)
	if err != nil {
		logger.Error("database connection initialization failed", "error", err.Error())
		os.Exit(1)
	}
	defer func() {
		if err := dbStore.Close(); err != nil {
			logger.Error("error closing database pool", "error", err)
		}
	}()

	// 3. Assemble domain services
	authReader := auth.NewDatabaseSessionStore(dbStore)
	authMiddleware := auth.NewMiddleware(auth.MiddlewareConfig{
		Reader:         authReader,
		Logger:         logger,
		AllowDevBypass: cfg.AllowDevAuth,
	})

	casesService := cases.NewService(dbStore)
	checkoutService := checkout.NewService(dbStore)
	evidenceService := evidence.NewService(dbStore)
	reconciler := reconciliation.NewReconciler(dbStore, logger)

	payPalClient := paypal.NewClient(paypal.ClientConfig{
		APIBase:      cfg.PayPalAPIBase,
		ClientID:     cfg.PayPalClientID,
		ClientSecret: cfg.PayPalClientSecret,
	})

	webhookID := os.Getenv("PAYPAL_WEBHOOK_ID")
	webhookGateway := webhook.NewGateway(payPalClient, dbStore, webhookID, logger)

	// 4. Assemble HTTP Router and register routes
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	serverCfg := httpx.DefaultServerConfig(":"+port, dbStore)
	serverCfg.Logger = logger

	router := httpx.NewRouter(serverCfg)

	api.RegisterRoutes(router, api.RouterDependencies{
		Cases:          casesService,
		Checkout:       checkoutService,
		Evidence:       evidenceService,
		Reconciler:     reconciler,
		PayPal:         payPalClient,
		WebhookGateway: webhookGateway,
		AuthMiddleware: authMiddleware,
		Logger:         logger,
	})

	srv := &http.Server{
		Addr:         serverCfg.Addr,
		Handler:      router,
		ReadTimeout:  serverCfg.ReadTimeout,
		WriteTimeout: serverCfg.WriteTimeout,
		IdleTimeout:  serverCfg.IdleTimeout,
	}

	// 5. Graceful Shutdown Listener
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("HTTP server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server failure", "error", err)
			os.Exit(1)
		}
	}()

	<-shutdownChan
	logger.Info("shutdown signal received, commencing graceful termination...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}

	logger.Info("PayBridge AI API exited cleanly")
}
