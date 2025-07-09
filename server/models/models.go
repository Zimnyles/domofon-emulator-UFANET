package models

import "time"

type LoginForm struct {
	Login    string
	Password string
}

type RegistrationForm struct {
	Login      string
	Email      string
	Password   string
	SecretCode string
}

type CreateUserCredential struct {
	Login    string
	Email    string
	Password string
}

type CreateIntercomCredentials struct {
	MAC                string `json:"mac"`
	Address            string `json:"address"`
	NumberOfApartments int    `json:"apartments"`
}

type CreateIntercomConnectionRequset struct {
	ID int `json:"id"`
}

type Intercom struct {
	ID                  int       `json:"id"`
	MAC                 string    `json:"mac_address"`
	IntercomStatus      bool      `json:"intercom_status"`
	DoorStatus          bool      `json:"door_status"`
	Address             string    `json:"address"`
	NumberOfApartments  int       `json:"number_of_apartments"`
	IsCalling           bool      `json:"is_calling"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	CalledApartment     int
	OpenedDoorApartment int
	IsActive            bool `json:"is_active"`
}
