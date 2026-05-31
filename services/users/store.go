package users

import (
	"database/sql"

	"github.com/YhomiAce/rest-api-pg/types"
	"github.com/YhomiAce/rest-api-pg/config"
)

type Store struct {
	db *sql.DB
	config *config.Config
}

func NewUserStore(db *sql.DB, config *config.Config) *Store {
	return &Store{db: db, config: config}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	var user types.User
	query := "SELECT * FROM users WHERE email = $1"
	err := s.db.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, nil
}

func (s *Store) CreateUser(user *types.CreateUserPayload) (error) {
	query := "INSERT INTO users (username, email, password, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id"
	_, err := s.db.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}