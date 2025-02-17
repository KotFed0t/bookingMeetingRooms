package service

import (
	"context"
	"time"

	"github.com/KotFed0t/booking_meeting_rooms/internal/model"
)

type IService interface {
	CreateBooking(ctx context.Context, booking model.Booking) (*model.BookingResponse, error)
	GetFreeTimeslots(ctx context.Context, date time.Time, roomId int64) (*model.TimeSlotsResponse, error)
}