package service

import (
	"database/sql"

	"github.com/cyberhawk12121/Saarthi/internal/repository"
	types "github.com/cyberhawk12121/Saarthi/internal/shared"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	DB            *sql.DB
	userRepo      *repository.UserRepository
	recordingRepo *repository.RecordingRepository
	config        *Config
}

func NewUserService(db *sql.DB) *UserService {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	return &UserService{DB: db, userRepo: repository.NewUserRepository(db), recordingRepo: repository.NewRecordingRepository(db), config: config}
}

func (us *UserService) RegisterUser(userData types.RegisterRequest) error {
	// 1. Check if the email is already registered - DONE
	// 2. convert password to a bcrypt hash - DONE
	// 3. use jwt library and create a jwt token for the user and return the object
	// 4. put the refresh token in the db and also in the cookie header
	// 5. return the user object with the jwt token

	_, err := us.userRepo.GetUserByEmail(userData.Email)
	if err != nil {
		return err
	}

	// For now just save in DB
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	rows, err := us.DB.Query("INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING true", userData.FirstName, userData.LastName, userData.Email, bcryptPassword)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (us *UserService) LoginUser(req types.LoginRequest) (types.LoginResponse, error) {
	// 1. check if the user exists in the db
	// 2. if the user exists then check if the password matches
	// 3. if the password matches then create a jwt token and return the object
	// 4. if the password doesn't match then return an error
	bcryptedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return types.LoginResponse{}, err
	}
	rows, err := us.DB.Query("SELECT * FROM users WHERE email=$1 & password=$2", req.Email, bcryptedPassword)

	if err != nil {
		return types.LoginResponse{}, err
	}

	defer rows.Close()

	return types.LoginResponse{}, nil
}
