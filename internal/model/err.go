package model

import "errors"

var (
	ErrBookingOutOfWorkingHours = errors.New("booking out of working hours")
	ErrInvalidInputData = errors.New("invalid input data")
	ErrBookingOverlaps = errors.New("booking overlaps")
)