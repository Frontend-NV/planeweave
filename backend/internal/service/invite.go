package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"

	"github.com/planeweave/planeweave/internal/domain"
	"github.com/planeweave/planeweave/internal/repository"
)

type InviteService struct {
	repo   *repository.Repository
	access *AccessService
}

func NewInviteService(repo *repository.Repository, access *AccessService) *InviteService {
	return &InviteService{repo: repo, access: access}
}

func (s *InviteService) Create(ctx context.Context, projectID, userID uuid.UUID) (domain.Invite, error) {
	if err := s.access.EnsureMember(ctx, projectID, userID); err != nil {
		return domain.Invite{}, err
	}

	token, err := randomToken(32)
	if err != nil {
		return domain.Invite{}, err
	}
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err := s.repo.CreateInvite(ctx, token, projectID, userID, expiresAt); err != nil {
		return domain.Invite{}, err
	}
	return domain.Invite{Token: token, ExpiresAt: expiresAt}, nil
}

func (s *InviteService) Accept(ctx context.Context, token string, userID uuid.UUID) (domain.Project, error) {
	row, err := s.repo.GetInvite(ctx, token)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.Project{}, domain.ErrInviteNotFound
		}
		return domain.Project{}, err
	}
	if !row.ExpiresAt.Valid || time.Now().After(row.ExpiresAt.Time) {
		return domain.Project{}, domain.ErrInviteExpired
	}

	if err := s.repo.AddProjectMember(ctx, row.ProjectID, userID, "member"); err != nil {
		return domain.Project{}, err
	}
	return s.repo.GetProject(ctx, row.ProjectID)
}

func randomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
