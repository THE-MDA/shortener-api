package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"rest_API/internal/config"
	"rest_API/internal/lib/logger/sl"
	"rest_API/internal/storage/sqlite"
	ssogrpc "rest_API/internal/clients/sso/grpc"
	deletehandler "rest_API/internal/http-server/handlers/delete"
	"rest_API/internal/http-server/handlers/redirect"
	"rest_API/internal/http-server/handlers/url/getall"
	"rest_API/internal/http-server/handlers/url/save"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// init config: cleanenv
	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)

	//logger: slog: import "log/slog"
	log := setupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env), slog.String("version", "123"))
	log.Debug("debug msg are enabled")

	ssoClient,err:=ssogrpc.New(
		context.Background(),
		log,
		cfg.Clients.SSO.Address,
		cfg.Clients.SSO.Timeout,
		cfg.Clients.SSO.RetriesCount,
	)

	if err!=nil{
		log.Error("failed to init sso client", sl.Err(err))
		os.Exit(1)
	}

	ssoClient.IsAdmin(context.Background(),1)
	//storage: sqlite
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage
	router := chi.NewRouter()
	//middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", deletehandler.New(log, storage))
		r.Get("/", getall.New(log, storage))
	})

	//router: chi, "chi render"
	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.IdleTimeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	//run server
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
