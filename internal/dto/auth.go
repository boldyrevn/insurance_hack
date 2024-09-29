package dto

import "insurance_hack/internal/model"

type GetTokenRequest struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type GetTokenResponse struct {
	Token string `json:"token,omitempty"`
}

type CreateUserRequest struct {
	model.User
	Password string `json:"password,omitempty"`
}
