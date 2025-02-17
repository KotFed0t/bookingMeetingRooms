package routes

import (
	"github.com/KotFed0t/booking_meeting_rooms/internal/transport/http/v1/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(engine *gin.Engine, bookingController *controllers.BookingController) {
	apiV1Group := engine.Group("api/v1")
	apiV1Group.GET("/test", bookingController.Test)
	apiV1Group.POST("/booking", bookingController.CreateBooking)
	apiV1Group.GET("/timeslotList", bookingController.GetFreeTimeSlots)
}
