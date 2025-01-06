package api

import (
	"api/internal/models"
	"context"
	"github.com/pkg/errors"
)

func (s *Service) Employees(ctx context.Context) ([]models.Employee, error) {
	employees, err := s.storage.Employees(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve employees")
	}

	return employees, nil
}

func (s *Service) GetEmployee(ctx context.Context, id int) (models.Employee, error) {
	employee, err := s.storage.GetEmployee(ctx, id)

	return employee, errors.Wrap(err, "failed to retrieve employee")
}

func (s *Service) AddEmployee(ctx context.Context, newEmployee models.NewEmployee) (models.Employee, error) {
	id, err := s.storage.AddEmployee(ctx, newEmployee)

	if err != nil {
		return models.Employee{}, errors.Wrap(err, "failed to add employee")
	}

	employee, err := s.GetEmployee(ctx, id)

	return employee, errors.Wrap(err, "failed to add employee")
}

func (s *Service) UpdateEmployee(ctx context.Context, id int, updateEmployee models.UpdateEmployee) (models.Employee, error) {
	err := s.storage.UpdateEmployee(ctx, id, updateEmployee)
	if err != nil {
		return models.Employee{}, errors.Wrap(err, "failed to update employee")
	}

	employee, err := s.GetEmployee(ctx, id)

	return employee, errors.Wrap(err, "failed to retrieve updated employee")
}

func (s *Service) RemoveEmployee(ctx context.Context, id int) error {
	err := s.storage.RemoveEmployee(ctx, id)

	return errors.Wrap(err, "failed to remove accommodation")
}
