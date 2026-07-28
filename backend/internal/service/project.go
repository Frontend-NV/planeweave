package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/planeweave/planeweave/internal/domain"
	"github.com/planeweave/planeweave/internal/repository"
)

type ProjectService struct {
	repo   *repository.Repository
	access *AccessService
}

func NewProjectService(repo *repository.Repository, access *AccessService) *ProjectService {
	return &ProjectService{repo: repo, access: access}
}

func (s *ProjectService) List(ctx context.Context, userID uuid.UUID) ([]domain.Project, error) {
	return s.repo.ListProjects(ctx, userID)
}

func (s *ProjectService) Create(ctx context.Context, userID uuid.UUID, title string) (domain.Project, error) {
	if title == "" {
		return domain.Project{}, domain.ErrInvalidRequest
	}
	return s.repo.CreateProject(ctx, title, userID)
}

func (s *ProjectService) Get(ctx context.Context, projectID, userID uuid.UUID) (domain.Project, error) {
	if err := s.access.EnsureMember(ctx, projectID, userID); err != nil {
		return domain.Project{}, err
	}
	return s.repo.GetProject(ctx, projectID)
}

func (s *ProjectService) UpdateTitle(ctx context.Context, projectID, userID uuid.UUID, title string) (domain.Project, error) {
	if title == "" {
		return domain.Project{}, domain.ErrInvalidRequest
	}
	if err := s.access.EnsureMember(ctx, projectID, userID); err != nil {
		return domain.Project{}, err
	}
	return s.repo.UpdateProjectTitle(ctx, projectID, title)
}

func (s *ProjectService) Delete(ctx context.Context, projectID, userID uuid.UUID) error {
	ownerID, err := s.repo.GetProjectOwnerID(ctx, projectID)
	if err != nil {
		return err
	}
	if ownerID != userID {
		return domain.ErrForbidden
	}
	return s.repo.DeleteProject(ctx, projectID)
}

func (s *ProjectService) Members(ctx context.Context, projectID, userID uuid.UUID) ([]domain.ProjectMember, error) {
	if err := s.access.EnsureMember(ctx, projectID, userID); err != nil {
		return nil, err
	}
	return s.repo.ListProjectMembers(ctx, projectID)
}

func (s *ProjectService) EnsureMember(ctx context.Context, projectID, userID uuid.UUID) error {
	return s.access.EnsureMember(ctx, projectID, userID)
}
