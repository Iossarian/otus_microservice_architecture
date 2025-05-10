package domain

type UserCreatedEvent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
