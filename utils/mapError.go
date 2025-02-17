package utils

import (
	"errors"
	"net/http"

	"github.com/KotFed0t/booking_meeting_rooms/internal/model"
)

func MapErrorToResponse(err error) (int, model.ErrorResponse) {
	switch {
	case errors.Is(err, model.ErrBookingOutOfWorkingHours):
		return http.StatusBadRequest, model.ErrorResponse{Error: "bad request", Message: model.ErrBookingOutOfWorkingHours.Error()}
	case errors.Is(err, model.ErrInvalidInputData):
		return http.StatusBadRequest, model.ErrorResponse{Error: "bad request", Message: model.ErrInvalidInputData.Error()}
	case errors.Is(err, model.ErrBookingOverlaps):
		return http.StatusBadRequest, model.ErrorResponse{Error: "bad request", Message: model.ErrBookingOverlaps.Error()}
	default:
		return http.StatusInternalServerError, model.ErrorResponse{Error: "internal error", Message: "something went wrong"}
	}
}