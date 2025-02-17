package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/KotFed0t/booking_meeting_rooms/data/db"
	"github.com/KotFed0t/booking_meeting_rooms/internal/model"
	dbmodel "github.com/KotFed0t/booking_meeting_rooms/internal/model/dbModel"
)

func (p *Postgres) GetRoomSchedule(ctx context.Context, roomId int64) (*model.RoomSchedule, error) {
	var roomSchedule dbmodel.RoomSchedule
	query := `
		SELECT rs.room_id, rs.start_at, rs.end_at, r.timezone FROM room_schedules rs
		JOIN rooms r using(room_id) 
		WHERE rs.room_id = $1
	`

	err := p.db.GetContext(ctx, &roomSchedule, query, roomId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, db.ErrNoRows
		}
		return nil, err
	}
	return p.convertfromPGRoomSchedule(roomSchedule), nil
}

func (p *Postgres) convertfromPGRoomSchedule(s dbmodel.RoomSchedule) *model.RoomSchedule {
	return &model.RoomSchedule{
		RoomId:         s.RoomId,
		StartAtMinutes: int(s.StartAt.Microseconds / 1e6 / 60),
		EndAtMinutes:   int(s.EndAt.Microseconds / 1e6 / 60),
		Timezone:       s.Timezone,
	}
}
