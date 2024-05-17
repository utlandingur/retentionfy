package memory

import (
	"errors"
	"sync"
	"time"

	"github.com/dchest/uniuri"
	"github.com/noona-hq/app-template/store"
	"github.com/noona-hq/app-template/store/entity"
)

type inMemoryStore struct {
	users map[string][]entity.User // Keyed by companyID
	mu    sync.RWMutex
}

// NewStore creates a new in-memory store.
// Useful for testing and development.
func NewStore() store.Store {
	return &inMemoryStore{
		users: make(map[string][]entity.User),
	}
}

// CreateUser adds a new user to the in-memory store.
func (s *inMemoryStore) CreateUser(user entity.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if user.ID == "" {
		user.ID = randomID()
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Append the user to the slice associated with the company ID.
	s.users[user.CompanyID] = append(s.users[user.CompanyID], user)

	return nil
}

func (s *inMemoryStore) UpdateUser(id string, user entity.User) (entity.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user.UpdatedAt = time.Now()

	users, exists := s.users[user.CompanyID]
	if !exists {
		return entity.User{}, errors.New("user not found")
	}

	for i, u := range users {
		if u.ID == id {
			users[i] = user
			s.users[user.CompanyID] = users
			return user, nil
		}
	}

	return entity.User{}, errors.New("user not found")
}

// GetUserForCompany retrieves the latest user for a given company ID from the in-memory store.
func (s *inMemoryStore) GetUserForCompany(companyID string) (entity.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users, exists := s.users[companyID]
	if !exists || len(users) == 0 {
		return entity.User{}, errors.New("user not found")
	}

	// Since we're always appending, the last user should be the latest one.
	latestUser := users[len(users)-1]

	return latestUser, nil
}

func (s *inMemoryStore) DeleteUser(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for companyID, users := range s.users {
		for i, user := range users {
			if user.ID == id {
				// Remove the user from the slice.
				s.users[companyID] = append(users[:i], users[i+1:]...)
				return nil
			}
		}
	}

	return errors.New("user not found")
}

func randomID() string {
	return uniuri.NewLen(24)
}
