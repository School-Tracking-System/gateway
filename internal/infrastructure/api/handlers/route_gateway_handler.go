package handlers

import (
	"encoding/json"
	"net/http"

	pb "github.com/fercho/school-tracking/proto/gen/fleet/v1"
	"github.com/fercho/school-tracking/services/gateway/internal/infrastructure/api/dtos"
	"github.com/go-chi/chi/v5"
)

// CreateRoute godoc
// @Summary      Create a route
// @Description  Register a new school transport route
// @Tags         fleet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dtos.CreateRouteRequest true "Create Route Request"
// @Success      201 {object} dtos.RouteResponse
// @Failure      400 {object} dtos.ErrorResponse
// @Router       /fleet/routes [post]
func (h *FleetHandler) CreateRoute(w http.ResponseWriter, r *http.Request) {
	var body dtos.CreateRouteRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.routes.CreateRoute(r.Context(), &pb.CreateRouteRequest{
		Name:         body.Name,
		Direction:    body.Direction,
		ScheduleTime: body.ScheduleTime,
		SchoolId:     body.SchoolID,
		VehicleId:    body.VehicleID,
		DriverId:     body.DriverID,
	})
	if err != nil {
		h.handleError(w, "CreateRoute", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mapRouteToResponse(res.Route))
}

// GetRoute godoc
// @Summary      Get route details
// @Description  Retrieve route info and its stops
// @Tags         fleet
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Route UUID"
// @Success      200 {object} dtos.RouteResponse
// @Failure      404 {object} dtos.ErrorResponse
// @Router       /fleet/routes/{id} [get]
func (h *FleetHandler) GetRoute(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing route ID", http.StatusBadRequest)
		return
	}

	res, err := h.routes.GetRoute(r.Context(), &pb.GetRouteRequest{Id: id})
	if err != nil {
		h.handleError(w, "GetRoute", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapRouteToResponse(res.Route))
}

// UpdateRoute godoc
// @Summary      Update a route
// @Description  Modify route basic info or assignments
// @Tags         fleet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Route UUID"
// @Param        request body dtos.UpdateRouteRequest true "Update Route Request"
// @Success      200 {object} dtos.RouteResponse
// @Failure      400 {object} dtos.ErrorResponse
// @Router       /fleet/routes/{id} [put]
func (h *FleetHandler) UpdateRoute(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body dtos.UpdateRouteRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.routes.UpdateRoute(r.Context(), &pb.UpdateRouteRequest{
		Id:           id,
		Name:         body.Name,
		Direction:    body.Direction,
		ScheduleTime: body.ScheduleTime,
		VehicleId:    body.VehicleID,
		DriverId:     body.DriverID,
	})
	if err != nil {
		h.handleError(w, "UpdateRoute", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapRouteToResponse(res.Route))
}

// ListRoutes godoc
// @Summary      List routes
// @Description  Retrieve a paginated list of routes
// @Tags         fleet
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "Limit (default 10)"
// @Param        offset query int false "Offset (default 0)"
// @Success      200 {object} dtos.ListRoutesResponse
// @Router       /fleet/routes [get]
func (h *FleetHandler) ListRoutes(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagination(r)
	res, err := h.routes.ListRoutes(r.Context(), &pb.ListRoutesRequest{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		h.handleError(w, "ListRoutes", err)
		return
	}

	resp := dtos.ListRoutesResponse{
		Total:  res.TotalCount,
		Routes: make([]*dtos.RouteResponse, len(res.Routes)),
	}
	for i, route := range res.Routes {
		resp.Routes[i] = mapRouteToResponse(route)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// AddStop godoc
// @Summary      Add a stop to a route
// @Description  Create a new stop for a student in a route
// @Tags         fleet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Route UUID"
// @Param        request body dtos.AddStopRequest true "Add Stop Request"
// @Success      201 {object} dtos.StopDTO
// @Failure      400 {object} dtos.ErrorResponse
// @Router       /fleet/routes/{id}/stops [post]
func (h *FleetHandler) AddStop(w http.ResponseWriter, r *http.Request) {
	routeID := chi.URLParam(r, "id")
	var body dtos.AddStopRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req := &pb.AddStopRequest{
		RouteId:   routeID,
		StudentId: body.StudentID,
		Order:     body.Order,
		Address:   body.Address,
		EstTime:   body.EstTime,
	}

	if body.Location != nil {
		req.Location = &pb.GeoPoint{
			Latitude:  body.Location.Latitude,
			Longitude: body.Location.Longitude,
		}
	}

	res, err := h.routes.AddStop(r.Context(), req)
	if err != nil {
		h.handleError(w, "AddStop", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mapStopToResponse(res.Stop))
}

// GetRouteStops godoc
// @Summary      Get route stops
// @Description  Retrieve all stops for a specific route
// @Tags         fleet
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Route UUID"
// @Success      200 {array} dtos.StopDTO
// @Router       /fleet/routes/{id}/stops [get]
func (h *FleetHandler) GetRouteStops(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	res, err := h.routes.GetRouteStops(r.Context(), &pb.GetRouteStopsRequest{RouteId: id})
	if err != nil {
		h.handleError(w, "ListStops", err)
		return
	}

	resp := make([]*dtos.StopDTO, len(res.Stops))
	for i, s := range res.Stops {
		resp[i] = mapStopToResponse(s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RemoveStop godoc
// @Summary      Remove a stop
// @Description  Delete a stop from a route
// @Tags         fleet
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Route UUID"
// @Param        stop_id path string true "Stop UUID"
// @Success      200 {object} map[string]string
// @Router       /fleet/routes/{id}/stops/{stop_id} [delete]
func (h *FleetHandler) RemoveStop(w http.ResponseWriter, r *http.Request) {
	stopID := chi.URLParam(r, "stop_id")
	_, err := h.routes.RemoveStop(r.Context(), &pb.RemoveStopRequest{Id: stopID})
	if err != nil {
		h.handleError(w, "RemoveStop", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "stop removed"})
}

// --- Mappers ---

func mapRouteToResponse(r *pb.Route) *dtos.RouteResponse {
	if r == nil {
		return nil
	}
	resp := &dtos.RouteResponse{
		ID:           r.Id,
		Name:         r.Name,
		Direction:    r.Direction,
		ScheduleTime: r.ScheduleTime,
		VehicleID:    r.VehicleId,
		DriverID:     r.DriverId,
		IsActive:     r.IsActive,
	}

	return resp
}

func mapStopToResponse(s *pb.RouteStop) *dtos.StopDTO {
	if s == nil {
		return nil
	}
	resp := &dtos.StopDTO{
		ID:        s.Id,
		StudentID: s.StudentId,
		Order:     s.Order,
		Address:   s.Address,
		EstTime:   s.EstTime,
	}
	if s.Location != nil {
		resp.Location = &dtos.LocationDTO{
			Latitude:  s.Location.Latitude,
			Longitude: s.Location.Longitude,
		}
	}
	return resp
}
