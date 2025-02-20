package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"log"
	"time"
)

type Config struct {
	Env             string `env:"ENV"`
	LogLevel        string `env:"LOG_LEVEL"`
	Postgres        Postgres
	HttpServer      HttpServer
	BookingDuration BookingDuration
}

type Postgres struct {
	Host            string `env:"PG_HOST"`
	Port            int    `env:"PG_PORT"`
	DbName          string `env:"PG_DB_NAME"`
	Password        string `env:"PG_PASSWORD"`
	User            string `env:"PG_USER"`
	PoolMax         int    `env:"PG_POOL_MAX"`
	MaxOpenConns    int    `env:"PG_MAX_OPEN_CONNS"`
	ConnMaxLifetime int    `env:"PG_CONN_MAX_LIFETIME"`
	MaxIdleConns    int    `env:"PG_MAX_IDLE_CONNS"`
	ConnMaxIdleTime int    `env:"PG_CONN_MAX_IDLE_TIME"`
}

type HttpServer struct {
	Address         string        `env:"HTTP_SERVER_ADDRESS"`
	Timeout         time.Duration `env:"HTTP_SERVER_TIMEOUT"`
	IdleTimeout     time.Duration `env:"HTTP_SERVER_IDLE_TIMEOUT"`
	ShutdownTimeout time.Duration `env:"HTTP_SERVER_SHUTDOWN_TIMEOUT"`
}

type BookingDuration struct {
	Min time.Duration `env:"BOOKING_MIN_DURATION"`
	Max time.Duration `env:"BOOKING_MAX_DURATION"`
}

func MustLoad() *Config {
	_ = godotenv.Load(".env")

	cfg := &Config{}

	opts := env.Options{RequiredIfNoDef: true}

	if err := env.ParseWithOptions(cfg, opts); err != nil {
		log.Fatalf("parse config error: %s", err)
	}

	return cfg
}
