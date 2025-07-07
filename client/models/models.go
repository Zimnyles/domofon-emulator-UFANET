package models

import (
	"time"
)

type NewIntercomProperties struct {
	MAC                string
	Address            string
	NumberOfApartments int
}

type IntercomConnectResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Intercom
}

type Intercom struct {
	ID                 int       `json:"id"`
	MAC                string    `json:"mac_address"`
	IntercomStatus     bool      `json:"intercom_status"`
	DoorStatus         bool      `json:"door_status"`
	Address            string    `json:"address"`
	NumberOfApartments int       `json:"number_of_apartments"`
	IsCalling          bool      `json:"is_calling"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type RegistrationRequest struct {
	Action             string `json:"action"`
	MAC                string `json:"mac"`
	Address            string `json:"address"`
	NumberOfApartments int    `json:"apartments"`
}

type RegistrationResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	Intercom *Intercom `json:"intercom,omitempty"`
}

type CreateIntercomResponse struct {
	Success bool
	ID      int
	IsNew   bool
	Mac     string
	Message string
}

type IntercomPowerOnOffRequest struct {
	ID     int
	Action string
}

type IntercomPowerOnOffResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Intercom `json:"intercom,omitempty"`
}

type IntercomOpenDoorResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Intercom `json:"intercom,omitempty"`
}

type IntercomCallResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Intercom `json:"intercom,omitempty"`
}
