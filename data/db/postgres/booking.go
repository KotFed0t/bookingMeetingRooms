package postgres

import (
	"context"
	"errors"
	"log/slog"

	"github.com/KotFed0t/booking_meeting_rooms/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func (p *Postgres) CreateBooking(ctx context.Context, booking model.Booking) (bookingId int64, err error) {
	query := `
		INSERT INTO bookings(room_id, booking_time, user_id)
		VALUES($1, $2, $3)
		RETURNING booking_id;
	`
	bookingRange := pgtype.Range[pgtype.Timestamp]{
		Lower:     pgtype.Timestamp{Time: booking.StartAt, Valid: true},
		Upper:     pgtype.Timestamp{Time: booking.EndAt, Valid: true},
		LowerType: pgtype.Inclusive,
		UpperType: pgtype.Exclusive,
		Valid:     true,
	}

	err = p.db.QueryRowContext(ctx, query, booking.RoomId, bookingRange, booking.UserId).Scan(&bookingId)
	if err != nil {
		slog.Error("Error in CreateBooking", slog.String("err", err.Error()))
		
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "unique_booking_time_exclude" {
				return 0, model.ErrBookingOverlaps
			}
			if pgErr.Code == "23503" { // foreign_key_violation
				return 0, model.ErrInvalidInputData
			}
		}
		return 0, err
	}

	return bookingId, nil
}
