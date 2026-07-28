package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/planeweave/planeweave/internal/repository/sqlc"
)

func UUIDToString(id uuid.UUID) string {
	return id.String()
}

func StringToUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}

func OptionalUUIDString(id *uuid.UUID) *string {
	if id == nil {
		return nil
	}
	s := id.String()
	return &s
}

func OptionalStringUUID(id *string) (*uuid.UUID, error) {
	if id == nil {
		return nil, nil
	}
	parsed, err := uuid.Parse(*id)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func Timestamptz(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

func UserFromCreateRow(row sqlc.CreateUserRow) User {
	return User{
		ID:          UUIDToString(row.ID),
		Email:       row.Email,
		DisplayName: row.DisplayName,
	}
}

func UserFromGetRow(row sqlc.GetUserByIDRow) User {
	return User{
		ID:          UUIDToString(row.ID),
		Email:       row.Email,
		DisplayName: row.DisplayName,
	}
}

func ProjectFromSQL(p sqlc.Project) Project {
	return Project{
		ID:        UUIDToString(p.ID),
		Title:     p.Title,
		OwnerID:   UUIDToString(p.OwnerID),
		CreatedAt: Timestamptz(p.CreatedAt),
	}
}

func TaskFromSQL(t sqlc.Task) Task {
	return taskFromFields(t.ID, t.ProjectID, t.Title, t.Status, t.AssigneeID, t.X, t.Y, t.UpdatedAt)
}

func TaskFromListRow(row sqlc.ListTasksByProjectRow) Task {
	return taskFromFields(row.ID, row.ProjectID, row.Title, row.Status, row.AssigneeID, row.X, row.Y, row.UpdatedAt)
}

func TaskFromCreateRow(row sqlc.CreateTaskRow) Task {
	return taskFromFields(row.ID, row.ProjectID, row.Title, row.Status, row.AssigneeID, row.X, row.Y, row.UpdatedAt)
}

func TaskFromGetRow(row sqlc.GetTaskRow) Task {
	return taskFromFields(row.ID, row.ProjectID, row.Title, row.Status, row.AssigneeID, row.X, row.Y, row.UpdatedAt)
}

func TaskFromUpdateRow(row sqlc.UpdateTaskRow) Task {
	return taskFromFields(row.ID, row.ProjectID, row.Title, row.Status, row.AssigneeID, row.X, row.Y, row.UpdatedAt)
}

func taskFromFields(
	id, projectID uuid.UUID,
	title, status string,
	assigneeID *uuid.UUID,
	x, y float64,
	updatedAt pgtype.Timestamptz,
) Task {
	return Task{
		ID:         UUIDToString(id),
		ProjectID:  UUIDToString(projectID),
		Title:      title,
		Status:     status,
		AssigneeID: OptionalUUIDString(assigneeID),
		X:          x,
		Y:          y,
		UpdatedAt:  Timestamptz(updatedAt),
	}
}

func DependencyFromSQL(d sqlc.Dependency) Dependency {
	return Dependency{
		ID:         UUIDToString(d.ID),
		ProjectID:  UUIDToString(d.ProjectID),
		FromTaskID: UUIDToString(d.FromTaskID),
		ToTaskID:   UUIDToString(d.ToTaskID),
	}
}

func MemberFromSQL(row sqlc.ListProjectMembersRow) ProjectMember {
	return ProjectMember{
		UserID:      UUIDToString(row.ID),
		Email:       row.Email,
		DisplayName: row.DisplayName,
		Role:        row.Role,
	}
}

func EdgeFromSQL(row sqlc.ListDependencyEdgesByProjectRow) GraphEdge {
	return GraphEdge{
		FromTaskID: UUIDToString(row.FromTaskID),
		ToTaskID:   UUIDToString(row.ToTaskID),
	}
}
