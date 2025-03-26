package repositories

import (
	"log"
	"time"

	"github.com/amirtalbi/examen_go/internal/domain/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type postgresUserRepository struct {
	db *sqlx.DB
}

func NewPostgresUserRepository(db *sqlx.DB) UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(user *models.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
        INSERT INTO users (id, name, email, password, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := r.db.Exec(query, user.ID, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)

	return err
}

func (r *postgresUserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE email = $1"
	
	log.Printf("Searching for user with email: %s", email)
	err := r.db.Get(&user, query, email)
	if err != nil {
		log.Printf("Error finding user by email: %v", err)
		return nil, ErrUserNotFound
	}
	log.Printf("Found user: %s with ID: %s", user.Email, user.ID)
	return &user, nil
}

func (r *postgresUserRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE id = $1"
	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (r *postgresUserRepository) SaveResetToken(email, token string, expiry time.Time) error {
	query := `
        UPDATE users 
        SET reset_token = $1, reset_token_expires = $2, updated_at = $3
        WHERE email = $4
    `
	result, err := r.db.Exec(query, token, expiry, time.Now(), email)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *postgresUserRepository) FindByResetToken(token string) (*models.User, error) {
	var user models.User
	query := `
        SELECT * FROM users 
        WHERE reset_token = $1 AND (reset_token_expires IS NULL OR reset_token_expires > $2)
    `
	err := r.db.Get(&user, query, token, time.Now())
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (r *postgresUserRepository) UpdatePassword(id, password string) error {
	query := `
        UPDATE users 
        SET password = $1, reset_token = NULL, reset_token_expires = NULL, updated_at = $2
        WHERE id = $3
    `
	result, err := r.db.Exec(query, password, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
