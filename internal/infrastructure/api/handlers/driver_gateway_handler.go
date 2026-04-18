package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	pb "github.com/fercho/school-tracking/proto/gen/fleet/v1"
	"github.com/fercho/school-tracking/services/gateway/internal/infrastructure/api/dtos"
	"github.com/go-chi/chi/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// RegisterDriver godoc
// @Summary      Register a driver
// @Description  Register a new driver linked to an Auth service user
// @Tags         fleet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dtos.RegisterDriverRequest true "Register Driver Request"
// @Success      201 {object} dtos.DriverResponse
// @Failure      400 {object} dtos.ErrorResponse
// @Router       /fleet/drivers [post]
func (h *FleetHandler) RegisterDriver(w http.ResponseWriter, r *http.Request) {
	var body dtos.RegisterDriverRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req := &pb.RegisterDriverRequest{
		UserId:         body.UserID,
		LicenseNumber:  body.LicenseNumber,
		LicenseType:    body.LicenseType,
		CedulaId:       body.CedulaID,
		EmergencyPhone: body.EmergencyPhone,
	}

	if body.LicenseExpiry != "" {
		if t, err := time.Parse("2006-01-02", body.LicenseExpiry); err == nil {
			req.LicenseExpiry = timestamppb.New(t)
		}
	}

	res, err := h.drivers.RegisterDriver(r.Context(), req)
	if err != nil {
		h.handleError(w, "RegisterDriver", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mapDriverToResponse(res.Driver))
}

// GetDriver godoc
// @Summary      Get a driver by ID
// @Description  Retrieve driver details
// @Tags         fleet
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Driver UUID"
// @Success      200 {object} dtos.DriverResponse
// @Failure      404 {object} dtos.ErrorResponse
// @Router       /fleet/drivers/{id} [get]
func (h *FleetHandler) GetDriver(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing driver ID", http.StatusBadRequest)
		return
	}

	res, err := h.drivers.GetDriver(r.Context(), &pb.GetDriverRequest{Id: id})
	if err != nil {
		h.handleError(w, "GetDriver", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapDriverToResponse(res.Driver))
}

// UpdateDriver godoc
// @Summary      Update a driver
// @Description  Modify driver details
// @Tags         fleet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Driver UUID"
// @Param        request body dtos.UpdateDriverRequest true "Update Driver Request"
// @Success      200 {object} dtos.DriverResponse
// @Failure      400 {object} dtos.ErrorResponse
// @Router       /fleet/drivers/{id} [put]
func (h *FleetHandler) UpdateDriver(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body dtos.UpdateDriverRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req := &pb.UpdateDriverRequest{
		Id:             id,
		LicenseType:    body.LicenseType,
		EmergencyPhone: body.EmergencyPhone,
		Status:         body.Status,
	}
	if body.LicenseExpiry != "" {
		if t, err := time.Parse("2006-01-02", body.LicenseExpiry); err == nil {
			req.LicenseExpiry = timestamppb.New(t)
		}
	}

	res, err := h.drivers.UpdateDriver(r.Context(), req)
	if err != nil {
		h.handleError(w, "UpdateDriver", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapDriverToResponse(res.Driver))
}

// ListDrivers godoc
// @Summary      List drivers
// @Description  Retrieve a paginated list of drivers
// @Tags         fleet
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "Limit (default 10)"
// @Param        offset query int false "Offset (default 0)"
// @Success      200 {object} dtos.ListDriversResponse
// @Router       /fleet/drivers [get]
func (h *FleetHandler) ListDrivers(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagination(r)
	res, err := h.drivers.ListDrivers(r.Context(), &pb.ListDriversRequest{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		h.handleError(w, "ListDrivers", err)
		return
	}

	resp := dtos.ListDriversResponse{
		Total:   res.TotalCount,
		Drivers: make([]*dtos.DriverResponse, len(res.Drivers)),
	}
	for i, d := range res.Drivers {
		resp.Drivers[i] = mapDriverToResponse(d)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// --- Mapper ---

func mapDriverToResponse(d *pb.Driver) *dtos.DriverResponse {
	if d == nil {
		return nil
	}
	resp := &dtos.DriverResponse{
		ID:            d.Id,
		UserID:        d.UserId,
		LicenseNumber: d.LicenseNumber,
		LicenseType:   d.LicenseType,
		Status:        d.Status,
	}
	if d.LicenseExpiry != nil {
		t := d.LicenseExpiry.AsTime()
		resp.LicenseExpiry = &t
	}
	return resp
}
