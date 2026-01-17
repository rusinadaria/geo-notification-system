package models

import (
	"time"
)

type WebhookPayload struct {
	Event     string                   `json:"event"`
	UserID    int                      `json:"user_id"`
	Lat       float64                  `json:"lat"`
	Lon       float64                  `json:"lon"`
	Incidents []NearbyIncidentResponse `json:"incidents"`
	CheckedAt time.Time                `json:"checked_at"`
}

type LocationCheckResponse struct {
	Danger    bool                     `json:"danger"`
	Incidents []NearbyIncidentResponse `json:"incidents"`
}

type NearbyIncidentResponse struct {
	ID             int64   `json:"id"`
	Type           string  `json:"type"`
	DistanceMeters float64 `json:"distance_meters"`
}

type LocationCheckRequest struct {
	UserID int     `json:"user_id"`
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
}

type Incident struct {
	ID          int64     `json:"id" db:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description" db:"description"`
	Active      bool      `json:"active" db:"active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type IncidentRequest struct {
	Type         string  `json:"type"`
	Description  string  `json:"description"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	RadiusMeters int     `json:"radius_meters"`
	Active       bool    `json:"active"`
}

type IncidentResponse struct {
	Type         string    `json:"type" db:"type"`
	Description  string    `json:"description" db:"description"`
	Latitude     float64   `json:"latitude" db:"latitude"`
	Longitude    float64   `json:"longitude" db:"longitude"`
	RadiusMeters int       `json:"radius_meters" db:"radius_meters"`
	Active       bool      `json:"active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type IncidentStatsResponse struct {
	UserCount    int64 `json:"user_count"`
	WindowMinute int   `json:"window_minutes"`
}

type HealthStatus string

const (
	HealthOK       HealthStatus = "ok"
	HealthDown     HealthStatus = "down"
	HealthDegraded HealthStatus = "degraded"
)

type HealthResponse struct {
	Status    HealthStatus      `json:"status"`
	Checks    map[string]string `json:"checks"`
	Timestamp time.Time         `json:"timestamp"`
}

type ErrorResponse struct {
	Errors string `json:"errors"`
}
