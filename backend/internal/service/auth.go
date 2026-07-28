package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"

	"github.com/planeweave/planeweave/internal/domain"
	"github.com/planeweave/planeweave/internal/pkg/jwt"
	"github.com/planeweave/planeweave/internal/repository"
)

type AuthService struct {
	repo   *repository.Repository
	tokens *jwt.Service
}

func NewAuthService(repo *repository.Repository, tokens *jwt.Service) *AuthService {
	return &AuthService{repo: repo, tokens: tokens}
}

func (s *AuthService) Register(ctx context.Context, email, password, displayName string) (domain.AuthResponse, error) {
	if email == "" || password == "" || displayName == "" {
		return domain.AuthResponse{}, domain.ErrInvalidRequest
	}
	if len(password) < 8 {
		return domain.AuthResponse{}, domain.ErrPasswordTooShort
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.AuthResponse{}, err
	}

	user, err := s.repo.CreateUser(ctx, email, string(hash), displayName)
	if err != nil {
		if repository.IsUniqueViolation(err) {
			return domain.AuthResponse{}, domain.ErrEmailTaken
		}
		return domain.AuthResponse{}, err
	}

	token, err := s.tokens.Generate(user.ID, user.Email)
	if err != nil {
		return domain.AuthResponse{}, err
	}
	return domain.AuthResponse{Token: token, User: user}, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (domain.AuthResponse, error) {
	if email == "" || password == "" {
		return domain.AuthResponse{}, domain.ErrInvalidRequest
	}

	row, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.AuthResponse{}, domain.ErrInvalidCredentials
		}
		return domain.AuthResponse{}, err
	}
	if bcrypt.CompareHashAndPassword([]byte(row.PasswordHash), []byte(password)) != nil {
		return domain.AuthResponse{}, domain.ErrInvalidCredentials
	}

	user := domain.User{
		ID:          domain.UUIDToString(row.ID),
		Email:       row.Email,
		DisplayName: row.DisplayName,
	}
	token, err := s.tokens.Generate(user.ID, user.Email)
	if err != nil {
		return domain.AuthResponse{}, err
	}
	return domain.AuthResponse{Token: token, User: user}, nil
}

func (s *AuthService) Me(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *AuthService) Tokens() *jwt.Service {
	return s.tokens
}
