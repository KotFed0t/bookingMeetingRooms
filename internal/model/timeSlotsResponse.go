package model

import "time"

type TimeSlotsResponse struct {
	TimeSlots   []TimeSlot `json:"time_slots"`
	RoomId      int64      `json:"room_id"`
}

type TimeSlot struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}
