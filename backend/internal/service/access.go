package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/planeweave/planeweave/internal/domain"
	"github.com/planeweave/planeweave/internal/repository"
)

type AccessService struct {
	repo *repository.Repository
}

func NewAccessService(repo *repository.Repository) *AccessService {
	return &AccessService{repo: repo}
}

func (s *AccessService) EnsureMember(ctx context.Context, projectID, userID uuid.UUID) error {
	ok, err := s.repo.IsProjectMember(ctx, projectID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return domain.ErrForbidden
	}
	return nil
}

func ParseUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}
