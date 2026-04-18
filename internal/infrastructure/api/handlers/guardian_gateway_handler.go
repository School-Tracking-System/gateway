package handlers

import (
	"encoding/json"
	"net/http"

	pb "github.com/fercho/school-tracking/proto/gen/fleet/v1"
	"github.com/fercho/school-tracking/services/gateway/internal/infrastructure/api/dtos"
	"github.com/go-chi/chi/v5"
)

// LinkGuardian godoc
// @Summary      Link a guardian to a student
// @Description  Create a relationship between a user (guardian role) and a student
// @Tags         fleet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dtos.LinkGuardianRequest true "Link Guardian Request"
// @Success      201 {object} dtos.GuardianResponse
// @Failure      400 {object} dtos.ErrorResponse
// @Router       /fleet/guardians [post]
func (h *FleetHandler) LinkGuardian(w http.ResponseWriter, r *http.Request) {
	var body dtos.LinkGuardianRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.guardians.LinkGuardian(r.Context(), &pb.LinkGuardianRequest{
		UserId:    body.UserID,
		StudentId: body.StudentID,
		Relation:  body.Relation,
		IsPrimary: body.IsPrimary,
	})
	if err != nil {
		h.handleError(w, "LinkGuardian", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mapGuardianToResponse(res.Guardian))
}

// UnlinkGuardian godoc
// @Summary      Unlink a guardian
// @Description  Remove the relationship between a guardian and a student
// @Tags         fleet
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Guardian Relationship UUID"
// @Success      200 {object} map[string]string
// @Router       /fleet/guardians/{id} [delete]
func (h *FleetHandler) UnlinkGuardian(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.guardians.UnlinkGuardian(r.Context(), &pb.UnlinkGuardianRequest{Id: id})
	if err != nil {
		h.handleError(w, "UnlinkGuardian", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "guardian unlinked"})
}

// GetGuardiansByStudent godoc
// @Summary      Get student guardians
// @Description  List all guardians linked to a specific student
// @Tags         fleet
// @Produce      json
// @Security     BearerAuth
// @Param        student_id path string true "Student UUID"
// @Success      200 {array} dtos.GuardianResponse
// @Router       /fleet/students/{student_id}/guardians [get]
func (h *FleetHandler) GetGuardiansByStudent(w http.ResponseWriter, r *http.Request) {
	studentID := chi.URLParam(r, "student_id")
	res, err := h.guardians.GetGuardiansByStudent(r.Context(), &pb.GetGuardiansByStudentRequest{StudentId: studentID})
	if err != nil {
		h.handleError(w, "GetGuardiansByStudent", err)
		return
	}

	resp := make([]*dtos.GuardianResponse, len(res.Guardians))
	for i, g := range res.Guardians {
		resp[i] = mapGuardianToResponse(g)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// --- Mapper ---

func mapGuardianToResponse(g *pb.Guardian) *dtos.GuardianResponse {
	if g == nil {
		return nil
	}
	return &dtos.GuardianResponse{
		ID:        g.Id,
		UserID:    g.UserId,
		StudentID: g.StudentId,
		Relation:  g.Relation,
		IsPrimary: g.IsPrimary,
	}
}
