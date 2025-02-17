package controllers

import (
	"net/http"

	"github.com/KotFed0t/booking_meeting_rooms/config"
	"github.com/KotFed0t/booking_meeting_rooms/internal/service"
	"github.com/gin-gonic/gin"
)

type BookingController struct{
	cfg *config.Config
	service service.IService
}

func NewBookingController(cfg *config.Config, service service.IService) *BookingController {
	return &BookingController{cfg:cfg, service: service}
}

func (c *BookingController) Test(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "hello world"})
}
