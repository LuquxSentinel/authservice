package types

import "time"

type User struct {
	UID       string    `json:"uid" bson:"uid"`
	FirstName string    `json:"first_name" bson:"first_name"`
	LastName  string    `json:"last_name" bson:"last_name"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"password" bson:"password"`
	Busket    *Busket   `json:"busket" bson:"busket"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
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
	FirstName string    `json:"first_name"`
	LastName  string    `json:"lastname"`
	Email     string    `json:"email"`
	Busket    *Busket   `json:"busket"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
