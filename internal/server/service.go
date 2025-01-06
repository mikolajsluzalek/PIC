package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"api/internal/api"
	"api/internal/models"
	"api/internal/storage"
	"github.com/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/jwtauth"
)

type Service struct {
	Config Config

	TokenAuth *jwtauth.JWTAuth
	Storage   storage.Service
	API       api.ServiceInterface
}

func NewService() (*Service, error) {
	svc := &Service{}

	cfg, err := readConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config")
	}

	svc.Config = cfg

	svc.Storage, err = storage.New()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create storage service")
	}

	svc.API, err = api.New()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create API service")
	}

	svc.TokenAuth = jwtauth.New("HS256", cfg.JWTSecret, nil)

	fmt.Println("Server initialized successfully!")

	return svc, nil
}
func (s *Service) Handler() http.Handler {
	logger := httplog.NewLogger("api", httplog.Options{
		LogLevel:         slog.LevelDebug,
		JSON:             true,
		Concise:          true,
		RequestHeaders:   true,
		ResponseHeaders:  true,
		MessageFieldName: "message",
		LevelFieldName:   "severity",
		TimeFieldFormat:  time.RFC3339,
		Tags: map[string]string{
			"env": "prod",
		},
		QuietDownRoutes: []string{
			"/",
			"/ping",
		},
		QuietDownPeriod: 10 * time.Second,
		// SourceFieldName: "source",
	})

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		//AllowedOrigins: []string{}, // Use this to allow specific origin hosts
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	router.Use(httplog.RequestLogger(logger, []string{"/ping"}))
	router.Use(middleware.Heartbeat("/ping"))

	// Router for routes requiring authorization
	router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(s.TokenAuth))

		r.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
			resp, err := s.API.Dashboard(r.Context())
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(resp)
		})

		r.Get("/cars", func(w http.ResponseWriter, r *http.Request) {
			cars, err := s.API.Cars(r.Context())
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(cars)
		})

		r.Get("/car/{id}", func(w http.ResponseWriter, r *http.Request) {
			stringId := chi.URLParam(r, "id")

			if stringId == "" {
				http.Error(w, "id is required", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(stringId)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			car, err := s.API.GetCar(r.Context(), id)
			if err != nil {
				// TODO 404
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(car)
		})

		r.Get("/car/numbers", func(w http.ResponseWriter, r *http.Request) {
			projects, err := s.API.GetCarNumbers(r.Context())
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(projects)
		})

		r.Delete("/car/{id}", func(w http.ResponseWriter, r *http.Request) {
			stringId := chi.URLParam(r, "id")

			if stringId == "" {
				http.Error(w, "id is required", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(stringId)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			err = s.API.RemoveCar(r.Context(), id)
			if err != nil {
				// TODO 404
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		})

		r.Post("/car", func(w http.ResponseWriter, r *http.Request) {
			var newCar models.NewCar

			err := json.NewDecoder(r.Body).Decode(&newCar)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			car, err := s.API.AddCar(r.Context(), newCar)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(car)
		})

		r.Post("/car/{id}/update", func(w http.ResponseWriter, r *http.Request) {
			var updateCar models.UpdateCar

			err := json.NewDecoder(r.Body).Decode(&updateCar)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			stringId := chi.URLParam(r, "id")

			if stringId == "" {
				http.Error(w, "id is required", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(stringId)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			car, err := s.API.UpdateCar(r.Context(), id, updateCar)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(car)
		})

		r.Get("/projects", func(w http.ResponseWriter, r *http.Request) {
			projects, err := s.API.Projects(r.Context())
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(projects)
		})

		r.Get("/project/names", func(w http.ResponseWriter, r *http.Request) {
			projects, err := s.API.GetProjectNames(r.Context())
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(projects)
		})

		r.Get("/project/{id}", func(w http.ResponseWriter, r *http.Request) {
			stringId := chi.URLParam(r, "id")

			if stringId == "" {
				http.Error(w, "id is required", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(stringId)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			project, err := s.API.GetProject(r.Context(), id)
			if err != nil {
				// TODO 404
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(project)
		})

		r.Post("/project", func(w http.ResponseWriter, r *http.Request) {
			var newProject models.NewProject

			err := json.NewDecoder(r.Body).Decode(&newProject)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			project, err := s.API.AddProject(r.Context(), newProject)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(project)
		})

		r.Post("/project/{id}/update", func(w http.ResponseWriter, r *http.Request) {
			var updateProject models.UpdateProject

			err := json.NewDecoder(r.Body).Decode(&updateProject)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			stringId := chi.URLParam(r, "id")

			if stringId == "" {
				http.Error(w, "id is required", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(stringId)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			car, err := s.API.UpdateProject(r.Context(), id, updateProject)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(car)
		})

		r.Delete("/project/{id}", func(w http.ResponseWriter, r *http.Request) {
			stringId := chi.URLParam(r, "id")

			if stringId == "" {
				http.Error(w, "id is required", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(stringId)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			err = s.API.RemoveProject(r.Context(), id)
			if err != nil {
				// TODO 404
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		})

		r.Get("/accommodations", func(w http.ResponseWriter, r *http.Request) {
			accommodations, err := s.API.Accommodations(r.Context())
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(accommodations)
		})

		r.Get("/accommodation/{id}", func(w http.ResponseWriter, r *http.Request) {
			stringId := chi.URLParam(r, "id")

			if stringId == "" {
				http.Error(w, "id is required", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(stringId)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			accommodation, err := s.API.GetAccommodation(r.Context(), id)
			if err != nil {
				// TODO 404
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(accommodation)
		})

		r.Get("/accommodation/addresses", func(w http.ResponseWriter, r *http.Request) {
			addresses, err := s.API.GetAccommodationAddresses(r.Context())
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(addresses)
		})

		r.Post("/accommodation", func(w http.ResponseWriter, r *http.Request) {
			var newAcc models.NewAccommodation

			err := json.NewDecoder(r.Body).Decode(&newAcc)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			acc, err := s.API.AddAccommodation(r.Context(), newAcc)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(acc)
		})

		r.Post("/accommodation/{id}/update", func(w http.ResponseWriter, r *http.Request) {
			var updateAccommodation models.UpdateAccommodation

			err := json.NewDecoder(r.Body).Decode(&updateAccommodation)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			stringId := chi.URLParam(r, "id")

			if stringId == "" {
				http.Error(w, "id is required", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(stringId)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			car, err := s.API.UpdateAccommodation(r.Context(), id, updateAccommodation)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(car)
		})

		r.Delete("/accommodation/{id}", func(w http.ResponseWriter, r *http.Request) {
			stringId := chi.URLParam(r, "id")

			if stringId == "" {
				http.Error(w, "id is required", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(stringId)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			err = s.API.RemoveAccommodation(r.Context(), id)
			if err != nil {
				// TODO 404
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		})

		r.Get("/employees", func(w http.ResponseWriter, r *http.Request) {
			employees, err := s.API.Employees(r.Context())
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(employees)
		})

		r.Get("/employee/{id}", func(w http.ResponseWriter, r *http.Request) {
			stringId := chi.URLParam(r, "id")

			if stringId == "" {
				http.Error(w, "id is required", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(stringId)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			employee, err := s.API.GetEmployee(r.Context(), id)
			if err != nil {
				// TODO 404
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(employee)
		})

		r.Post("/employee", func(w http.ResponseWriter, r *http.Request) {
			var newEmployee models.NewEmployee

			err := json.NewDecoder(r.Body).Decode(&newEmployee)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			employee, err := s.API.AddEmployee(r.Context(), newEmployee)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(employee)
		})

		r.Post("/employee/{id}/update", func(w http.ResponseWriter, r *http.Request) {
			var updateEmployee models.UpdateEmployee

			err := json.NewDecoder(r.Body).Decode(&updateEmployee)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			stringId := chi.URLParam(r, "id")

			if stringId == "" {
				http.Error(w, "id is required", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(stringId)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			employee, err := s.API.UpdateEmployee(r.Context(), id, updateEmployee)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(employee)
		})

		r.Delete("/employee/{id}", func(w http.ResponseWriter, r *http.Request) {
			stringId := chi.URLParam(r, "id")

			if stringId == "" {
				http.Error(w, "id is required", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(stringId)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			err = s.API.RemoveEmployee(r.Context(), id)
			if err != nil {
				// TODO 404
				logger.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		})

	})

	router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		var req models.LoginRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		resp, err := s.API.Login(r.Context(), req.Username, req.Password)
		if err != nil {
			if errors.Is(err, api.ErrUnauthorized) {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			logger.Error(err.Error())
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	})

	return router
}
