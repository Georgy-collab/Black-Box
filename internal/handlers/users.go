package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Georgy-collab/Black-Box/internal/models"
	"github.com/Georgy-collab/Black-Box/internal/pyvalidate"
	"github.com/Georgy-collab/Black-Box/internal/store"
)

type UserHandler struct {
	store *store.UserStore
}

func NewUserHandler(s *store.UserStore) *UserHandler {
	return &UserHandler{store: s}
}

func (h *UserHandler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, models.HealthResponse{Status: "ok"})
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var body models.UserCreate
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusUnprocessableEntity, "Invalid request body")
		return
	}

	validated, err := pyvalidate.Validate(map[string]any{
		"name":  body.Name,
		"email": body.Email,
		"age":   body.Age,
	})
	if err != nil {
		writeValidationError(w, err)
		return
	}
	if validated.Name == nil || validated.Email == nil || validated.Age == nil {
		writeError(w, http.StatusUnprocessableEntity, "Incomplete validation result")
		return
	}

	user, err := h.store.Create(*validated.Name, *validated.Email, *validated.Age)
	if err != nil {
		if errors.Is(err, store.ErrEmailExists) {
			writeError(w, http.StatusBadRequest, "Email already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	var minAge *int
	if raw := r.URL.Query().Get("min_age"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value < 0 {
			writeError(w, http.StatusUnprocessableEntity, "Invalid min_age")
			return
		}
		minAge = &value
	}

	users := h.store.List(minAge)
	writeJSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUserID(w, r)
	if !ok {
		return
	}

	user, err := h.store.Get(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "User not found")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUserID(w, r)
	if !ok {
		return
	}

	var body models.UserUpdate
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusUnprocessableEntity, "Invalid request body")
		return
	}

	fields := map[string]any{}
	if body.Name != nil {
		fields["name"] = *body.Name
	}
	if body.Email != nil {
		fields["email"] = *body.Email
	}
	if body.Age != nil {
		fields["age"] = *body.Age
	}

	var name, email *string
	var age *int
	if len(fields) > 0 {
		validated, err := pyvalidate.Validate(fields)
		if err != nil {
			writeValidationError(w, err)
			return
		}
		name = validated.Name
		email = validated.Email
		age = validated.Age
	}

	user, err := h.store.Update(id, name, email, age)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			writeError(w, http.StatusNotFound, "User not found")
		case errors.Is(err, store.ErrEmailExists):
			writeError(w, http.StatusBadRequest, "Email already exists")
		default:
			writeError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUserID(w, r)
	if !ok {
		return
	}

	if err := h.store.Delete(id); err != nil {
		writeError(w, http.StatusNotFound, "User not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseUserID(w http.ResponseWriter, r *http.Request) (int, bool) {
	raw := r.PathValue("user_id")
	id, err := strconv.Atoi(raw)
	if err != nil || id < 1 {
		writeError(w, http.StatusUnprocessableEntity, "Invalid user id")
		return 0, false
	}
	return id, true
}

func writeValidationError(w http.ResponseWriter, err error) {
	var validationErr *pyvalidate.ValidationError
	if errors.As(err, &validationErr) {
		writeError(w, http.StatusUnprocessableEntity, validationErr.Detail)
		return
	}
	writeError(w, http.StatusInternalServerError, "Validator unavailable")
}
