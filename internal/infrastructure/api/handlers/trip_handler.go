package handlers

import (
	"encoding/json"
	"net/http"

	pb "github.com/fercho/school-tracking/proto/gen/trip/v1"
	"github.com/fercho/school-tracking/services/gateway/internal/infrastructure/api/dtos"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TripHandler handles all HTTP requests for trip entities, delegating to
// the appropriate gRPC client on the Trip service.
type TripHandler struct {
	client pb.TripServiceClient
	log    *zap.Logger
}

func NewTripHandler(client pb.TripServiceClient, log *zap.Logger) *TripHandler {
	return &TripHandler{
		client: client,
		log:    log,
	}
}

// StartTrip godoc
// @Summary      Start a new trip
// @Description  Create a new trip instance based on a route ID
// @Tags         trips
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dtos.StartTripRequest true "Start Trip Request"
// @Success      201 {object} dtos.TripResponse
// @Failure      400 {object} dtos.ErrorResponse
// @Router       /trips [post]
func (h *TripHandler) StartTrip(w http.ResponseWriter, r *http.Request) {
	var body dtos.StartTripRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.log.Error("Failed to decode JSON request", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.client.StartTrip(r.Context(), &pb.StartTripRequest{
		RouteId:   body.RouteID,
		VehicleId: body.VehicleID,
		DriverId:  body.DriverID,
	})
	if err != nil {
		h.handleError(w, "StartTrip", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mapTripToResponse(res.Trip))
}

// EndTrip godoc
// @Summary      End a trip
// @Description  Mark an in-progress trip as completed
// @Tags         trips
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Trip UUID"
// @Success      200 {object} dtos.TripResponse
// @Failure      400 {object} dtos.ErrorResponse
// @Failure      404 {object} dtos.ErrorResponse
// @Router       /trips/{id}/end [put]
func (h *TripHandler) EndTrip(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing trip ID", http.StatusBadRequest)
		return
	}

	res, err := h.client.EndTrip(r.Context(), &pb.EndTripRequest{Id: id})
	if err != nil {
		h.handleError(w, "EndTrip", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapTripToResponse(res.Trip))
}

// GetTrip godoc
// @Summary      Get trip details
// @Description  Retrieve full trip data by UUID
// @Tags         trips
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Trip UUID"
// @Success      200 {object} dtos.TripResponse
// @Failure      404 {object} dtos.ErrorResponse
// @Router       /trips/{id} [get]
func (h *TripHandler) GetTrip(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing trip ID", http.StatusBadRequest)
		return
	}

	res, err := h.client.GetTrip(r.Context(), &pb.GetTripRequest{Id: id})
	if err != nil {
		h.handleError(w, "GetTrip", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapTripToResponse(res.Trip))
}

// ListTrips godoc
// @Summary      List trips
// @Description  Paginated list of trips, filterable by status or route
// @Tags         trips
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "Limit results (default 10)"
// @Param        offset query int false "Offset results (default 0)"
// @Param        status query string false "Filter by status"
// @Param        route_id query string false "Filter by Route ID"
// @Success      200 {object} dtos.ListTripsResponse
// @Router       /trips [get]
func (h *TripHandler) ListTrips(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagination(r)

	req := &pb.ListTripsRequest{
		Limit:  limit,
		Offset: offset,
	}

	if status := r.URL.Query().Get("status"); status != "" {
		req.Status = status
	}
	if routeID := r.URL.Query().Get("route_id"); routeID != "" {
		req.RouteId = routeID
	}

	res, err := h.client.ListTrips(r.Context(), req)
	if err != nil {
		h.handleError(w, "ListTrips", err)
		return
	}

	resp := dtos.ListTripsResponse{
		Total: res.TotalCount,
		Trips: make([]*dtos.TripResponse, len(res.Trips)),
	}

	for i, t := range res.Trips {
		resp.Trips[i] = mapTripToResponse(t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// ListActiveTrips godoc
// @Summary      List active trips
// @Description  Retrieve all trips currently in progress
// @Tags         trips
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} dtos.TripResponse
// @Router       /trips/active [get]
func (h *TripHandler) ListActiveTrips(w http.ResponseWriter, r *http.Request) {
	res, err := h.client.ListActiveTrips(r.Context(), &pb.ListActiveTripsRequest{})
	if err != nil {
		h.handleError(w, "ListActiveTrips", err)
		return
	}

	resp := make([]*dtos.TripResponse, len(res.Trips))
	for i, t := range res.Trips {
		resp[i] = mapTripToResponse(t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// CheckinStudent godoc
// @Summary      Student check-in
// @Description  Register a boarding/exiting event for a student on a trip
// @Tags         trips
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Trip UUID"
// @Param        request body dtos.CheckinRequest true "Checkin Data"
// @Success      201 {object} dtos.CheckinResponse
// @Failure      400 {object} dtos.ErrorResponse
// @Router       /trips/{id}/checkins [post]
func (h *TripHandler) CheckinStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing trip ID", http.StatusBadRequest)
		return
	}

	var body dtos.CheckinRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.log.Error("Failed to decode JSON request", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req := &pb.CheckinStudentRequest{
		TripId:    id,
		StudentId: body.StudentID,
		Action:    body.Action,
	}

	if body.Location != nil {
		req.Latitude = body.Location.Latitude
		req.Longitude = body.Location.Longitude
	}

	res, err := h.client.CheckinStudent(r.Context(), req)
	if err != nil {
		h.handleError(w, "CheckinStudent", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mapCheckinToResponse(res.Checkin))
}

// GetTripCheckins godoc
// @Summary      Get trip check-ins
// @Description  Retrieve all student check-in events for a specific trip
// @Tags         trips
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Trip UUID"
// @Success      200 {array} dtos.CheckinResponse
// @Router       /trips/{id}/checkins [get]
func (h *TripHandler) GetTripCheckins(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing trip ID", http.StatusBadRequest)
		return
	}

	res, err := h.client.GetTripCheckins(r.Context(), &pb.GetTripCheckinsRequest{TripId: id})
	if err != nil {
		h.handleError(w, "GetTripCheckins", err)
		return
	}

	resp := make([]*dtos.CheckinResponse, len(res.Checkins))
	for i, c := range res.Checkins {
		resp[i] = mapCheckinToResponse(c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// --- Mappers ---

func mapTripToResponse(t *pb.Trip) *dtos.TripResponse {
	if t == nil {
		return nil
	}

	resp := &dtos.TripResponse{
		ID:        t.Id,
		RouteID:   t.RouteId,
		DriverID:  t.DriverId,
		VehicleID: t.VehicleId,
		Status:    t.Status,
	}

	if t.StartedAt != nil {
		pt := t.StartedAt.AsTime()
		resp.StartedAt = &pt
	}
	if t.EndedAt != nil {
		pt := t.EndedAt.AsTime()
		resp.EndedAt = &pt
	}

	return resp
}

func mapCheckinToResponse(c *pb.TripCheckin) *dtos.CheckinResponse {
	if c == nil {
		return nil
	}

	resp := &dtos.CheckinResponse{
		ID:        c.Id,
		TripID:    c.TripId,
		StudentID: c.StudentId,
		Action:    c.Action,
		Location: &dtos.LocationDTO{
			Latitude:  c.Latitude,
			Longitude: c.Longitude,
		},
	}

	if c.Timestamp != nil {
		resp.Timestamp = c.Timestamp.AsTime()
	}

	return resp
}

func (h *TripHandler) handleError(w http.ResponseWriter, operation string, err error) {
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
