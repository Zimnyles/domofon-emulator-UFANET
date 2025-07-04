package models

import (
	"fmt"
	"strings"
	"time"
)

type NewIntercomProperties struct {
	MAC                string
	Address            string
	NumberOfApartments int
}

type Intercom struct {
	ID                 int       `json:"id"`
	MAC                string    `json:"mac_address"`
	DomofonStatus      bool      `json:"domofon_status"`
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
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	Intercom *Intercom `json:"intercom,omitempty"`
}

type CreateIntercomResponse struct {
	Success bool
	ID      int
	IsNew   bool
	Message string
}

func (d *Intercom) Topic() string {
	return fmt.Sprintf("%s/status", strings.ReplaceAll(d.MAC, ":", ""))
}
