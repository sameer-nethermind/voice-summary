package repository

import (
	"database/sql"

	"github.com/cyberhawk12121/Saarthi/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (rs *UserRepository) NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) GetAllUsers() ([]model.User, error) {
	rows, err := ur.db.Query(`
		SELECT id, first_name, last_name, email, password, created_at, updated_at 
		FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, rows.Err() // Check for row iteration errors
}

func (ur *UserRepository) GetUserByEmail(email string) (model.User, error) {
	rows, err := ur.db.Query("SELECT * FROM users WHERE email=$1", email)
	if err != nil {
		return model.User{}, err
	}

	defer rows.Close()

	// Check if rows are empty or not
	if !rows.Next() {
		return model.User{}, nil
	}

	var user model.User
	if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (ur *UserRepository) CreateUser(user model.User) error {
	_, err := ur.db.Exec(`
		INSERT INTO users (first_name, last_name, email, password) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at, updated_at
	`, user.FirstName, user.LastName, user.Email, user.Password)
	return err
}
