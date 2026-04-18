package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	pb "github.com/fercho/school-tracking/proto/gen/fleet/v1"
	"github.com/fercho/school-tracking/services/gateway/internal/infrastructure/api/dtos"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// FleetHandler handles all HTTP requests for fleet entities, delegating to
// the appropriate gRPC client on the Fleet service.
type FleetHandler struct {
	vehicles  pb.VehicleServiceClient
	schools   pb.SchoolServiceClient
	drivers   pb.DriverServiceClient
	students  pb.StudentServiceClient
	guardians pb.GuardianServiceClient
	routes    pb.RouteServiceClient
	log       *zap.Logger
}

func NewFleetHandler(
	vehicles pb.VehicleServiceClient,
	schools pb.SchoolServiceClient,
	drivers pb.DriverServiceClient,
	students pb.StudentServiceClient,
	guardians pb.GuardianServiceClient,
	routes pb.RouteServiceClient,
	log *zap.Logger,
) *FleetHandler {
	return &FleetHandler{
		vehicles:  vehicles,
		schools:   schools,
		drivers:   drivers,
		students:  students,
		guardians: guardians,
		routes:    routes,
		log:       log,
	}
}

// parsePagination reads limit and offset from query params, defaulting to 10/0.
func parsePagination(r *http.Request) (limit, offset int32) {
	limit = 10
	offset = 0
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = int32(n)
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = int32(n)
		}
	}
	return
}

// CreateVehicle godoc
// @Summary      Create a new vehicle
// @Description  Register a new vehicle in the fleet
// @Tags         fleet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dtos.CreateVehicleRequest true "Create Vehicle Request"
// @Success      201 {object} dtos.VehicleResponse
// @Failure      400 {object} dtos.ErrorResponse
// @Router       /fleet/vehicles [post]
func (h *FleetHandler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	var body dtos.CreateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.log.Error("Failed to decode JSON request", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req := &pb.CreateVehicleRequest{
		Plate:       body.Plate,
		Brand:       body.Brand,
		Model:       body.Model,
		Year:        body.Year,
		Capacity:    body.Capacity,
		Color:       body.Color,
		VehicleType: body.VehicleType,
		ChassisNum:  body.ChassisNum,
	}

	if body.InsuranceExp != "" {
		if t, err := time.Parse("2006-01-02", body.InsuranceExp); err == nil {
			req.InsuranceExp = timestamppb.New(t)
		}
	}
	if body.TechReviewExp != "" {
		if t, err := time.Parse("2006-01-02", body.TechReviewExp); err == nil {
			req.TechReviewExp = timestamppb.New(t)
		}
	}

	res, err := h.vehicles.CreateVehicle(r.Context(), req)
	if err != nil {
		h.handleError(w, "CreateVehicle", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mapVehicleToResponse(res.Vehicle))
}

// GetVehicle godoc
// @Summary      Get a vehicle by ID
// @Description  Retrieve a single vehicle by its UUID
// @Tags         fleet
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Vehicle UUID"
// @Success      200 {object} dtos.VehicleResponse
// @Failure      404 {object} dtos.ErrorResponse
// @Router       /fleet/vehicles/{id} [get]
func (h *FleetHandler) GetVehicle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing vehicle ID", http.StatusBadRequest)
		return
	}

	res, err := h.vehicles.GetVehicle(r.Context(), &pb.GetVehicleRequest{Id: id})
	if err != nil {
		h.handleError(w, "GetVehicle", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapVehicleToResponse(res.Vehicle))
}

// ListVehicles godoc
// @Summary      List all vehicles
// @Description  Retrieve a paginated list of vehicles in the fleet
// @Tags         fleet
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "Limit (default 10)"
// @Param        offset query int false "Offset (default 0)"
// @Success      200 {object} dtos.ListVehiclesResponse
// @Router       /fleet/vehicles [get]
func (h *FleetHandler) ListVehicles(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagination(r)
	res, err := h.vehicles.ListVehicles(r.Context(), &pb.ListVehiclesRequest{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		h.handleError(w, "ListVehicles", err)
		return
	}

	resp := dtos.ListVehiclesResponse{
		Total:    res.TotalCount,
		Vehicles: make([]*dtos.VehicleResponse, len(res.Vehicles)),
	}
	for i, v := range res.Vehicles {
		resp.Vehicles[i] = mapVehicleToResponse(v)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// --- Mappers ---

func mapVehicleToResponse(v *pb.Vehicle) *dtos.VehicleResponse {
	if v == nil {
		return nil
	}

	resp := &dtos.VehicleResponse{
		ID:          v.Id,
		Plate:       v.Plate,
		Brand:       v.Brand,
		Model:       v.Model,
		Year:        v.Year,
		Capacity:    v.Capacity,
		Status:      v.Status,
		Color:       v.Color,
		VehicleType: v.VehicleType,
		ChassisNum:  v.ChassisNum,
	}

	if v.InsuranceExp != nil {
		t := v.InsuranceExp.AsTime()
		resp.InsuranceExp = &t
	}
	if v.TechReviewExp != nil {
		t := v.TechReviewExp.AsTime()
		resp.TechReviewExp = &t
	}
	if v.CreatedAt != nil {
		resp.CreatedAt = v.CreatedAt.AsTime()
	}
	if v.UpdatedAt != nil {
		resp.UpdatedAt = v.UpdatedAt.AsTime()
	}

	return resp
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
		http.Error(w, st.Message(), http.StatusConflict)
	case codes.InvalidArgument:
		http.Error(w, st.Message(), http.StatusBadRequest)
	case codes.Unauthenticated:
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	default:
		http.Error(w, "downstream service error", http.StatusInternalServerError)
	}
}
