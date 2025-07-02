package models

type LoginForm struct{
	Login string
	Password string
}

type RegistrationForm struct{
	Login string
	Email string
	Password string
	SecretCode string
}

type CreateUserCredential struct{
	Login string
	Email string
	Password string
}