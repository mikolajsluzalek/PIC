package api

import (
	"context"

	"api/internal/models"
	"github.com/pkg/errors"
)

func (s *Service) Cars(ctx context.Context) ([]models.Car, error) {
	cars, err := s.storage.Cars(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve cars")
	}

	return cars, nil
}

func (s *Service) GetCar(ctx context.Context, id int) (models.Car, error) {
	car, err := s.storage.GetCar(ctx, id)

	return car, errors.Wrap(err, "failed to retrieve car")
}

func (s *Service) AddCar(ctx context.Context, newCar models.NewCar) (models.Car, error) {
	id, err := s.storage.AddCar(ctx, newCar)

	if err != nil {
		return models.Car{}, errors.Wrap(err, "failed to add car")
	}

	car, err := s.GetCar(ctx, id)

	return car, errors.Wrap(err, "failed to add car")
}

func (s *Service) RemoveCar(ctx context.Context, id int) error {
	err := s.storage.RemoveCar(ctx, id)

	return errors.Wrap(err, "failed to remove car")
}

func (s *Service) UpdateCar(ctx context.Context, id int, updateCar models.UpdateCar) (models.Car, error) {
	err := s.storage.UpdateCar(ctx, id, updateCar)
	if err != nil {
		return models.Car{}, errors.Wrap(err, "failed to update car")
	}

	car, err := s.GetCar(ctx, id)

	return car, errors.Wrap(err, "failed to retrieve updated car")
}

func (s *Service) GetCarNumbers(ctx context.Context) ([]models.CarNumbers, error) {
	carNumbers, err := s.storage.GetCarNumbers(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve car numbers")
	}

	return carNumbers, nil
}
