package identity

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	v1 "github.com/ShapleyIO/iam/api/v1"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type ServiceIdentity struct {
	ctx         context.Context
	redisClient *redis.Client
}

func NewServiceIdentity(redisClient *redis.Client) *ServiceIdentity {
	return &ServiceIdentity{
		ctx:         context.Background(),
		redisClient: redisClient,
	}
}

// Create a User
// (POST /v1/user)
func (s *ServiceIdentity) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Implementation goes here
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Unmarshal the request body
	var user v1.User
	if err := json.Unmarshal(body, &user); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the user already exists
	if s.redisClient.Exists(s.ctx, string(user.Email)).Val() == 1 {
		log.Error().Str("email", string(user.Email)).Msg("user already exists")
		w.WriteHeader(http.StatusConflict)
	}

	// Create the user
	if err := s.redisClient.Set(s.ctx, string(user.Email), user, 0).Err(); err != nil {
		log.Error().Err(err).Msg("failed to create user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Update a User's Password
// (PUT /v1/user/password/{user_id})
func (s *ServiceIdentity) UpdateUserPassword(w http.ResponseWriter, r *http.Request, params v1.UpdateUserPasswordParams) {
	// Implementation goes here
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Unmarshal the request body
	var password v1.Password
	if err := json.Unmarshal(body, &password); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the user
	userJson, err := s.redisClient.Get(s.ctx, string(params.Email)).Result()
	if err != nil {
		log.Error().Err(err).Msg("failed to get user")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Unmarshal the user
	var user UserWithPassword
	if err := json.Unmarshal([]byte(userJson), &user); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update the user's password
	user.Password = password.Password
	if err := s.redisClient.Set(s.ctx, string(params.Email), user, 0).Err(); err != nil {
		log.Error().Err(err).Msg("failed to update user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Delete a User
// (DELETE /v1/user/{user_id})
func (s *ServiceIdentity) DeleteUser(w http.ResponseWriter, r *http.Request, params v1.DeleteUserParams) {
	// Delete user
	if err := s.redisClient.Del(s.ctx, string(params.Email)).Err(); err != nil {
		log.Error().Err(err).Msg("failed to delete user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Get a User
// (GET /v1/user/{user_id})
func (s *ServiceIdentity) GetUser(w http.ResponseWriter, r *http.Request, params v1.GetUserParams) {
	// Get user
	userJson, err := s.redisClient.Get(s.ctx, string(params.Email)).Result()
	if err != nil {
		log.Error().Err(err).Msg("failed to get user")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(userJson))
}

// Update a User
// (PUT /v1/user/{user_id})
func (s *ServiceIdentity) UpdateUser(w http.ResponseWriter, r *http.Request, params v1.UpdateUserParams) {
	// Update User
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Unmarshal the request body
	var user v1.User
	if err := json.Unmarshal(body, &user); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the user
	userJson, err := s.redisClient.Get(s.ctx, string(params.Email)).Result()
	if err != nil {
		log.Error().Err(err).Msg("failed to get user")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Unmarshal the user
	var userWithPassword UserWithPassword
	if err := json.Unmarshal([]byte(userJson), &userWithPassword); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update the user
	userWithPassword.FirstName = user.FirstName
	userWithPassword.LastName = user.LastName
	userWithPassword.Email = user.Email
	if err := s.redisClient.Set(s.ctx, string(params.Email), userWithPassword, 0).Err(); err != nil {
		log.Error().Err(err).Msg("failed to update user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
