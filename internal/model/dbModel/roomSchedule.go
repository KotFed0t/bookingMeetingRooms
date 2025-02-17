package dbmodel

import "github.com/jackc/pgx/v5/pgtype"

type RoomSchedule struct {
	RoomId   int64       `db:"room_id"`
	StartAt  pgtype.Time `db:"start_at"`
	EndAt    pgtype.Time `db:"end_at"`
	Timezone string      `db:"timezone"`
}
