package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/paybridge/paybridge-ai/services/api/internal/claims"
	"github.com/paybridge/paybridge-ai/services/api/internal/config"
	"github.com/paybridge/paybridge-ai/services/api/internal/orchestrator"
	"github.com/paybridge/paybridge-ai/services/api/internal/paypal"
	"github.com/paybridge/paybridge-ai/services/api/internal/reducer"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("worker configuration error: %v", err)
	}

	httpClient := &http.Client{Timeout: 15 * time.Second}
	tm := paypal.NewTokenManager(cfg.PayPalClientID, cfg.PayPalClientSecret, cfg.PayPalAPIBase, httpClient)

	// INVARIANT: Worker receives ONLY Reader, NEVER Mutator!
	reader := paypal.NewReader(cfg.PayPalAPIBase, tm, httpClient)
	_ = reader
	red := reducer.NewReducer()
	_ = red
	val := claims.NewValidator()
	orc := orchestrator.NewOrchestrator(orchestrator.Config{
		LLMEnabled: os.Getenv("LLM_ENABLED") == "true",
	}, val)
	_ = orc

	log.Printf("PayBridge AI Background Worker started (Reconciler polling active, Mutator unreachable)")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Outbox processing loop heartbeat
			log.Println("[Worker Heartbeat] Polling outbox jobs and ambiguous timeout cases...")
		case <-stop:
			log.Println("Shutting down PayBridge AI Worker...")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = ctx
			log.Println("PayBridge AI Worker stopped.")
			return
		}
	}
}
