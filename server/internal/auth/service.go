package auth

import (
	"domofonEmulator/pkg/validator"
	"domofonEmulator/server/models"
	"os"
	"regexp"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
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

func (s *AuthService) RegisterUser(registerForm models.RegistrationForm) (bool, string) {
	envSecretKey := s.GetRegistrationKey()
	if envSecretKey != registerForm.SecretCode {
		return false, "Неверный код приглашение, обратитесь к вашему руководителю"
	}
	//login validation
	IsLoginExists, _ := s.repository.IsLoginExists(registerForm.Login)
	if IsLoginExists {
		return false, "Пользователя с таким логином уже существует"
	}
	if len(registerForm.Login) > 13 {
		return false, "Длина логина не может превышать 13 символов"
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, registerForm.Login)
	if !matched {
		return false, "Логин может содержать только буквы, цифры и подчеркивания"
	}
	//email validation
	IsEmailExists, _ := s.repository.IsEmailExists(registerForm.Email)
	if IsEmailExists {
		return false, "Пользователь с такой почтой уже существует"
	}

	//password validation
	if len(registerForm.Password) < 6 {
		return false, "Длина пароля не может быть меньше 6 символов"
	}

	//overall validation
	errors := validate.Validate(
		&validators.EmailIsPresent{Name: "Email", Field: registerForm.Email, Message: "Email не задан или задан неверно"},
		&validators.StringIsPresent{Name: "Password", Field: registerForm.Password, Message: "Пароль не задан или задан неверно"},
		&validators.StringIsPresent{Name: "Login", Field: registerForm.Login, Message: "Логин не задан или задан неверно"},
	)
	if len(errors.Error()) > 0 {
		return false, validator.FormatErrors(errors)
	}

	createUserForm := models.CreateUserCredential{
		Login:    registerForm.Login,
		Email:    registerForm.Email,
		Password: registerForm.Password,
	}

	success, err := s.repository.AddUser(createUserForm, s.logger)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("Failed to create new account")
	}

	return success, ""
}
