package api

import (
	"context"

	"api/internal/models"
	"github.com/pkg/errors"
)

func (s *Service) Projects(ctx context.Context) ([]models.Project, error) {
	projects, err := s.storage.Projects(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve projects")
	}

	return projects, nil
}

func (s *Service) GetProject(ctx context.Context, id int) (models.Project, error) {
	project, err := s.storage.GetProject(ctx, id)

	return project, errors.Wrap(err, "failed to retrieve project")
}

func (s *Service) AddProject(ctx context.Context, newProject models.NewProject) (models.Project, error) {
	id, err := s.storage.AddProject(ctx, newProject)

	if err != nil {
		return models.Project{}, errors.Wrap(err, "failed to add project")
	}

	project, err := s.GetProject(ctx, id)

	return project, errors.Wrap(err, "failed to add project")
}

func (s *Service) RemoveProject(ctx context.Context, id int) error {
	err := s.storage.RemoveProject(ctx, id)

	return errors.Wrap(err, "failed to remove project")
}

func (s *Service) UpdateProject(ctx context.Context, id int, updateProject models.UpdateProject) (models.Project, error) {
	err := s.storage.UpdateProject(ctx, id, updateProject)
	if err != nil {
		return models.Project{}, errors.Wrap(err, "failed to update project")
	}

	project, err := s.GetProject(ctx, id)

	return project, errors.Wrap(err, "failed to retrieve updated project")
}

func (s *Service) GetProjectNames(ctx context.Context) ([]models.ProjectNames, error) {
	projectNames, err := s.storage.GetProjectNames(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve projects names")
	}

	return projectNames, nil
}
