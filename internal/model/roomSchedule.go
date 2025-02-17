package model

import "time"

type RoomSchedule struct {
	RoomId         int64
	StartAtMinutes int
	EndAtMinutes   int
	Timezone       string
}

func (r *RoomSchedule) GetStartTime(year int, month time.Month, day int, location time.Location) time.Time {
	return time.Date(year, month, day, r.StartAtMinutes / 60, r.StartAtMinutes % 60, 0, 0, &location)
}

func (r *RoomSchedule) GetEndTime(year int, month time.Month, day int, location time.Location) time.Time {
	return time.Date(year, month, day, r.EndAtMinutes / 60, r.EndAtMinutes % 60, 0, 0, &location)
}

func (r *RoomSchedule) IsMidnigntRoom() bool {
	return r.StartAtMinutes > r.EndAtMinutes
}