package model

import "time"

type Booking struct {
	UserId  int64
	StartAt time.Time
	EndAt   time.Time
	RoomId  int64
}
