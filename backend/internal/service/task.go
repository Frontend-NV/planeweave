package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/planeweave/planeweave/internal/domain"
	"github.com/planeweave/planeweave/internal/repository"
	"github.com/planeweave/planeweave/internal/ws"
)

type EventPublisher interface {
	Broadcast(projectID string, event ws.Event)
}

type TaskService struct {
	repo   *repository.Repository
	access *AccessService
	events EventPublisher
}

func NewTaskService(repo *repository.Repository, access *AccessService, events EventPublisher) *TaskService {
	return &TaskService{repo: repo, access: access, events: events}
}

type CreateTaskInput struct {
	Title      string
	Status     string
	AssigneeID *string
	X          float64
	Y          float64
}

func (s *TaskService) List(ctx context.Context, projectID, userID uuid.UUID) ([]domain.Task, error) {
	if err := s.access.EnsureMember(ctx, projectID, userID); err != nil {
		return nil, err
	}
	return s.repo.ListTasks(ctx, projectID)
}

func (s *TaskService) Create(ctx context.Context, projectID, userID uuid.UUID, input CreateTaskInput) (domain.Task, error) {
	if input.Title == "" {
		return domain.Task{}, domain.ErrInvalidRequest
	}
	if err := s.access.EnsureMember(ctx, projectID, userID); err != nil {
		return domain.Task{}, err
	}

	status := input.Status
	if status == "" {
		status = "todo"
	}
	assigneeID, err := domain.OptionalStringUUID(input.AssigneeID)
	if err != nil {
		return domain.Task{}, domain.ErrInvalidRequest
	}

	task, err := s.repo.CreateTask(ctx, projectID, input.Title, status, assigneeID, input.X, input.Y)
	if err != nil {
		return domain.Task{}, err
	}
	s.events.Broadcast(projectID.String(), ws.Event{Type: "task.created", Payload: task})
	return task, nil
}

func (s *TaskService) Update(ctx context.Context, taskID, userID uuid.UUID, patch domain.TaskPatch) (domain.Task, string, error) {
	current, err := s.repo.GetTask(ctx, taskID)
	if err != nil {
		return domain.Task{}, "", err
	}

	projectID, err := domain.StringToUUID(current.ProjectID)
	if err != nil {
		return domain.Task{}, "", domain.ErrInvalidRequest
	}
	if err := s.access.EnsureMember(ctx, projectID, userID); err != nil {
		return domain.Task{}, "", err
	}

	title := current.Title
	if patch.Title != nil {
		title = *patch.Title
	}
	status := current.Status
	if patch.Status != nil {
		status = *patch.Status
	}
	assigneeID := current.AssigneeID
	if patch.AssigneeID != nil {
		assigneeID = patch.AssigneeID
	}
	x := current.X
	if patch.X != nil {
		x = *patch.X
	}
	y := current.Y
	if patch.Y != nil {
		y = *patch.Y
	}

	assigneeUUID, err := domain.OptionalStringUUID(assigneeID)
	if err != nil {
		return domain.Task{}, "", domain.ErrInvalidRequest
	}

	task, err := s.repo.UpdateTask(ctx, taskID, title, status, assigneeUUID, x, y)
	if err != nil {
		return domain.Task{}, "", err
	}

	eventType := "task.updated"
	if patch.X != nil || patch.Y != nil {
		eventType = "task.moved"
	}
	s.events.Broadcast(projectID.String(), ws.Event{Type: eventType, Payload: task})
	return task, eventType, nil
}

func (s *TaskService) Delete(ctx context.Context, taskID, userID uuid.UUID) error {
	projectID, err := s.repo.GetTaskProjectID(ctx, taskID)
	if err != nil {
		return err
	}
	if err := s.access.EnsureMember(ctx, projectID, userID); err != nil {
		return err
	}
	if err := s.repo.DeleteTask(ctx, taskID); err != nil {
		return err
	}
	s.events.Broadcast(projectID.String(), ws.Event{
		Type: "task.deleted",
		Payload: map[string]string{
			"id":        taskID.String(),
			"projectId": projectID.String(),
		},
	})
	return nil
}
