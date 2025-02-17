package service

import (
	"context"
	"time"

	"github.com/KotFed0t/booking_meeting_rooms/internal/model"
)

func (s *BookingService) GetFreeTimeslots(ctx context.Context, date time.Time, roomId int64) (*model.TimeSlotsResponse, error) {
	roomSchedule, err := s.repo.GetRoomSchedule(ctx, roomId)
	if err != nil {
		return nil, err
	}

	roomLocation, err := time.LoadLocation(roomSchedule.Timezone)
	if err != nil {
		return nil, err
	}

	start, end, ok := s.getSelectingTimeRangeInUTC(*roomSchedule, roomLocation, date)
	if !ok {
		return nil, model.ErrBookingOutOfWorkingHours
	}

	busySlots, err := s.repo.GetBusyTimeSlots(ctx, start, end, roomId)
	if err != nil {
		return nil, err
	}

	var timeSlots []model.TimeSlot

	for _, slot := range busySlots {
		if slot.Start.Before(start) || slot.Start.Equal(start) {
			start = slot.End
		} else if start.Sub(slot.Start) >= s.cfg.BookingDuration.Min {
			timeSlots = append(timeSlots, model.TimeSlot{Start: start.In(roomLocation), End: slot.Start.In(roomLocation)})
			start = slot.End
		}
	}

	if start.Before(end) {
		if end.Sub(start) >= s.cfg.BookingDuration.Min {
			timeSlots = append(timeSlots, model.TimeSlot{Start: start.In(roomLocation), End: end.In(roomLocation)})
		}
	}

	return &model.TimeSlotsResponse{
		TimeSlots: timeSlots,
		RoomId:    roomSchedule.RoomId,
	}, nil
}

func (s *BookingService) getSelectingTimeRangeInUTC(roomSchedule model.RoomSchedule, roomLocation *time.Location, date time.Time) (start time.Time, end time.Time, ok bool) {
	roomStart := roomSchedule.GetStartTime(date.Year(), date.Month(), date.Day(), *roomLocation).UTC()
	roomEnd := roomSchedule.GetEndTime(date.Year(), date.Month(), date.Day(), *roomLocation).UTC()
	now := time.Now().UTC().Truncate(time.Minute)

	if roomSchedule.IsMidnigntRoom() {
		roomEnd = roomEnd.Add(24 * time.Hour)
	}

	if now.After(roomEnd) {
		return start, end, false
	}

	if now.After(roomStart) {
		return now, roomEnd, true
	}

	return roomStart, roomEnd, true
}
