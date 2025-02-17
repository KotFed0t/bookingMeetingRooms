package controllers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/KotFed0t/booking_meeting_rooms/internal/model"
	"github.com/KotFed0t/booking_meeting_rooms/utils"
	"github.com/gin-gonic/gin"
)

func (ctrl *BookingController) CreateBooking(c *gin.Context) {
	var request model.BookingRequest
	ctx := context.Background()

	err := c.ShouldBindJSON(&request)
	if err != nil {
		slog.Error("error while unmarshall request", slog.String("err", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Error: "bad request", Message: err.Error()})
	}

	booking := model.Booking{UserId: request.UserId, RoomId: request.RoomId}

	booking.StartAt, err = time.Parse(time.RFC3339, request.StartAt)
	if err != nil {
		slog.Error("error parsing date", slog.String("err", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Error: "bad request", Message: err.Error()})
		return
	}
	booking.StartAt = booking.StartAt.UTC()

	booking.EndAt, err = time.Parse(time.RFC3339, request.EndAt)
	if err != nil {
		slog.Error("error parsing date", slog.String("err", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Error: "bad request", Message: err.Error()})
		return
	}
	booking.EndAt = booking.EndAt.UTC()

	if booking.StartAt.Before(time.Now()) || booking.EndAt.Before(booking.StartAt) {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Error: "bad request", Message: "invalid date"})
		return
	}

	if dur := booking.EndAt.Sub(booking.StartAt); dur < ctrl.cfg.BookingDuration.Min || dur > ctrl.cfg.BookingDuration.Max {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Error: "bad request", Message: "invalid duration"})
		return
	}

	resp, err := ctrl.service.CreateBooking(ctx, booking)
	if err != nil {
		slog.Error("error from service.CreateBooking", slog.String("err", err.Error()))
		code, errResponse := utils.MapErrorToResponse(err)
		c.AbortWithStatusJSON(code, errResponse)
		return
	}
	c.JSON(http.StatusOK, resp)
}
