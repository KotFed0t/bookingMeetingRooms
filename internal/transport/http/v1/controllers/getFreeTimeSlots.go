package controllers

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/KotFed0t/booking_meeting_rooms/internal/model"
	"github.com/KotFed0t/booking_meeting_rooms/utils"
	"github.com/gin-gonic/gin"
)

func (ctrl *BookingController) GetFreeTimeSlots(c *gin.Context) {
	ctx := context.Background()
	if c.Query("room_id") == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Error: "bad request", Message: "room_id is required"})
		return
	}
	roomId, err := strconv.ParseInt(c.Query("room_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Error: "bad request", Message: "invalid room_id"})
		return
	}

	if c.Query("date") == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Error: "bad request", Message: "date is required"})
		return
	}
	date, err := time.Parse(time.DateOnly, c.Query("date"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Error: "bad request", Message: "invalid date"})
		return
	}

	res, err := ctrl.service.GetFreeTimeslots(ctx, date, roomId)
	if err != nil {
		slog.Error("error from service.GetFreeTimeslots", slog.String("err", err.Error()))
		code, errResponse := utils.MapErrorToResponse(err)
		c.AbortWithStatusJSON(code, errResponse)
		return
	}
	c.JSON(http.StatusOK, res) 
}
