package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	authmw "github.com/raul/glovebox-go-migration/go-service/internal/middleware"
	"github.com/raul/glovebox-go-migration/go-service/internal/repo"
	"github.com/raul/glovebox-go-migration/go-service/internal/service"
	"github.com/raul/glovebox-go-migration/go-service/internal/telemetry"
	transport "github.com/raul/glovebox-go-migration/go-service/internal/transport/http"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal server error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	tp, logger, err := telemetry.InitTelemetry(ctx, "go-service")
	if err != nil {
		return err
	}
	defer func() {
		_ = tp.Shutdown(context.Background())
	}()

	dbURL := mustEnv("DATABASE_URL")
	jwtSecret := mustEnv("JWT_SECRET")
	port := envOr("PORT", "8080")

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return err
	}

	// Wiring: repo (Postgres adapter) -> service (business logic) -> handler (HTTP).
	// Each layer only knows about the one below it via an interface, so swapping
	// e.g. the DB or adding a gRPC transport later doesn't ripple through the stack.
	phRepo := repo.NewPolicyholderPG(pool)
	phSvc := service.NewPolicyholderService(phRepo)
	phHandler := transport.NewPolicyholderHandler(phSvc)

	policyRepo := repo.NewPolicyPG(pool)
	policySvc := service.NewPolicyService(policyRepo)
	policyHandler := transport.NewPolicyHandler(policySvc)

	noteRepo := repo.NewNotePG(pool)
	noteSvc := service.NewNoteService(noteRepo)
	noteHandler := transport.NewNoteHandler(noteSvc)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(telemetry.TraceMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r.Route("/api/v1", func(api chi.Router) {
		api.Use(authmw.Auth([]byte(jwtSecret)))
		phHandler.Routes(api)
		policyHandler.Routes(api)
		noteHandler.Routes(api)
	})

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("go-service listening", "port", port)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
	case <-ctx.Done():
		logger.Info("shutdown signal received")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
	return nil
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		slog.Error("missing required env var", "key", key)
		os.Exit(1)
	}
	return v
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
