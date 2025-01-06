package api

import (
	"api/internal/models"
	"context"
	"github.com/pkg/errors"
)

func (s *Service) Accommodations(ctx context.Context) ([]models.Accommodation, error) {
	accommodations, err := s.storage.Accommodations(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve accommodations")
	}

	return accommodations, nil
}

func (s *Service) GetAccommodation(ctx context.Context, id int) (models.Accommodation, error) {
	accommodation, err := s.storage.GetAccommodation(ctx, id)

	return accommodation, errors.Wrap(err, "failed to retrieve accommodation")
}

func (s *Service) AddAccommodation(ctx context.Context, newAccommodation models.NewAccommodation) (models.Accommodation, error) {
	id, err := s.storage.AddAccommodation(ctx, newAccommodation)

	if err != nil {
		errors.Wrap(err, "failed to add accommodations")
	}

	accommodation, err := s.GetAccommodation(ctx, id)

	return accommodation, errors.Wrap(err, "failed to add accommodation")
}

func (s *Service) RemoveAccommodation(ctx context.Context, id int) error {
	err := s.storage.RemoveAccommodation(ctx, id)

	return errors.Wrap(err, "failed to remove accommodation")
}

func (s *Service) UpdateAccommodation(ctx context.Context, id int, updateAccommodation models.UpdateAccommodation) (models.Accommodation, error) {
	err := s.storage.UpdateAccommodation(ctx, id, updateAccommodation)
	if err != nil {
		return models.Accommodation{}, errors.Wrap(err, "failed to update accommodation")
	}

	accommodation, err := s.GetAccommodation(ctx, id)

	return accommodation, errors.Wrap(err, "failed to retrieve updated accommodation")
}

func (s *Service) GetAccommodationAddresses(ctx context.Context) ([]models.AccommodationAddresses, error) {
	accommodationAddresses, err := s.storage.GetAccommodationAddresses(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve accommodation addresses")
	}

	return accommodationAddresses, nil
}
