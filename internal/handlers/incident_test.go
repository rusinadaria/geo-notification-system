package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/rusinadaria/geo-notification-system/internal/models"
	"github.com/rusinadaria/geo-notification-system/internal/services"
	mock_service "github.com/rusinadaria/geo-notification-system/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestHandler_CreateIncident(t *testing.T) {
	type mockBehavior func(s *mock_service.MockIncident, req models.IncidentRequest)

	testTable := []struct {
		name               string
		inputBody          string
		inputReq           models.IncidentRequest
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name: "OK",
			inputBody: `{
				"type": "danger",
				"description": "desc_inc_4",
				"latitude": 55.751244,
				"longitude": 37.618423,
				"radius_meters": 100,
				"active": true
			}`,
			inputReq: models.IncidentRequest{
				Type:         "danger",
				Description:  "desc_inc_4",
				Latitude:     55.751244,
				Longitude:    37.618423,
				RadiusMeters: 100,
				Active:       true,
			},
			mockBehavior: func(s *mock_service.MockIncident, req models.IncidentRequest) {
				s.EXPECT().CreateIncident(req)
			},
			expectedStatusCode: 200,
		},
		{
			name:               "Invalid JSON",
			inputBody:          `{ "type": "danger", `,
			mockBehavior:       func(s *mock_service.MockIncident, req models.IncidentRequest) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Empty body",
			inputBody:          ``,
			mockBehavior:       func(s *mock_service.MockIncident, req models.IncidentRequest) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Service error",
			inputBody: `{
				"type": "danger",
				"description": "desc",
				"latitude": 55.751244,
				"longitude": 37.618423,
				"radius_meters": 100,
				"active": true
			}`,
			inputReq: models.IncidentRequest{
				Type:         "danger",
				Description:  "desc",
				Latitude:     55.751244,
				Longitude:    37.618423,
				RadiusMeters: 100,
				Active:       true,
			},
			mockBehavior: func(s *mock_service.MockIncident, req models.IncidentRequest) {
				s.EXPECT().CreateIncident(req).
					Return(errors.New("service error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			incident := mock_service.NewMockIncident(c)
			testCase.mockBehavior(incident, testCase.inputReq)

			services := &services.Service{Incident: incident}
			handler := NewHandler(services)

			r := chi.NewRouter()
			r.Post("/api/v1/incidents/", handler.CreateIncidentHandler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/incidents/", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}

func TestHandler_GetIncident(t *testing.T) {
	type mockBehavior func(s *mock_service.MockIncident, id int)

	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	testTable := []struct {
		name                 string
		id                   int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			id:   2,
			mockBehavior: func(s *mock_service.MockIncident, id int) {
				s.EXPECT().
					GetIncidentById(id).
					Return(models.IncidentResponse{
						Type:         "danger",
						Description:  "desc_inc_4",
						Latitude:     55.751244,
						Longitude:    37.618423,
						RadiusMeters: 100,
						Active:       true,
						CreatedAt:    fixedTime,
						UpdatedAt:    fixedTime,
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: func() string {
				body, _ := json.Marshal(models.IncidentResponse{
					Type:         "danger",
					Description:  "desc_inc_4",
					Latitude:     55.751244,
					Longitude:    37.618423,
					RadiusMeters: 100,
					Active:       true,
					CreatedAt:    fixedTime,
					UpdatedAt:    fixedTime,
				})
				return string(body)
			}(),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			incident := mock_service.NewMockIncident(ctrl)
			testCase.mockBehavior(incident, testCase.id)

			services := &services.Service{Incident: incident}
			handler := NewHandler(services)

			r := chi.NewRouter()
			r.Get("/api/v1/incidents/{id}", handler.GetIncident)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/api/v1/incidents/%d", testCase.id),
				nil,
			)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_ListIncidents(t *testing.T) {
	type mockBehavior func(
		s *mock_service.MockIncident,
		limit int,
		offset int,
	)

	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	testTable := []struct {
		name                 string
		query                string
		limit                int
		offset               int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:   "OK",
			query:  "?limit=5&offset=10",
			limit:  5,
			offset: 10,
			mockBehavior: func(s *mock_service.MockIncident, limit, offset int) {
				s.EXPECT().
					GetAllIncidents(limit, offset).
					Return([]models.IncidentResponse{
						{
							Type:         "danger",
							Description:  "desc1",
							Latitude:     1,
							Longitude:    2,
							RadiusMeters: 100,
							Active:       true,
							CreatedAt:    fixedTime,
							UpdatedAt:    fixedTime,
						},
						{
							Type:         "warning",
							Description:  "desc2",
							Latitude:     3,
							Longitude:    4,
							RadiusMeters: 200,
							Active:       false,
							CreatedAt:    fixedTime,
							UpdatedAt:    fixedTime,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: func() string {
				body, _ := json.Marshal([]models.IncidentResponse{
					{
						Type:         "danger",
						Description:  "desc1",
						Latitude:     1,
						Longitude:    2,
						RadiusMeters: 100,
						Active:       true,
						CreatedAt:    fixedTime,
						UpdatedAt:    fixedTime,
					},
					{
						Type:         "warning",
						Description:  "desc2",
						Latitude:     3,
						Longitude:    4,
						RadiusMeters: 200,
						Active:       false,
						CreatedAt:    fixedTime,
						UpdatedAt:    fixedTime,
					},
				})
				return string(body)
			}(),
		},
		{
			name:   "Default params",
			query:  "",
			limit:  10,
			offset: 0,
			mockBehavior: func(s *mock_service.MockIncident, limit, offset int) {
				s.EXPECT().
					GetAllIncidents(limit, offset).
					Return([]models.IncidentResponse{}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[]`,
		},
		{
			name:   "Invalid params",
			query:  "?limit=-1&offset=-5",
			limit:  10,
			offset: 0,
			mockBehavior: func(s *mock_service.MockIncident, limit, offset int) {
				s.EXPECT().
					GetAllIncidents(limit, offset).
					Return([]models.IncidentResponse{}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[]`,
		},
		{
			name:   "Service error",
			query:  "?limit=5&offset=0",
			limit:  5,
			offset: 0,
			mockBehavior: func(s *mock_service.MockIncident, limit, offset int) {
				s.EXPECT().
					GetAllIncidents(limit, offset).
					Return(nil, errors.New("db error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			incident := mock_service.NewMockIncident(ctrl)

			if testCase.mockBehavior != nil {
				testCase.mockBehavior(
					incident,
					testCase.limit,
					testCase.offset,
				)
			}

			services := &services.Service{Incident: incident}
			handler := NewHandler(services)

			r := chi.NewRouter()
			r.Get("/api/v1/incidents", handler.ListIncidents)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodGet,
				"/api/v1/incidents"+testCase.query,
				nil,
			)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)

			if testCase.expectedResponseBody != "" {
				assert.JSONEq(t, testCase.expectedResponseBody, w.Body.String())
			}
		})
	}
}

func TestHandler_UpdateIncident(t *testing.T) {
	type mockBehavior func(
		s *mock_service.MockIncident,
		id int,
		req models.IncidentRequest,
	)

	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	testTable := []struct {
		name                 string
		id                   string
		inputBody            string
		inputReq             models.IncidentRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			id:   "2",
			inputBody: `{
				"type": "danger",
				"description": "updated desc",
				"latitude": 55.75,
				"longitude": 37.61,
				"radius_meters": 150,
				"active": false
			}`,
			inputReq: models.IncidentRequest{
				Type:         "danger",
				Description:  "updated desc",
				Latitude:     55.75,
				Longitude:    37.61,
				RadiusMeters: 150,
				Active:       false,
			},
			mockBehavior: func(
				s *mock_service.MockIncident,
				id int,
				req models.IncidentRequest,
			) {
				s.EXPECT().
					UpdateIncident(id, req).
					Return(models.IncidentResponse{
						Type:         req.Type,
						Description:  req.Description,
						Latitude:     req.Latitude,
						Longitude:    req.Longitude,
						RadiusMeters: req.RadiusMeters,
						Active:       req.Active,
						CreatedAt:    fixedTime,
						UpdatedAt:    fixedTime,
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: func() string {
				body, _ := json.Marshal(models.IncidentResponse{
					Type:         "danger",
					Description:  "updated desc",
					Latitude:     55.75,
					Longitude:    37.61,
					RadiusMeters: 150,
					Active:       false,
					CreatedAt:    fixedTime,
					UpdatedAt:    fixedTime,
				})
				return string(body)
			}(),
		},
		{
			name:               "Invalid ID",
			id:                 "abc",
			inputBody:          `{}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Invalid Body",
			id:                 "1",
			inputBody:          `{`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid Radius",
			id:   "1",
			inputBody: `{
				"type": "danger",
				"radius_meters": 0
			}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Service Error",
			id:   "3",
			inputBody: `{
				"type": "danger",
				"description": "desc",
				"latitude": 1,
				"longitude": 1,
				"radius_meters": 10,
				"active": true
			}`,
			inputReq: models.IncidentRequest{
				Type:         "danger",
				Description:  "desc",
				Latitude:     1,
				Longitude:    1,
				RadiusMeters: 10,
				Active:       true,
			},
			mockBehavior: func(
				s *mock_service.MockIncident,
				id int,
				req models.IncidentRequest,
			) {
				s.EXPECT().
					UpdateIncident(id, req).
					Return(models.IncidentResponse{}, errors.New("update error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			incident := mock_service.NewMockIncident(ctrl)

			if testCase.mockBehavior != nil {
				idInt, _ := strconv.Atoi(testCase.id)
				testCase.mockBehavior(incident, idInt, testCase.inputReq)
			}

			services := &services.Service{Incident: incident}
			handler := NewHandler(services)

			r := chi.NewRouter()
			r.Put("/api/v1/incidents/{id}", handler.UpdateIncident)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodPut,
				"/api/v1/incidents/"+testCase.id,
				bytes.NewBufferString(testCase.inputBody),
			)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)

			if testCase.expectedResponseBody != "" {
				assert.JSONEq(t, testCase.expectedResponseBody, w.Body.String())
			}
		})
	}
}

func TestHandler_DeleteIncident(t *testing.T) {
	type mockBehavior func(s *mock_service.MockIncident, id int)

	testTable := []struct {
		name               string
		id                 string
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name: "OK",
			id:   "2",
			mockBehavior: func(s *mock_service.MockIncident, id int) {
				s.EXPECT().DeleteIncident(id).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Invalid ID",
			id:   "abc",
			mockBehavior: func(s *mock_service.MockIncident, id int) {

			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Service Error",
			id:   "3",
			mockBehavior: func(s *mock_service.MockIncident, id int) {
				s.EXPECT().DeleteIncident(id).Return(errors.New("delete error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			incident := mock_service.NewMockIncident(ctrl)

			if testCase.name != "Invalid ID" {
				idInt, _ := strconv.Atoi(testCase.id)
				testCase.mockBehavior(incident, idInt)
			}

			services := &services.Service{Incident: incident}
			handler := NewHandler(services)

			r := chi.NewRouter()
			r.Delete("/api/v1/incidents/{id}", handler.DeleteIncident)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodDelete,
				"/api/v1/incidents/"+testCase.id,
				nil,
			)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}
