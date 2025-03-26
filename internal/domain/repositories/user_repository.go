package repositories

import (
	"errors"
	"sync"
	"time"

	"github.com/amirtalbi/examen_go/internal/domain/models"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id string) (*models.User, error)
	SaveResetToken(email, token string, expiry time.Time) error
	FindByResetToken(token string) (*models.User, error)
	UpdatePassword(id, password string) error
}

type inMemoryUserRepository struct {
	users map[string]*models.User
	mutex sync.RWMutex
}

func NewUserRepository() UserRepository {
	return &inMemoryUserRepository{
		users: make(map[string]*models.User),
	}
}

func (r *inMemoryUserRepository) Create(user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, existingUser := range r.users {
		if existingUser.Email == user.Email {
			return ErrEmailAlreadyExists
		}
	}

	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	r.users[user.ID] = user
	return nil
}

func (r *inMemoryUserRepository) FindByEmail(email string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

func (r *inMemoryUserRepository) FindByID(id string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if user, exists := r.users[id]; exists {
		return user, nil
	}
	return nil, ErrUserNotFound
}

func (r *inMemoryUserRepository) SaveResetToken(email, token string, expiry time.Time) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, user := range r.users {
		if user.Email == email {
			tokenCopy := token
			user.ResetToken = &tokenCopy
			expiryCopy := expiry
			user.ResetTokenExpires = &expiryCopy
			user.UpdatedAt = time.Now()
			return nil
		}
	}
	return ErrUserNotFound
}

func (r *inMemoryUserRepository) FindByResetToken(token string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.ResetToken != nil && *user.ResetToken == token && 
		   user.ResetTokenExpires != nil && user.ResetTokenExpires.After(time.Now()) {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

func (r *inMemoryUserRepository) UpdatePassword(id, password string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if user, exists := r.users[id]; exists {
		user.Password = password
		user.ResetToken = nil
		user.ResetTokenExpires = nil
		user.UpdatedAt = time.Now()
		return nil
	}
	return ErrUserNotFound
}
