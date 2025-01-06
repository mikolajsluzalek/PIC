package api

import (
	"context"

	"api/internal/models"
	"api/internal/storage"
	"github.com/pkg/errors"
)

type ServiceInterface interface {
	Login(ctx context.Context, username, password string) (models.LoginResponse, error)

	Dashboard(ctx context.Context) (models.Dashboard, error)

	Cars(ctx context.Context) ([]models.Car, error)
	GetCar(ctx context.Context, id int) (models.Car, error)
	AddCar(ctx context.Context, newCar models.NewCar) (models.Car, error)
	UpdateCar(ctx context.Context, id int, updateCar models.UpdateCar) (models.Car, error)
	RemoveCar(ctx context.Context, id int) error
	GetCarNumbers(ctx context.Context) ([]models.CarNumbers, error)

	Projects(ctx context.Context) ([]models.Project, error)
	GetProject(ctx context.Context, id int) (models.Project, error)
	AddProject(ctx context.Context, newProject models.NewProject) (models.Project, error)
	UpdateProject(ctx context.Context, id int, updateProject models.UpdateProject) (models.Project, error)
	RemoveProject(ctx context.Context, id int) error
	GetProjectNames(ctx context.Context) ([]models.ProjectNames, error)

	Accommodations(ctx context.Context) ([]models.Accommodation, error)
	GetAccommodation(ctx context.Context, id int) (models.Accommodation, error)
	AddAccommodation(ctx context.Context, newAccommodation models.NewAccommodation) (models.Accommodation, error)
	UpdateAccommodation(ctx context.Context, id int, updateAccommodation models.UpdateAccommodation) (models.Accommodation, error)
	RemoveAccommodation(ctx context.Context, id int) error
	GetAccommodationAddresses(ctx context.Context) ([]models.AccommodationAddresses, error)

	Employees(ctx context.Context) ([]models.Employee, error)
	GetEmployee(ctx context.Context, id int) (models.Employee, error)
	AddEmployee(ctx context.Context, newEmployee models.NewEmployee) (models.Employee, error)
	UpdateEmployee(ctx context.Context, id int, updateEmployee models.UpdateEmployee) (models.Employee, error)
	RemoveEmployee(ctx context.Context, id int) error
}

type Service struct {
	Config  Config
	storage storage.Service
}

func New() (*Service, error) {
	svc := &Service{}

	cfg, err := readConfig()
	if err != nil {
		return svc, errors.Wrap(err, "failed to read config")
	}

	svc.Config = cfg

	svc.storage, err = storage.New()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create storage service")
	}

	return svc, nil
}
