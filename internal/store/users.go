package store

import (
	"errors"
	"sort"
	"sync"

	"github.com/Georgy-collab/Black-Box/internal/models"
)

var (
	ErrNotFound    = errors.New("user not found")
	ErrEmailExists = errors.New("email already exists")
)

// UserStore is an in-memory thread-safe user repository.
type UserStore struct {
	mu     sync.RWMutex
	users  map[int]models.User
	nextID int
}

func NewUserStore() *UserStore {
	return &UserStore{
		users:  make(map[int]models.User),
		nextID: 1,
	}
}

func (s *UserStore) Create(name, email string, age int) (models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.ensureUniqueEmailLocked(email, 0); err != nil {
		return models.User{}, err
	}

	user := models.User{
		ID:    s.nextID,
		Name:  name,
		Email: email,
		Age:   age,
	}
	s.users[user.ID] = user
	s.nextID++
	return user, nil
}

func (s *UserStore) List(minAge *int) []models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.User, 0, len(s.users))
	for _, user := range s.users {
		if minAge != nil && user.Age < *minAge {
			continue
		}
		result = append(result, user)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result
}

func (s *UserStore) Get(id int) (models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return models.User{}, ErrNotFound
	}
	return user, nil
}

func (s *UserStore) Update(id int, name, email *string, age *int) (models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[id]
	if !ok {
		return models.User{}, ErrNotFound
	}

	if name != nil {
		user.Name = *name
	}
	if email != nil {
		if err := s.ensureUniqueEmailLocked(*email, id); err != nil {
			return models.User{}, err
		}
		user.Email = *email
	}
	if age != nil {
		user.Age = *age
	}

	s.users[id] = user
	return user, nil
}

func (s *UserStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[id]; !ok {
		return ErrNotFound
	}
	delete(s.users, id)
	return nil
}

func (s *UserStore) ensureUniqueEmailLocked(email string, excludeID int) error {
	for id, user := range s.users {
		if excludeID != 0 && id == excludeID {
			continue
		}
		if user.Email == email {
			return ErrEmailExists
		}
	}
	return nil
}
