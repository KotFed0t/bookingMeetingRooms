package service

import (
	"github.com/KotFed0t/booking_meeting_rooms/config"
	repository "github.com/KotFed0t/booking_meeting_rooms/data/db"
)

type BookingService struct {
	cfg *config.Config
	repo repository.IRepository
}

func NewBookingService(cfg *config.Config, repo repository.IRepository) *BookingService {
	return &BookingService{cfg: cfg, repo: repo}
}