package models

import (
	"time"
)

// type NearbyIncident struct {
//     ID             int64
//     Type           string
//     DistanceMeters float64
// }

type LocationCheckResponse struct { // Сервис возвращает СПИСОК БЛИЖАЙШИХ ОПАСНЫХ ЗОН (т.е. все активные инциденты, в радиусе которых user находится прямо сейчас)
    Danger    bool                    `json:"danger"`
    Incidents []NearbyIncidentResponse `json:"incidents"`
}

type NearbyIncidentResponse struct {
    ID             int64   `json:"id"`
    Type           string  `json:"type"`
    DistanceMeters float64 `json:"distance_meters"`
}


type LocationCheckRequest struct { // Пользователь отправляет свои координаты
    UserID string  `json:"user_id"`
    Lat    float64 `json:"lat"`
    Lon    float64 `json:"lon"`
}

type Incident struct {
	ID          int64     `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Active      bool      `json:"active" db:"active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type IncidentRequest struct {
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Active      bool      `json:"active" db:"active"`
}

type IncidentResponse struct {
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Active      bool      `json:"active" db:"active"`
}


type ErrorResponse struct {
	Errors string `json:"errors"`
}