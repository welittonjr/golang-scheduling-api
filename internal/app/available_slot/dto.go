package availableslot

import "time"

type AvailableSlotsInput struct {
	StaffID int       `json:"staff_id"`
	Date    time.Time `json:"date"`
}

type AvailableSlotOutput struct {
	Time time.Time `json:"time"`
}