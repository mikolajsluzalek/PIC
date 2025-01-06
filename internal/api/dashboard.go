package api

import (
	"context"

	"api/internal/models"
	"github.com/pkg/errors"
)

func (s *Service) Dashboard(ctx context.Context) (models.Dashboard, error) {
	projects, err := s.storage.DashboardEmployeeProjects(ctx)
	if err != nil {
		return models.Dashboard{}, errors.Wrap(err, "failed to retrieve dashboard data")
	}

	accommodations, err := s.storage.Accommodation(ctx)
	if err != nil {
		return models.Dashboard{}, errors.Wrap(err, "failed to retrieve dashboard data")
	}

	carInspections, err := s.storage.CarInspections(ctx)
	if err != nil {
		return models.Dashboard{}, errors.Wrap(err, "failed to retrieve dashboard data")
	}

	employeePermits, err := s.storage.EmployeePermits(ctx)
	if err != nil {
		return models.Dashboard{}, errors.Wrap(err, "failed to retrieve dashboard data")
	}

	return models.Dashboard{
		EmployeesProject: projects,
		Accommodations:   accommodations,
		CarInspections:   carInspections,
		EmployeePermits:  employeePermits,
	}, nil
}
