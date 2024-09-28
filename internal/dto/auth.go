package dto

type AuthRequest struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}
