package handlers

import (
	"encoding/json"
	"net/http"

	pb "github.com/fercho/school-tracking/proto/gen/fleet/v1"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FleetHandler struct {
	client pb.VehicleServiceClient
	log    *zap.Logger
}

func NewFleetHandler(client pb.VehicleServiceClient, log *zap.Logger) *FleetHandler {
	return &FleetHandler{
		client: client,
		log:    log,
	}
}

func (h *FleetHandler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	var req pb.CreateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode JSON request", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.client.CreateVehicle(r.Context(), &req)
	if err != nil {
		h.handleError(w, "CreateVehicle", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res.Vehicle)
}

func (h *FleetHandler) GetVehicle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing vehicle ID", http.StatusBadRequest)
		return
	}

	res, err := h.client.GetVehicle(r.Context(), &pb.GetVehicleRequest{Id: id})
	if err != nil {
		h.handleError(w, "GetVehicle", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res.Vehicle)
}

func (h *FleetHandler) ListVehicles(w http.ResponseWriter, r *http.Request) {
	res, err := h.client.ListVehicles(r.Context(), &pb.ListVehiclesRequest{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		h.handleError(w, "ListVehicles", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *FleetHandler) handleError(w http.ResponseWriter, operation string, err error) {
	st, ok := status.FromError(err)
	if !ok {
		h.log.Error("Unexpected non-gRPC error", zap.String("op", operation), zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.log.Warn("gRPC error occurred",
		zap.String("op", operation),
		zap.String("code", st.Code().String()),
		zap.String("msg", st.Message()),
	)

	switch st.Code() {
	case codes.NotFound:
		http.Error(w, st.Message(), http.StatusNotFound)
	case codes.AlreadyExists:
		http.Error(w, st.Message(), http.StatusConflict) // 409 Conflict
	case codes.InvalidArgument:
		http.Error(w, st.Message(), http.StatusBadRequest)
	case codes.Unauthenticated:
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	default:
		http.Error(w, "downstream service error", http.StatusInternalServerError)
	}
}
