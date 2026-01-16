package services

import (
	"database/sql"
	"errors"
	"log"
	// "mime/quotedprintable"
	"context"
	"time"
	"github.com/rusinadaria/geo-notification-system/internal/models"
	// "github.com/rusinadaria/geo-notification-system/internal/queue"
	"github.com/rusinadaria/geo-notification-system/internal/repository"
)

type IncidentService struct {
	repo repository.Incident
	windowMin int	
}

func NewIncidentService(repo repository.Incident, windowMin int) *IncidentService {
	return &IncidentService{repo: repo, windowMin: windowMin}
}

var ErrIncidentAlreadyExists = errors.New("incident already exists")

func (s *IncidentService) CheckLocation(checkReq models.LocationCheckRequest) (models.LocationCheckResponse, error) {

	nearbyResp, err := s.repo.CheckLocation(checkReq)
	if err != nil {
		log.Println(err)
		return models.LocationCheckResponse{}, err
	}

	err = s.repo.SaveCheck(checkReq.UserID, checkReq.Lat, checkReq.Lon, true) // ФУНКЦИЯ СОХРАНЕНИЯ ФАКТА ПРОВЕРКИ
	if err != nil {
		log.Println(err)
		return models.LocationCheckResponse{}, err
    }

	// if nearbyResp.Danger {
	// 	err := s.webhookQueue.Enqueue(models.WebhookPayload{
	// 		Event:     "danger_detected",
	// 		UserID:    checkReq.UserID,
	// 		Lat:       checkReq.Lat,
	// 		Lon:       checkReq.Lon,
	// 		Incidents: nearbyResp.Incidents,
	// 		CheckedAt: time.Now().UTC(),
	// 	})

	// 	if err != nil {
	// 		log.Println("failed to enqueue webhook:", err)
	// 	}
	// }

	return nearbyResp, nil
}

func (s *IncidentService) CreateIncident(req models.IncidentRequest) error {
	// id, err := s.repo.GetIncident(req)
	// if id != 0 {
	// 	return ErrIncidentAlreadyExists
	// }

	err := s.repo.CreateIncident(req)
	if err != nil {
		return err
	}

	return nil
}

func (s *IncidentService) GetAllIncidents(limit, offset int) ([]models.IncidentResponse, error) {
	incidents, err := s.repo.GetAllIncidents(limit, offset)
	if err != nil {
		return []models.IncidentResponse{}, err
	}

	return incidents, nil
}

func (s *IncidentService) GetIncidentById(id int) (models.IncidentResponse, error) {
	incident, err := s.repo.GetIncidentById(id)
	if err != nil {
		return models.IncidentResponse{}, err
	}

	return incident, nil
}

func (s *IncidentService) UpdateIncident(id int, req models.IncidentRequest) (models.IncidentResponse, error) {
	_, err := s.repo.GetIncidentById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println(err)
			return models.IncidentResponse{}, err
		}
		log.Println(err)
		return models.IncidentResponse{}, err
	}

	updateIncident, err := s.repo.UpdateIncident(id, req)
	if err != nil {
		log.Println(err)
		return models.IncidentResponse{}, err
	}

	return updateIncident, nil
}

func (s *IncidentService) DeleteIncident(id int) error {
	err := s.repo.DeleteIncident(id)
	if err != nil {
		return err
	}
	
	return nil
}

func (s *IncidentService) GetIncidentStats(ctx context.Context) (models.IncidentStatsResponse, error) {

    window := time.Duration(s.windowMin) * time.Minute

    count, err := s.repo.GetDangerStats(ctx, window)
    if err != nil {
        return models.IncidentStatsResponse{}, err
    }

    return models.IncidentStatsResponse{
        UserCount:    count,
        WindowMinute: s.windowMin,
    }, nil
}