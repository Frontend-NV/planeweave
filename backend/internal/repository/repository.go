package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/planeweave/planeweave/internal/domain"
	"github.com/planeweave/planeweave/internal/repository/sqlc"
)

type Repository struct {
	pool *pgxpool.Pool
	q    *sqlc.Queries
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool, q: sqlc.New(pool)}
}

func (r *Repository) WithTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(r.q.WithTx(tx)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func IsUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "duplicate key")
}

func MapNotFound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	return err
}

func (r *Repository) CreateUser(ctx context.Context, email, passwordHash, displayName string) (domain.User, error) {
	row, err := r.q.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
		DisplayName:  displayName,
	})
	if err != nil {
		return domain.User{}, err
	}
	return domain.UserFromCreateRow(row), nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (sqlc.GetUserByEmailRow, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	return row, MapNotFound(err)
}

func (r *Repository) GetUserByID(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	row, err := r.q.GetUserByID(ctx, userID)
	if err != nil {
		return domain.User{}, MapNotFound(err)
	}
	return domain.UserFromGetRow(row), nil
}

func (r *Repository) ListProjects(ctx context.Context, userID uuid.UUID) ([]domain.Project, error) {
	rows, err := r.q.ListProjectsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.Project, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.ProjectFromSQL(row))
	}
	return out, nil
}

func (r *Repository) CreateProject(ctx context.Context, title string, ownerID uuid.UUID) (domain.Project, error) {
	var project domain.Project
	err := r.WithTx(ctx, func(q *sqlc.Queries) error {
		p, err := q.CreateProject(ctx, sqlc.CreateProjectParams{Title: title, OwnerID: ownerID})
		if err != nil {
			return err
		}
		if err := q.AddProjectMember(ctx, sqlc.AddProjectMemberParams{
			ProjectID: p.ID,
			UserID:    ownerID,
			Role:      "owner",
		}); err != nil {
			return err
		}
		project = domain.ProjectFromSQL(p)
		return nil
	})
	return project, err
}

func (r *Repository) GetProject(ctx context.Context, projectID uuid.UUID) (domain.Project, error) {
	row, err := r.q.GetProject(ctx, projectID)
	if err != nil {
		return domain.Project{}, MapNotFound(err)
	}
	return domain.ProjectFromSQL(row), nil
}

func (r *Repository) UpdateProjectTitle(ctx context.Context, projectID uuid.UUID, title string) (domain.Project, error) {
	row, err := r.q.UpdateProjectTitle(ctx, sqlc.UpdateProjectTitleParams{Title: title, ID: projectID})
	if err != nil {
		return domain.Project{}, MapNotFound(err)
	}
	return domain.ProjectFromSQL(row), nil
}

func (r *Repository) DeleteProject(ctx context.Context, projectID uuid.UUID) error {
	return r.q.DeleteProject(ctx, projectID)
}

func (r *Repository) GetProjectOwnerID(ctx context.Context, projectID uuid.UUID) (uuid.UUID, error) {
	ownerID, err := r.q.GetProjectOwnerID(ctx, projectID)
	return ownerID, MapNotFound(err)
}

func (r *Repository) IsProjectMember(ctx context.Context, projectID, userID uuid.UUID) (bool, error) {
	return r.q.IsProjectMember(ctx, sqlc.IsProjectMemberParams{
		ProjectID: projectID,
		UserID:    userID,
	})
}

func (r *Repository) ListProjectMembers(ctx context.Context, projectID uuid.UUID) ([]domain.ProjectMember, error) {
	rows, err := r.q.ListProjectMembers(ctx, projectID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.ProjectMember, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.MemberFromSQL(row))
	}
	return out, nil
}

