package postgres

import (
	"context"
	"time"

	"github.com/KotFed0t/booking_meeting_rooms/internal/model"
	"github.com/jackc/pgx/v5/pgtype"
)

func (p *Postgres) GetBusyTimeSlots(ctx context.Context, start time.Time, end time.Time, roomId int64) ([]model.TimeSlot, error) {
	query := `
		select lower(booking_time), upper(booking_time) from bookings 
		WHERE booking_time && tsrange($1 , $2, '()')
		and room_id = $3
		order by lower(booking_time) asc
	`

	rows, err := p.db.QueryContext(
		ctx,
		query,
		pgtype.Timestamp{Time: start, Valid: true},
		pgtype.Timestamp{Time: end, Valid: true},
		roomId,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var (
		startBooking time.Time
		endBooking   time.Time
		timeSlots    []model.TimeSlot
	)

	for rows.Next() {
		err := rows.Scan(&startBooking, &endBooking)
		if err != nil {
			return nil, err
		}

		timeSlots = append(timeSlots, model.TimeSlot{Start: startBooking, End: endBooking})
	}

	return timeSlots, nil
}
