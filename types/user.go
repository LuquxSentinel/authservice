package types

import "time"

type User struct {
	UID       string
	FirstName string
	LastName  string
	Email     string
	Password  string
	Busket    *Busket
	CreatedAt time.Time
}

func (u *User) ResponseUser() *ResponseUser {
	return &ResponseUser{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Busket:    u.Busket,
		CreatedAt: u.CreatedAt,
	}
}

type ResponseUser struct {
	FirstName string
	LastName  string
	Email     string
	Busket    *Busket
	CreatedAt time.Time
}

type CreateUserInput struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type LoginInput struct {
	Email    string
	Password string
}
