package httpapi

import (
	"encoding/json"
	"io"
	"net/http"

	"insurance_hack/internal/db"
	"insurance_hack/internal/dto"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db db.DB
}

type UserCustomClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *AuthService) GetUserToken(writer http.ResponseWriter, r *http.Request) {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		handleError(writer, http.StatusInternalServerError, err)
	}

	var authReq dto.AuthRequest
	if err := json.Unmarshal(reqBody, &authReq); err != nil {
		handleError(writer, http.StatusBadRequest, err)
	}

}
