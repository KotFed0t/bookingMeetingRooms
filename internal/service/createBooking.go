package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/KotFed0t/booking_meeting_rooms/data/db"
	"github.com/KotFed0t/booking_meeting_rooms/internal/model"
)

func (s *BookingService) CreateBooking(ctx context.Context, booking model.Booking) (*model.BookingResponse, error) {
	isBookingInWorkingHours, err := s.validateWorkingHours(ctx, booking)
	if err != nil {
		return nil, err
	}

	if !isBookingInWorkingHours {
		return nil, model.ErrBookingOutOfWorkingHours
	}
	
	id, err := s.repo.CreateBooking(ctx, booking)
	if err != nil {
		slog.Error("error from repo.CreateBooking", slog.String("err", err.Error()))
		return nil, err
	}

	return &model.BookingResponse{BookingId: id}, nil
}

func (s *BookingService) validateWorkingHours(ctx context.Context, booking model.Booking) (bool, error) {
	// - запросить в PG часы работы
	roomSchedule, err := s.repo.GetRoomSchedule(ctx, booking.RoomId)
	if err != nil {
		slog.Error("error from repo.GetRoomSchedule", slog.String("err", err.Error()))
		if errors.Is(err, db.ErrNoRows) {
			return false, model.ErrInvalidInputData
		}
		return false, err
	}

	// - переводим время запроса в таймзону комнаты
	roomLocation, err := time.LoadLocation(roomSchedule.Timezone)
	if err != nil {
		slog.Error("error while load room location", slog.String("err", err.Error()))
		return false, err
	}

	bookingStart := booking.StartAt.In(roomLocation)
	bookingEnd := booking.EndAt.In(roomLocation)
	roomStart := roomSchedule.GetStartTime(bookingStart.Year(), bookingStart.Month(), bookingStart.Day(), *roomLocation)
	roomEnd := roomSchedule.GetEndTime(bookingEnd.Year(), bookingEnd.Month(), bookingEnd.Day(), *roomLocation)

	if roomSchedule.StartAtMinutes < roomSchedule.EndAtMinutes { // расписания находится в 1 дне
		// работаем только если букинг с 1 днем, если два дня - отклоняем
		if bookingStart.Day() == bookingEnd.Day() {
			if (bookingStart.After(roomStart) || bookingStart.Equal(roomStart)) && (bookingEnd.Before(roomEnd) || bookingEnd.Equal(roomEnd)) {
				return true, nil
			}
		}
	} else { // расписания пересекает 2 дня
		if bookingStart.Day() != bookingEnd.Day() { // букинг тоже пересекает 2 дня
			// тут чисто берем дни из start и end и считаем
			if (bookingStart.After(roomStart) || bookingStart.Equal(roomStart)) && (bookingEnd.Before(roomEnd) || bookingEnd.Equal(roomEnd)) {
				return true, nil
			}
		} else { // букинг 1 день
			// так как букинг входит в 1 день, а расписание в 2 дня - то тут высчитываем вхождение букинга в одну из частей расписания до или после midnight
			if bookingStart.After(roomStart) || bookingStart.Equal(roomStart) {
				midnight := time.Date(bookingStart.Year(), bookingStart.Month(), bookingStart.Day()+1, 0, 0, 0, 0, roomLocation)
				if bookingEnd.Before(midnight) || bookingEnd.Equal(midnight) {
					return true, nil
				}
			} else {
				midnight := time.Date(bookingEnd.Year(), bookingEnd.Month(), bookingEnd.Day(), 0, 0, 0, 0, roomLocation)
				if (bookingStart.After(midnight) || bookingStart.Equal(midnight)) && (bookingEnd.Before(roomEnd) || bookingEnd.Equal(roomEnd)) {
					return true, nil
				}
			}
		}
	}

	return false, nil
}
