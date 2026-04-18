package dtos

import "time"

// StartTripRequest defines the input for starting a new trip.
type StartTripRequest struct {
	RouteID   string `json:"route_id" validate:"required,uuid"`
	VehicleID string `json:"vehicle_id" validate:"required,uuid"`
	DriverID  string `json:"driver_id" validate:"required,uuid"`
}

// EndTripRequest defines the input for ending a trip.
type EndTripRequest struct {
	EndLocation *LocationDTO `json:"end_location,omitempty"`
}

// LocationDTO represents a geographic coordinate.
type LocationDTO struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// TripResponse represents the public data for a trip.
type TripResponse struct {
	ID        string       `json:"id"`
	RouteID   string       `json:"route_id"`
	DriverID  string       `json:"driver_id"`
	VehicleID string       `json:"vehicle_id"`
	Status    string       `json:"status"`
	StartedAt *time.Time   `json:"started_at,omitempty"`
	EndedAt   *time.Time   `json:"ended_at,omitempty"`
	StartLoc  *LocationDTO `json:"start_location,omitempty"`
	EndLoc    *LocationDTO `json:"end_location,omitempty"`
}

// CheckinRequest defines the input for a student check-in.
type CheckinRequest struct {
	StudentID string       `json:"student_id" validate:"required,uuid"`
	Action    string       `json:"action" validate:"required,oneof=board exit school_recv"`
	Location  *LocationDTO `json:"location,omitempty"`
}

// CheckinResponse represents the public data for a check-in event.
type CheckinResponse struct {
	ID        string       `json:"id"`
	TripID    string       `json:"trip_id"`
	StudentID string       `json:"student_id"`
	Action    string       `json:"action"`
	Location  *LocationDTO `json:"location,omitempty"`
	Timestamp time.Time    `json:"timestamp"`
}

// ListTripsResponse represents a paginated list of trips.
type ListTripsResponse struct {
	Trips []*TripResponse `json:"trips"`
	Total int32           `json:"total"`
}
