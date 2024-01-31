package main

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/luqus/authservice/types"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Login(ctx context.Context, email, password string) (*types.ResponseUser, string, error)
	CreateUser(ctx context.Context, createUserInput *types.CreateUserInput) error
}

type ServiceImpl struct {
	storage Storage
}

func NewServiceImpl(storage Storage) *ServiceImpl {
	return &ServiceImpl{
		storage: storage,
	}
}

func (s *ServiceImpl) Login(ctx context.Context, email, password string) (*types.ResponseUser, string, error) {
	user, err := s.storage.Get(ctx, email)
	if err != nil {
		return nil, "", err
	}

	// verify password

	if err := isPasswordValid(user.Password, password); err != nil {
		return nil, "", fmt.Errorf("wrong email or password")
	}

	// generate Jwt
	token, err := GenerateJWT(user.UID, user.Email)
	if err != nil {
		return nil, "", err
	}

	return user.ResponseUser(), token, nil
}

func (s *ServiceImpl) CreateUser(ctx context.Context, createUserInput *types.CreateUserInput) error {
	// TODO: check if email exists
	count, err := s.storage.CountEmail(ctx, createUserInput.Email)
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("email already in use")
	}

	// TODO: create new user struct
	newUser := new(types.User)

	uid := s.storage.GenerateID()
	newUser.UID = uid
	newUser.Email = createUserInput.Email
	newUser.FirstName = createUserInput.FirstName
	newUser.LastName = createUserInput.LastName
	newUser.CreatedAt = time.Now().UTC()
	newUser.Busket = types.NewBusket(uid)
	newUser.Password, err = generateHash(createUserInput.Password)
	if err != nil {
		return err
	}

	// TODO: persist user in storage
	err = s.storage.Create(ctx, newUser)
	if err != nil {
		return err
	}

	return nil
}

func generateHash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(b), err
}

func isPasswordValid(foundPassword, givenPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(foundPassword), []byte(givenPassword))

}

type Claims struct {
	UID   string
	Email string
	jwt.RegisteredClaims
}

var (
	secret_key = "kdjhsfoiuhW[OIWEUYR09[2308]]"
)

func GenerateJWT(uid string, email string) (string, error) {
	claims := &Claims{
		UID:   uid,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Hour).UTC()),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret_key))
	if err != nil {
		return "", err
	}

	return token, err
}

func VerifyJWT(signedString string) (string, error) {
	claims := new(Claims)

	token, err := jwt.ParseWithClaims(signedString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret_key), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid authorization token")
	}

	return claims.UID, nil
}
