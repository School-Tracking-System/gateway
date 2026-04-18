package handlers

import (
	"encoding/json"
	"net/http"

	pb "github.com/fercho/school-tracking/proto/gen/fleet/v1"
	"github.com/fercho/school-tracking/services/gateway/internal/infrastructure/api/dtos"
	"github.com/go-chi/chi/v5"
)

// CreateSchool godoc
// @Summary      Create a new school
// @Description  Register a new school in the system
// @Tags         fleet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dtos.CreateSchoolRequest true "Create School Request"
// @Success      201 {object} dtos.SchoolResponse
// @Failure      400 {object} dtos.ErrorResponse
// @Router       /fleet/schools [post]
func (h *FleetHandler) CreateSchool(w http.ResponseWriter, r *http.Request) {
	var body dtos.CreateSchoolRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req := &pb.CreateSchoolRequest{
		Name:    body.Name,
		Address: body.Address,
		Phone:   body.Phone,
		Email:   body.Email,
	}

	if body.Location != nil {
		req.Location = &pb.GeoPoint{
			Latitude:  body.Location.Latitude,
			Longitude: body.Location.Longitude,
		}
	}

	res, err := h.schools.CreateSchool(r.Context(), req)
	if err != nil {
		h.handleError(w, "CreateSchool", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mapSchoolToResponse(res.School))
}

// GetSchool godoc
// @Summary      Get a school by ID
// @Description  Retrieve school details by its UUID
// @Tags         fleet
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "School UUID"
// @Success      200 {object} dtos.SchoolResponse
// @Failure      404 {object} dtos.ErrorResponse
// @Router       /fleet/schools/{id} [get]
func (h *FleetHandler) GetSchool(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing school ID", http.StatusBadRequest)
		return
	}

	res, err := h.schools.GetSchool(r.Context(), &pb.GetSchoolRequest{Id: id})
	if err != nil {
		h.handleError(w, "GetSchool", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapSchoolToResponse(res.School))
}

// UpdateSchool godoc
// @Summary      Update a school
// @Description  Modify school details
// @Tags         fleet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "School UUID"
// @Param        request body dtos.UpdateSchoolRequest true "Update School Request"
// @Success      200 {object} dtos.SchoolResponse
// @Failure      400 {object} dtos.ErrorResponse
// @Router       /fleet/schools/{id} [put]
func (h *FleetHandler) UpdateSchool(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body dtos.UpdateSchoolRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req := &pb.UpdateSchoolRequest{
		Id:      id,
		Name:    body.Name,
		Address: body.Address,
		Phone:   body.Phone,
		Email:   body.Email,
	}

	if body.Location != nil {
		req.Location = &pb.GeoPoint{
			Latitude:  body.Location.Latitude,
			Longitude: body.Location.Longitude,
		}
	}

	res, err := h.schools.UpdateSchool(r.Context(), req)
	if err != nil {
		h.handleError(w, "UpdateSchool", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapSchoolToResponse(res.School))
}

// ListSchools godoc
// @Summary      List schools
// @Description  Retrieve a paginated list of schools
// @Tags         fleet
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "Limit (default 10)"
// @Param        offset query int false "Offset (default 0)"
// @Success      200 {object} dtos.ListSchoolsResponse
// @Router       /fleet/schools [get]
func (h *FleetHandler) ListSchools(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagination(r)
	res, err := h.schools.ListSchools(r.Context(), &pb.ListSchoolsRequest{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		h.handleError(w, "ListSchools", err)
		return
	}

	resp := dtos.ListSchoolsResponse{
		Total:   res.TotalCount,
		Schools: make([]*dtos.SchoolResponse, len(res.Schools)),
	}
	for i, s := range res.Schools {
		resp.Schools[i] = mapSchoolToResponse(s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// --- Mapper ---

func mapSchoolToResponse(s *pb.School) *dtos.SchoolResponse {
	if s == nil {
		return nil
	}
	resp := &dtos.SchoolResponse{
		ID:      s.Id,
		Name:    s.Name,
		Address: s.Address,
		Phone:   s.Phone,
		Email:   s.Email,
	}
	if s.Location != nil {
		resp.Location = &dtos.LocationDTO{
			Latitude:  s.Location.Latitude,
			Longitude: s.Location.Longitude,
		}
	}
	if s.CreatedAt != nil {
		resp.CreatedAt = s.CreatedAt.AsTime()
	}
	return resp
}
