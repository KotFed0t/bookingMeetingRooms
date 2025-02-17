package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/KotFed0t/booking_meeting_rooms/config"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	db  *sqlx.DB
	cfg *config.Config
}

const (
	defaultConnAttemts = 10
	connTimeout        = time.Second
)

func NewPostgres(ctx context.Context, c *config.Config) (*Postgres, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.DbName,
		c.Postgres.Password,
	)

	connAttempts := defaultConnAttemts
	var db *sqlx.DB
	var err error

	for connAttempts > 0 {
		db, err = sqlx.Connect("pgx", dataSourceName)
		if err == nil {
			break
		}

		slog.Info(fmt.Sprintf("Postgres is trying to connect, attempts left: %d", connAttempts))

		time.Sleep(connTimeout)

		connAttempts--
	}

	if err != nil {
		slog.Error("Postgres connAttempts = 0")
		panic(err)
	}

	db.SetMaxOpenConns(c.Postgres.MaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(c.Postgres.ConnMaxLifetime) * time.Second)
	db.SetMaxIdleConns(c.Postgres.MaxIdleConns)
	db.SetConnMaxIdleTime(time.Duration(c.Postgres.ConnMaxIdleTime) * time.Second)
	if err = db.Ping(); err != nil {
		slog.Error("Postgres dbPing error")
		panic(err)
	}

 	err = migrate(ctx, db)
	if err != nil {
		return &Postgres{}, err
	}

	return &Postgres{db: db, cfg: c}, nil
}
