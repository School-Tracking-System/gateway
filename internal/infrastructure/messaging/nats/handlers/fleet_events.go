package handlers

import (
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

// VehicleCreatedEvent mirrors the payload published by the fleet service.
type VehicleCreatedEvent struct {
	ID        string    `json:"id"`
	Plate     string    `json:"plate"`
	Brand     string    `json:"brand"`
	Model     string    `json:"model"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// VehicleUpdatedEvent mirrors the payload published by the fleet service.
type VehicleUpdatedEvent struct {
	ID        string    `json:"id"`
	Plate     string    `json:"plate"`
	Brand     string    `json:"brand"`
	Model     string    `json:"model"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FleetEventHandler handles fleet domain events received via NATS.
type FleetEventHandler struct {
	log *zap.Logger
}

// NewFleetEventHandler creates a new FleetEventHandler.
func NewFleetEventHandler(log *zap.Logger) *FleetEventHandler {
	return &FleetEventHandler{log: log}
}

// HandleVehicleCreated processes a fleet.vehicle.created event.
func (h *FleetEventHandler) HandleVehicleCreated(subject string, payload []byte) {
	var event VehicleCreatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		h.log.Error("Failed to unmarshal vehicle.created event", zap.Error(err))
		return
	}
	h.log.Info("fleet.vehicle.created received",
		zap.String("id", event.ID),
		zap.String("plate", event.Plate),
		zap.String("brand", event.Brand),
		zap.String("status", event.Status),
	)
}

// HandleVehicleUpdated processes a fleet.vehicle.updated event.
func (h *FleetEventHandler) HandleVehicleUpdated(subject string, payload []byte) {
	var event VehicleUpdatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		h.log.Error("Failed to unmarshal vehicle.updated event", zap.Error(err))
		return
	}
	h.log.Info("fleet.vehicle.updated received",
		zap.String("id", event.ID),
		zap.String("plate", event.Plate),
		zap.String("status", event.Status),
	)
}
