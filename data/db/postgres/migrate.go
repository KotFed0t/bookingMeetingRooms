package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func migrate(ctx context.Context, db *sqlx.DB) error {
	query := `
		CREATE EXTENSION IF NOT EXISTS btree_gist;

		CREATE TABLE IF NOT EXISTS rooms
		(
			room_id BIGSERIAL PRIMARY KEY,
			address TEXT NOT NULL,
			timezone TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS users
		(
			user_id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			timezone TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS room_schedules
		(
			room_id BIGINT NOT NULL REFERENCES rooms(room_id) ON DELETE CASCADE,
			start_at TIME NOT NULL,
			end_at TIME NOT NULL
		);

		CREATE TABLE IF NOT EXISTS bookings
		(
			booking_id BIGSERIAL PRIMARY KEY,
			room_id BIGINT NOT NULL REFERENCES rooms(room_id) ON DELETE CASCADE,
			booking_time TSRANGE NOT NULL,
			user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
			CONSTRAINT unique_booking_time_exclude
        		EXCLUDE USING gist (room_id WITH =, booking_time WITH &&)
		);
	`

	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}