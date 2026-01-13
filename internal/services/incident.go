package services


import (
	"errors"
	"database/sql"
	"github.com/rusinadaria/geo-notification-system/internal/repository"
	"github.com/rusinadaria/geo-notification-system/internal/models"
)

type IncidentService struct {
	repo repository.Incident
}

func NewIncidentService(repo repository.Incident) *IncidentService {
	return &IncidentService{repo: repo}
}

var ErrIncidentAlreadyExists = errors.New("incident already exists")

func (s *IncidentService) CheckLocation(checkReq models.LocationCheckRequest) (models.LocationCheckResponse, error) {
	// Алгоритм нахождения ближайших зон

	nearbyResp, err := s.repo.CheckLocation(checkReq)
	if err != nil {
		return models.LocationCheckResponse{}, err
	}

	err = s.repo.SaveCheck(checkReq.UserID, checkReq.Lat, checkReq.Lon, true) // ФУНКЦИЯ СОХРАНЕНИЯ ФАКТА ПРОВЕРКИ
	if err != nil {
        // log.Errorf("failed to save location check: %v", err)
		return models.LocationCheckResponse{}, err
        // TODO: можно поставить в Redis queue для retry
    }

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

func (s *IncidentService) GetIncidentById(id int) (models.Incident, error) {
	incident, err := s.repo.GetIncidentById(id)
	if err != nil {
		return models.Incident{}, err
	}

	return incident, nil
}

func (s *IncidentService) UpdateIncident(id int, req models.IncidentRequest) (models.Incident, error) {
	_, err := s.repo.GetIncidentById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Incident{}, err
		}
		return models.Incident{}, err
	}

	updateIncident, err := s.repo.UpdateIncident(id, req)
	if err != nil {
		return models.Incident{}, err
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