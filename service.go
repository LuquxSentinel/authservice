package main

import (
	"context"

	"github.com/luqus/authservice/types"
)

type Service interface {
	Login(ctx context.Context, email, password string) (*types.ResponseUser, error)
	CreateUser(ctx context.Context, createUserInput *types.CreateUserInput) error
}
