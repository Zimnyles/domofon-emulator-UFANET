package auth

import (
	"domofonEmulator/server/models"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	logger     *zerolog.Logger
	mqqtClient mqtt.Client
	repository IAuthRepository
}

func NewAuthService(logger *zerolog.Logger, mqqtClient mqtt.Client, repository IAuthRepository) *AuthService {
	return &AuthService{
		logger:     logger,
		mqqtClient: mqqtClient,
		repository: repository,
	}
}

func (s *AuthService) GetRegistrationKey() string {
	return getString("REG_KEY", "")
}

func getString(key, defValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defValue
	}
	return value
}

func (s *AuthService) LoginUser(loginForm models.LoginForm) (bool, string) {
	IsLoginExists, _ := s.repository.IsLoginExists(loginForm.Login)
	if !IsLoginExists {
		return false, "Пользователя с таким логином не существует"
	}

	userHashedPassword, err := s.repository.GetPasswordByLogin(loginForm.Login)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("Failed to get user hashed password from user login")
		return false, "Не удалось авторизоваться. Обратитесь к системному администратору"
	}

	err = bcrypt.CompareHashAndPassword([]byte(userHashedPassword), []byte(loginForm.Password))
	if err != nil {
		return false, "Неверный пароль. Попробуйте еще раз"
	}

	return true, "success"

}
