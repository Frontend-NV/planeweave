package domain

import "time"

type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type Project struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	OwnerID   string    `json:"ownerId"`
	CreatedAt time.Time `json:"createdAt"`
}

type ProjectMember struct {
	UserID      string `json:"userId"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	Role        string `json:"role"`
}

type Task struct {
	ID         string    `json:"id"`
	ProjectID  string    `json:"projectId"`
	Title      string    `json:"title"`
	Status     string    `json:"status"`
	AssigneeID *string   `json:"assigneeId"`
	X          float64   `json:"x"`
	Y          float64   `json:"y"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type Dependency struct {
	ID         string `json:"id"`
	ProjectID  string `json:"projectId"`
	FromTaskID string `json:"fromTaskId"`
	ToTaskID   string `json:"toTaskId"`
}

type Invite struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type TaskPatch struct {
	Title      *string
	Status     *string
	AssigneeID *string
	X          *float64
	Y          *float64
}

type GraphEdge struct {
	FromTaskID string
	ToTaskID   string
}
