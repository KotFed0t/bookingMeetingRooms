package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/KotFed0t/booking_meeting_rooms/config"
	"github.com/KotFed0t/booking_meeting_rooms/data/db/postgres"
	"github.com/KotFed0t/booking_meeting_rooms/httpserver"
	"github.com/KotFed0t/booking_meeting_rooms/internal/service"
	"github.com/KotFed0t/booking_meeting_rooms/internal/transport/http/v1/controllers"
	"github.com/KotFed0t/booking_meeting_rooms/internal/transport/http/v1/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.MustLoad()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var logLevel slog.Level

	switch cfg.LogLevel {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warning":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(log)

	slog.Debug("config", slog.Any("cfg", cfg))

	postgresDb, err := postgres.NewPostgres(ctx, cfg)
	if err != nil {
		slog.Error("postgres init error", slog.String("err", err.Error()))
		os.Exit(1)
	}

	bookingService := service.NewBookingService(cfg, postgresDb)

	bookingController := controllers.NewBookingController(cfg, bookingService)

	engine := gin.Default()
	routes.SetupRoutes(engine, bookingController)
	httpServer := httpserver.New(engine, cfg)

	// Waiting interruption signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case s := <-interrupt:
		slog.Info("got interruption signal: " + s.String())
	case err := <-httpServer.Notify():
		slog.Error("got httpServer.Notify", slog.String("err", err.Error()))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		slog.Error("httpServer.Shutdown error", slog.String("err", err.Error()))
	}
}
