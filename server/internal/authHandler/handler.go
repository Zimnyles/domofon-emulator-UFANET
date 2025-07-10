package auth

import (
	"domofonEmulator/pkg/tadapter"
	"domofonEmulator/server/models"
	mqttserver "domofonEmulator/server/mqttServer"
	"domofonEmulator/server/web/views/components"
	"domofonEmulator/server/web/views/pages"
	"net/http"
	"strings"

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
	GetIntercomByID(id int, logger *zerolog.Logger) (models.Intercom, error)
}

type IAuthService interface {
	GetRegistrationKey() string
	LoginUser(loginForm models.LoginForm) (bool, string)
	RegisterUser(registerForm models.RegistrationForm) (bool, string)
}

func NewHandler(router fiber.Router, logger *zerolog.Logger, mqttServer mqttserver.Server, service IAuthService, repository IAuthRepository, store *session.Store) {
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
	h.router.Get("api/logout", h.apiLogout)

}

func (h *AuthHandler) apiLogout(c *fiber.Ctx) error {

	sess, err := h.store.Get(c)
	if err != nil {
		panic(err)
	}

	sess.Delete("login")

	if err := sess.Save(); err != nil {
		panic(err)
	}

	c.Response().Header.Add("Hx-Redirect", "/login")
	return c.Redirect("/", http.StatusOK)

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
		c.Response().Header.Add("Hx-Redirect", "/connect")
		return c.Redirect("/", http.StatusOK)
	}
	component := components.NewIntercomResponse(msg)
	return tadapter.Render(c, component, fiber.StatusOK)
}

func (h *AuthHandler) apiRegister(c *fiber.Ctx) error {
	registerForm := models.RegistrationForm{
		Login:      c.FormValue("login"),
		Email:      c.FormValue("email"),
		Password:   c.FormValue("password"),
		SecretCode: c.FormValue("secretcode"),
	}

	isRegistrated, msg := h.service.RegisterUser(registerForm)
	if isRegistrated {
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
		c.Response().Header.Add("Hx-Redirect", "/connect")
		return c.Redirect("/", http.StatusOK)
	}
	component := components.NewIntercomResponse(msg)
	return tadapter.Render(c, component, fiber.StatusOK)
}
