package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/luqus/authservice/types"
)

type APIFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) *types.Error

type APIServer struct {
	listenAddress string
	router        *mux.Router
	service       Service
}

func NewAPIServer(listenAddress string, service Service) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
		router:        mux.NewRouter(),
		service:       service,
	}
}

func (api *APIServer) Run() error {

	return http.ListenAndServe(api.listenAddress, api.router)
}

func handler(fn APIFunc) http.HandlerFunc {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request-id", uuid.New())

	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(ctx, w, r)
		if err != nil {
			w.WriteHeader(err.StatusCode)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		}
	}
}

func (api *APIServer) login(ctx context.Context, w http.ResponseWriter, r *http.Request) *types.Error {
	loginInput := new(types.LoginInput)

	err := requestDecoder(r.Body, loginInput)
	if err != nil {
		return &types.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid login data",
		}
	}

	responseUser, err := api.service.Login(ctx, loginInput.Email, loginInput.Password)
	if err != err {
		return &types.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "wrong email or password",
		}
	}

	err = writeResponse(w, responseUser)
	if err != nil {
		return &types.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "encoding error",
		}
	}

	return nil
}

func (api *APIServer) createUser(ctx context.Context, w http.ResponseWriter, r *http.Request) *types.Error {
	createUserInput := new(types.CreateUserInput)
	err := requestDecoder(r.Body, createUserInput)
	if err != nil {
		return &types.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid create data format",
		}
	}

	err = api.service.CreateUser(ctx, createUserInput)
	if err != nil {
		return &types.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	response := map[string]string{"message": "succefully signed up"}
	err = writeResponse(w, response)
	if err != nil {
		return &types.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "encoding response error",
		}
	}

	return nil
}

func writeResponse(w http.ResponseWriter, v any) error {
	return json.NewEncoder(w).Encode(v)
}

func requestDecoder(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)

}
