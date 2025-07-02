package auth

import (
	"domofonEmulator/pkg/tadapter"
	"domofonEmulator/pkg/validator"
	"domofonEmulator/server/models"
	"domofonEmulator/server/web/views/components"
	"domofonEmulator/server/web/views/pages"
	"net/http"
	"regexp"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	router     fiber.Router
	logger     *zerolog.Logger
	service    IAuthService
	repository IAuthRepository
	store      *session.Store
}

type IAuthRepository interface {
	IsLoginExists(login string) (bool, error)
	IsEmailExists(email string) (bool, error)
	AddUser(form models.CreateUserCredential, logger *zerolog.Logger) (bool, error)
	GetPasswordByLogin(login string) (string, error)
}

type IAuthService interface {
	GetRegistrationKey() string
	LoginUser(loginForm models.LoginForm) (bool, string)
}

func NewHandler(router fiber.Router, logger *zerolog.Logger, mqqtClient mqtt.Client, service IAuthService, repository IAuthRepository, store *session.Store) {
	h := &AuthHandler{
		router:     router,
		logger:     logger,
		service:    service,
		repository: repository,
		store:      store,
	}

	h.router.Get("/login", h.login)
	h.router.Get("/register", h.register)

	h.router.Post("/api/login", h.apiLogin)
	h.router.Post("/api/register", h.apiRegister)

}

func (h *AuthHandler) login(c *fiber.Ctx) error {
	component := pages.LoginPage()
	return tadapter.Render(c, component, fiber.StatusOK)
}

func (h *AuthHandler) register(c *fiber.Ctx) error {
	component := pages.RegisterPage()
	return tadapter.Render(c, component, fiber.StatusOK)
}

func (h *AuthHandler) apiLogin(c *fiber.Ctx) error {
	loginForm := models.LoginForm{
		Login:    c.FormValue("login"),
		Password: c.FormValue("password"),
	}

	isSuccessToLoginUser, msg := h.service.LoginUser(loginForm)
	if isSuccessToLoginUser {
		sess, err := h.store.Get(c)
		if err != nil {
			h.logger.Fatal().Err(err).Msg("Failed to get session store")
			panic(err)
		}
		sess.Set("login", strings.ToLower(loginForm.Login))
		if err := sess.Save(); err != nil {
			h.logger.Fatal().Err(err).Msg("Failed to set user session")
			panic(err)
		}
		c.Response().Header.Add("Hx-Redirect", "/")
		return c.Redirect("/", http.StatusOK)
	}
	component := components.Notification(msg, false)
	return tadapter.Render(c, component, fiber.StatusOK)
}

func (h *AuthHandler) apiRegister(c *fiber.Ctx) error {
	registerForm := models.RegistrationForm{
		Login:      c.FormValue("login"),
		Email:      c.FormValue("email"),
		Password:   c.FormValue("password"),
		SecretCode: c.FormValue("secretcode"),
	}
	//key validation
	envSecretKey := h.service.GetRegistrationKey()
	if envSecretKey != registerForm.SecretCode {
		component := components.Notification("Неверный код приглашение, обратитесь к вашему руководителю", false)
		return tadapter.Render(c, component, fiber.StatusOK)
	}
	//login validation
	IsLoginExists, _ := h.repository.IsLoginExists(registerForm.Login)
	if IsLoginExists {
		component := components.Notification("Пользователя с таким логином уже существует", false)
		return tadapter.Render(c, component, fiber.StatusOK)
	}
	if len(registerForm.Login) > 13 {
		component := components.Notification("Длина логина не может превышать 13 символов", false)
		return tadapter.Render(c, component, fiber.StatusOK)
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, registerForm.Login)
	if !matched {
		component := components.Notification("Логин может содержать только буквы, цифры и подчеркивания", false)
		return tadapter.Render(c, component, http.StatusBadRequest)
	}
	//email validation
	IsEmailExists, _ := h.repository.IsEmailExists(registerForm.Email)
	if IsEmailExists {
		component := components.Notification("Пользователь с такой почтой уже существует", false)
		return tadapter.Render(c, component, fiber.StatusOK)
	}

	//password validation
	if len(registerForm.Password) < 6 {
		component := components.Notification("Длина пароля не может быть меньше 6 символов", false)
		return tadapter.Render(c, component, fiber.StatusOK)
	}

	//overall validation
	errors := validate.Validate(
		&validators.EmailIsPresent{Name: "Email", Field: registerForm.Email, Message: "Email не задан или задан неверно"},
		&validators.StringIsPresent{Name: "Password", Field: registerForm.Password, Message: "Пароль не задан или задан неверно"},
		&validators.StringIsPresent{Name: "Login", Field: registerForm.Login, Message: "Логин не задан или задан неверно"},
	)
	if len(errors.Error()) > 0 {
		component := components.Notification(validator.FormatErrors(errors), false)
		return tadapter.Render(c, component, fiber.StatusOK)
	}

	createUserForm := models.CreateUserCredential{
		Login:    registerForm.Login,
		Email:    registerForm.Email,
		Password: registerForm.Password,
	}

	success, err := h.repository.AddUser(createUserForm, h.logger)
	if err != nil {
		h.logger.Fatal().Err(err).Msg("Failed to create new account")
	}

	sess, err := h.store.Get(c)
	if err != nil {
		h.logger.Fatal().Err(err).Msg("Failed to get session store")
		panic(err)
	}
	sess.Set("login", strings.ToLower(registerForm.Login))
	if err := sess.Save(); err != nil {
		h.logger.Fatal().Err(err).Msg("Failed to set user session")
		panic(err)
	}

	if success {
		c.Response().Header.Add("Hx-Redirect", "/")
		return c.Redirect("/", http.StatusOK)
	}

	component := components.Notification("При регистрации произошла ошибка на сервере. Обратитесть к системному администратору", false)
	return tadapter.Render(c, component, fiber.StatusOK)
}
