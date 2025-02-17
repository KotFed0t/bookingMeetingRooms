package db

import (
	"context"
	"time"

	"github.com/KotFed0t/booking_meeting_rooms/internal/model"
)

type IRepository interface {
	GetRoomSchedule(ctx context.Context, roomId int64) (*model.RoomSchedule, error)
	CreateBooking(ctx context.Context, booking model.Booking) (bookingId int64, err error)
	GetBusyTimeSlots(ctx context.Context, start time.Time, end time.Time, roomId int64) ([]model.TimeSlot, error)
}
