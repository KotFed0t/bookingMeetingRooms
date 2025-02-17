package model

type BookingRequest struct {
	UserId  int64  `json:"user_id"`
	StartAt string `json:"start_at"`
	EndAt   string `json:"end_at"`
	RoomId  int64  `json:"room_id"`
}
