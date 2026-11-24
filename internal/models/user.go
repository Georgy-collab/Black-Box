package models

// User is the stored/response user entity.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// UserCreate is the payload for creating a user.
type UserCreate struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// UserUpdate is the payload for partial user updates.
// Pointer fields distinguish "not provided" from zero values.
type UserUpdate struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Age   *int    `json:"age"`
}

// HealthResponse is returned by GET /health.
type HealthResponse struct {
	Status string `json:"status"`
}

// ErrorResponse matches the FastAPI-style error body.
type ErrorResponse struct {
	Detail string `json:"detail"`
}