func (r *Repository) CreateInvite(ctx context.Context, token string, projectID, createdBy uuid.UUID, expiresAt time.Time) error {
	return r.q.CreateInvite(ctx, sqlc.CreateInviteParams{
		Token:     token,
		ProjectID: projectID,
		CreatedBy: createdBy,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
}

func (r *Repository) GetInvite(ctx context.Context, token string) (sqlc.GetInviteRow, error) {
	row, err := r.q.GetInvite(ctx, token)
	return row, MapNotFound(err)
}

func (r *Repository) AddProjectMember(ctx context.Context, projectID, userID uuid.UUID, role string) error {
	return r.q.AddProjectMember(ctx, sqlc.AddProjectMemberParams{
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
	})
}

func (r *Repository) ListTasks(ctx context.Context, projectID uuid.UUID) ([]domain.Task, error) {
	rows, err := r.q.ListTasksByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.Task, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.TaskFromListRow(row))
	}
	return out, nil
}

func (r *Repository) CreateTask(ctx context.Context, projectID uuid.UUID, title, status string, assigneeID *uuid.UUID, x, y float64) (domain.Task, error) {
	row, err := r.q.CreateTask(ctx, sqlc.CreateTaskParams{
		ProjectID:  projectID,
		Title:      title,
		Status:     status,
		AssigneeID: assigneeID,
		X:          x,
		Y:          y,
	})
	if err != nil {
		return domain.Task{}, err
	}
	return domain.TaskFromCreateRow(row), nil
}

func (r *Repository) GetTask(ctx context.Context, taskID uuid.UUID) (domain.Task, error) {
	row, err := r.q.GetTask(ctx, taskID)
	if err != nil {
		return domain.Task{}, MapNotFound(err)
	}
	return domain.TaskFromGetRow(row), nil
}

func (r *Repository) UpdateTask(ctx context.Context, taskID uuid.UUID, title, status string, assigneeID *uuid.UUID, x, y float64) (domain.Task, error) {
	row, err := r.q.UpdateTask(ctx, sqlc.UpdateTaskParams{
		ID:         taskID,
		Title:      title,
		Status:     status,
		AssigneeID: assigneeID,
		X:          x,
		Y:          y,
	})
	if err != nil {
		return domain.Task{}, MapNotFound(err)
	}
	return domain.TaskFromUpdateRow(row), nil
}

func (r *Repository) DeleteTask(ctx context.Context, taskID uuid.UUID) error {
	return r.q.DeleteTask(ctx, taskID)
}

func (r *Repository) GetTaskProjectID(ctx context.Context, taskID uuid.UUID) (uuid.UUID, error) {
	projectID, err := r.q.GetTaskProjectID(ctx, taskID)
	return projectID, MapNotFound(err)
}

func (r *Repository) ListDependencies(ctx context.Context, projectID uuid.UUID) ([]domain.Dependency, error) {
	rows, err := r.q.ListDependenciesByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.Dependency, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.DependencyFromSQL(row))
	}
	return out, nil
}

func (r *Repository) ListDependencyEdges(ctx context.Context, projectID uuid.UUID) ([]domain.GraphEdge, error) {
	rows, err := r.q.ListDependencyEdgesByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.GraphEdge, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.EdgeFromSQL(row))
	}
	return out, nil
}

func (r *Repository) CreateDependency(ctx context.Context, projectID, fromTaskID, toTaskID uuid.UUID) (domain.Dependency, error) {
	row, err := r.q.CreateDependency(ctx, sqlc.CreateDependencyParams{
		ProjectID:  projectID,
		FromTaskID: fromTaskID,
		ToTaskID:   toTaskID,
	})
	if err != nil {
		return domain.Dependency{}, err
	}
	return domain.DependencyFromSQL(row), nil
}

func (r *Repository) GetDependencyProjectID(ctx context.Context, dependencyID uuid.UUID) (uuid.UUID, error) {
	projectID, err := r.q.GetDependencyProjectID(ctx, dependencyID)
	return projectID, MapNotFound(err)
}

func (r *Repository) DeleteDependency(ctx context.Context, dependencyID uuid.UUID) error {
	return r.q.DeleteDependency(ctx, dependencyID)
}
