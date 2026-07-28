package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/planeweave/planeweave/internal/domain"
	"github.com/planeweave/planeweave/internal/pkg/graph"
	"github.com/planeweave/planeweave/internal/repository"
	"github.com/planeweave/planeweave/internal/ws"
)

type DependencyService struct {
	repo   *repository.Repository
	access *AccessService
	events EventPublisher
}

func NewDependencyService(repo *repository.Repository, access *AccessService, events EventPublisher) *DependencyService {
	return &DependencyService{repo: repo, access: access, events: events}
}

func (s *DependencyService) List(ctx context.Context, projectID, userID uuid.UUID) ([]domain.Dependency, error) {
	if err := s.access.EnsureMember(ctx, projectID, userID); err != nil {
		return nil, err
	}
	return s.repo.ListDependencies(ctx, projectID)
}

func (s *DependencyService) Create(ctx context.Context, projectID, userID uuid.UUID, fromTaskID, toTaskID string) (domain.Dependency, error) {
	if fromTaskID == "" || toTaskID == "" {
		return domain.Dependency{}, domain.ErrInvalidRequest
	}
	if err := s.access.EnsureMember(ctx, projectID, userID); err != nil {
		return domain.Dependency{}, err
	}

	fromUUID, err := uuid.Parse(fromTaskID)
	if err != nil {
		return domain.Dependency{}, domain.ErrInvalidRequest
	}
	toUUID, err := uuid.Parse(toTaskID)
	if err != nil {
		return domain.Dependency{}, domain.ErrInvalidRequest
	}

	edges, err := s.repo.ListDependencyEdges(ctx, projectID)
	if err != nil {
		return domain.Dependency{}, err
	}
	if graph.WouldCreateCycle(edges, fromTaskID, toTaskID) {
		return domain.Dependency{}, domain.ErrCycleDetected
	}

	dep, err := s.repo.CreateDependency(ctx, projectID, fromUUID, toUUID)
	if err != nil {
		return domain.Dependency{}, err
	}
	s.events.Broadcast(projectID.String(), ws.Event{Type: "dependency.created", Payload: dep})
	return dep, nil
}

func (s *DependencyService) Delete(ctx context.Context, dependencyID, userID uuid.UUID) error {
	projectID, err := s.repo.GetDependencyProjectID(ctx, dependencyID)
	if err != nil {
		return err
	}
	if err := s.access.EnsureMember(ctx, projectID, userID); err != nil {
		return err
	}
	if err := s.repo.DeleteDependency(ctx, dependencyID); err != nil {
		return err
	}
	s.events.Broadcast(projectID.String(), ws.Event{
		Type: "dependency.deleted",
		Payload: map[string]string{
			"id":        dependencyID.String(),
			"projectId": projectID.String(),
		},
	})
	return nil
}
