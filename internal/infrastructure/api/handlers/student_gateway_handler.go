package handlers

import (
	"encoding/json"
	"net/http"

	pb "github.com/fercho/school-tracking/proto/gen/fleet/v1"
	"github.com/fercho/school-tracking/services/gateway/internal/infrastructure/api/dtos"
	"github.com/go-chi/chi/v5"
)

// RegisterStudent godoc
// @Summary      Register a student
// @Description  Create a new student record
// @Tags         fleet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dtos.RegisterStudentRequest true "Register Student Request"
// @Success      201 {object} dtos.StudentResponse
// @Failure      400 {object} dtos.ErrorResponse
// @Router       /fleet/students [post]
func (h *FleetHandler) RegisterStudent(w http.ResponseWriter, r *http.Request) {
	var body dtos.RegisterStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req := &pb.RegisterStudentRequest{
		FirstName:     body.FirstName,
		LastName:      body.LastName,
		Grade:         body.Grade,
		SchoolId:      body.SchoolID,
		PickupAddress: body.PickupAddress,
		CedulaId:      body.CedulaID,
	}

	if body.PickupLocation != nil {
		req.PickupLocation = &pb.GeoPoint{
			Latitude:  body.PickupLocation.Latitude,
			Longitude: body.PickupLocation.Longitude,
		}
	}

	res, err := h.students.RegisterStudent(r.Context(), req)
	if err != nil {
		h.handleError(w, "RegisterStudent", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mapStudentToResponse(res.Student))
}

// GetStudent godoc
// @Summary      Get a student by ID
// @Description  Retrieve student details by UUID
// @Tags         fleet
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Student UUID"
// @Success      200 {object} dtos.StudentResponse
// @Failure      404 {object} dtos.ErrorResponse
// @Router       /fleet/students/{id} [get]
func (h *FleetHandler) GetStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing student ID", http.StatusBadRequest)
		return
	}

	res, err := h.students.GetStudent(r.Context(), &pb.GetStudentRequest{Id: id})
	if err != nil {
		h.handleError(w, "GetStudent", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapStudentToResponse(res.Student))
}

// UpdateStudent godoc
// @Summary      Update a student
// @Description  Modify student details
// @Tags         fleet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Student UUID"
// @Param        request body dtos.UpdateStudentRequest true "Update Student Request"
// @Success      200 {object} dtos.StudentResponse
// @Failure      400 {object} dtos.ErrorResponse
// @Router       /fleet/students/{id} [put]
func (h *FleetHandler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body dtos.UpdateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req := &pb.UpdateStudentRequest{
		Id:            id,
		FirstName:     body.FirstName,
		LastName:      body.LastName,
		Grade:         body.Grade,
		PickupAddress: body.PickupAddress,
		CedulaId:      body.CedulaID,
	}

	if body.PickupLocation != nil {
		req.PickupLocation = &pb.GeoPoint{
			Latitude:  body.PickupLocation.Latitude,
			Longitude: body.PickupLocation.Longitude,
		}
	}

	res, err := h.students.UpdateStudent(r.Context(), req)
	if err != nil {
		h.handleError(w, "UpdateStudent", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapStudentToResponse(res.Student))
}

// ListStudents godoc
// @Summary      List students
// @Description  Retrieve a paginated list of students
// @Tags         fleet
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "Limit (default 10)"
// @Param        offset query int false "Offset (default 0)"
// @Success      200 {object} dtos.ListStudentsResponse
// @Router       /fleet/students [get]
func (h *FleetHandler) ListStudents(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagination(r)
	res, err := h.students.ListStudents(r.Context(), &pb.ListStudentsRequest{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		h.handleError(w, "ListStudents", err)
		return
	}

	resp := dtos.ListStudentsResponse{
		Total:    res.TotalCount,
		Students: make([]*dtos.StudentResponse, len(res.Students)),
	}
	for i, s := range res.Students {
		resp.Students[i] = mapStudentToResponse(s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeactivateStudent godoc
// @Summary      Deactivate a student
// @Description  Mark a student as inactive
// @Tags         fleet
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Student UUID"
// @Success      200 {object} map[string]string
// @Router       /fleet/students/{id} [delete]
func (h *FleetHandler) DeactivateStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.students.DeactivateStudent(r.Context(), &pb.DeactivateStudentRequest{Id: id})
	if err != nil {
		h.handleError(w, "DeactivateStudent", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "student deactivated"})
}

// --- Mapper ---

func mapStudentToResponse(s *pb.Student) *dtos.StudentResponse {
	if s == nil {
		return nil
	}
	resp := &dtos.StudentResponse{
		ID:            s.Id,
		FirstName:     s.FirstName,
		LastName:      s.LastName,
		Grade:         s.Grade,
		SchoolID:      s.SchoolId,
		PickupAddress: s.PickupAddress,
		IsActive:      s.IsActive,
		CedulaID:      s.CedulaId,
	}
	if s.PickupLocation != nil {
		resp.PickupLocation = &dtos.LocationDTO{
			Latitude:  s.PickupLocation.Latitude,
			Longitude: s.PickupLocation.Longitude,
		}
	}
	return resp
}
