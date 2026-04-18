package dtos

import "time"

type NotificationResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	Channel   string    `json:"channel"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Status    string    `json:"status"`
	SentAt    *time.Time `json:"sent_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type SendPushRequest struct {
	UserID string            `json:"user_id" validate:"required,uuid"`
	Title  string            `json:"title" validate:"required"`
	Body   string            `json:"body" validate:"required"`
	Data   map[string]string `json:"data,omitempty"`
}

type SendSMSRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	Body   string `json:"body" validate:"required"`
}

type ListNotificationsResponse struct {
	Notifications []*NotificationResponse `json:"notifications"`
	Total         int32                   `json:"total"`
}
