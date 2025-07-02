package auth

import (
	"context"
	"domofonEmulator/server/models"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository struct {
	Dbpool       *pgxpool.Pool
	CustomLogger *zerolog.Logger
}

func NewAuthRepository(dbpool *pgxpool.Pool, customLogger *zerolog.Logger) *AuthRepository {
	return &AuthRepository{
		Dbpool:       dbpool,
		CustomLogger: customLogger,
	}
}

func (r *AuthRepository) IsLoginExists(login string) (bool, error) {
	var exists bool
	err := r.Dbpool.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)",
		login,
	).Scan(&exists)

	return exists, err
}

func (r *AuthRepository) IsEmailExists(email string) (bool, error) {
	var exists bool
	err := r.Dbpool.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)",
		email,
	).Scan(&exists)

	return exists, err
}

func (r *AuthRepository) AddUser(form models.CreateUserCredential, logger *zerolog.Logger) (bool, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}
	query := `
		INSERT INTO users (email, login, hashed_password, created_at) 
		VALUES (@email, @login, @hashed_password, @createdat)
	`
	args := pgx.NamedArgs{
		"email":           form.Email,
		"login":           strings.ToLower(form.Login),
		"hashed_password": hashedPassword,
		"created_at":      time.Now(),
	}
	_, err = r.Dbpool.Exec(context.Background(), query, args)
	if err != nil {
		return false, err
	}
	logger.Info().
		Str("login", form.Login).
		Str("email", form.Email).
		Time("created_at", time.Now()).
		Msg("New account created")
	return true, nil
}

func (r *AuthRepository) GetPasswordByLogin(login string) (string, error) {
	query := `
        SELECT   
            hashed_password   
        FROM users 
        WHERE login = @login`
	args := pgx.NamedArgs{
    	"login": login,
	}
	row := r.Dbpool.QueryRow(context.Background(), query, args)
	var userHashedPassword string
	err := row.Scan(&userHashedPassword)
	if err != nil {
    	return "", err		
	}
	return userHashedPassword, nil

}

