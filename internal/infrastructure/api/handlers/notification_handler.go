package handlers

import (
	"encoding/json"
	"net/http"

	pb "github.com/fercho/school-tracking/proto/gen/notification/v1"
	"github.com/fercho/school-tracking/services/gateway/internal/infrastructure/api/dtos"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NotificationHandler handles HTTP requests for notifications,
// delegating to the Notification gRPC service.
type NotificationHandler struct {
	client pb.NotificationServiceClient
	log    *zap.Logger
}

func NewNotificationHandler(client pb.NotificationServiceClient, log *zap.Logger) *NotificationHandler {
	return &NotificationHandler{
		client: client,
		log:    log,
	}
}

// SendPush godoc
// @Summary      Send a push notification
// @Description  Enqueue a new push notification via FCM
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dtos.SendPushRequest true "Push Request"
// @Success      201 {object} dtos.NotificationResponse
// @Failure      400 {object} dtos.ErrorResponse
// @Router       /notifications/push [post]
func (h *NotificationHandler) SendPush(w http.ResponseWriter, r *http.Request) {
	var body dtos.SendPushRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.log.Error("Failed to decode JSON request", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	dataJSON, _ := json.Marshal(body.Data)

	res, err := h.client.SendPush(r.Context(), &pb.SendPushRequest{
		UserId: body.UserID,
		Title:  body.Title,
		Body:   body.Body,
		Data:   string(dataJSON),
	})
	if err != nil {
		h.handleError(w, "SendPush", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mapNotificationToResponse(res.Notification))
}

// SendSMS godoc
// @Summary      Send an SMS notification
// @Description  Enqueue a new SMS via Twilio
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dtos.SendSMSRequest true "SMS Request"
// @Success      201 {object} dtos.NotificationResponse
// @Failure      400 {object} dtos.ErrorResponse
// @Router       /notifications/sms [post]
func (h *NotificationHandler) SendSMS(w http.ResponseWriter, r *http.Request) {
	var body dtos.SendSMSRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.log.Error("Failed to decode JSON request", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.client.SendSMS(r.Context(), &pb.SendSMSRequest{
		UserId: body.UserID,
		Body:   body.Body,
	})
	if err != nil {
		h.handleError(w, "SendSMS", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mapNotificationToResponse(res.Notification))
}

// GetNotification godoc
// @Summary      Get notification details
// @Description  Retrieve a notification log entry by UUID
// @Tags         notifications
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Notification UUID"
// @Success      200 {object} dtos.NotificationResponse
// @Failure      404 {object} dtos.ErrorResponse
// @Router       /notifications/{id} [get]
func (h *NotificationHandler) GetNotification(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing notification ID", http.StatusBadRequest)
		return
	}

	res, err := h.client.GetNotification(r.Context(), &pb.GetNotificationRequest{Id: id})
	if err != nil {
		h.handleError(w, "GetNotification", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapNotificationToResponse(res.Notification))
}

// ListNotifications godoc
// @Summary      List notifications
// @Description  Retrieve a paginated list of notification logs
// @Tags         notifications
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "Limit (default 10)"
// @Param        offset query int false "Offset (default 0)"
// @Param        user_id query string false "Filter by User ID"
// @Param        status query string false "Filter by Status"
// @Success      200 {object} dtos.ListNotificationsResponse
// @Router       /notifications [get]
func (h *NotificationHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagination(r)

	req := &pb.ListNotificationsRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	if userID := r.URL.Query().Get("user_id"); userID != "" {
		req.UserId = userID
	}
	if status := r.URL.Query().Get("status"); status != "" {
		req.Status = status
	}

	res, err := h.client.ListNotifications(r.Context(), req)
	if err != nil {
		h.handleError(w, "ListNotifications", err)
		return
	}

	resp := dtos.ListNotificationsResponse{
		Total:         res.TotalCount,
		Notifications: make([]*dtos.NotificationResponse, len(res.Notifications)),
	}

	for i, n := range res.Notifications {
		resp.Notifications[i] = mapNotificationToResponse(n)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RetryFailed godoc
// @Summary      Retry failed notifications
// @Description  Manually trigger a retry for all notifications in 'failed' status
// @Tags         notifications
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{}
// @Router       /notifications/retry [post]
func (h *NotificationHandler) RetryFailed(w http.ResponseWriter, r *http.Request) {
	res, err := h.client.RetryFailed(r.Context(), &pb.RetryFailedRequest{})
	if err != nil {
		h.handleError(w, "RetryFailed", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// --- Mappers ---

func mapNotificationToResponse(n *pb.Notification) *dtos.NotificationResponse {
	if n == nil {
		return nil
	}

	resp := &dtos.NotificationResponse{
		ID:      n.Id,
		UserID:  n.UserId,
		Type:    n.Type,
		Channel: n.Channel,
		Title:   n.Title,
		Body:    n.Body,
		Status:  n.Status,
	}

	if n.CreatedAt != nil {
		resp.CreatedAt = n.CreatedAt.AsTime()
	}
	if n.UpdatedAt != nil {
		t := n.UpdatedAt.AsTime()
		resp.SentAt = &t
	}

	return resp
}

func (h *NotificationHandler) handleError(w http.ResponseWriter, operation string, err error) {
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
