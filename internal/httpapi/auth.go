package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"insurance_hack/internal/db"
	"insurance_hack/internal/dto"
	"insurance_hack/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey string

func init() {
	secretKey = os.Getenv("SECRET_KEY")
}

type AuthHandler struct {
	DB db.DB
}

type UserCustomClaims struct {
	jwt.RegisteredClaims
	Login string `json:"login"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateUser creates new user
//
//	@Router		/user/create	[post]
//	@Summary	Create new user
//	@Accept		json
//	@Produce	json
//	@Param		user	body		dto.CreateUserRequest	true	"New user"
//	@Success	200		{object}	model.User
//	@Failure	400		{object}	httpapi.ErrorResponse
func (s *AuthHandler) CreateUser(writer http.ResponseWriter, r *http.Request) {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		handleError(writer, http.StatusInternalServerError, err)
		return
	}

	var req dto.CreateUserRequest
	if err := json.Unmarshal(reqBody, &req); err != nil {
		handleError(writer, http.StatusBadRequest, err)
		return
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		handleError(writer, http.StatusInternalServerError, err)
		return
	}

	err = s.DB.CreateUser(r.Context(), req.User, hashedPassword)
	if err != nil {
		handleError(writer, http.StatusBadRequest, err)
		return
	}

	sendResponse(writer, req.User)
}

func getLoginFromToken(tokenString string) (string, error) {
	var claims UserCustomClaims
	if _, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	}); err != nil {
		return "", err
	}
	return claims.Login, nil
}

func (s *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		arr := strings.Split(authHeader, " ")
		if len(arr) != 2 || arr[0] != "Bearer" {
			handleError(writer, http.StatusUnauthorized, fmt.Errorf("wrong auth header format"))
			return
		}

		login, err := getLoginFromToken(arr[1])
		if err != nil {
			handleError(writer, http.StatusUnauthorized, fmt.Errorf("failed to get login from token: %w", err))
			return
		}

		fmt.Println(login)

		request = request.WithContext(context.WithValue(request.Context(), model.LoginCtxKey, login))
		next.ServeHTTP(writer, request)
	})
}

func createToken(login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserCustomClaims{
		Login: login,
	})
	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return ss, nil
}

// GetUserToken get user's token
//
//	@Router		/user/token	[post]
//	@Summary	Get user's token
//	@Accept		json
//	@Produce	json
//	@Param		user	body		dto.GetTokenRequest	true	"User's login and password"
//	@Success	200		{object}	dto.GetTokenResponse
//	@Failure	400		{object}	httpapi.ErrorResponse
//	@Failure	401		{object}	httpapi.ErrorResponse
func (s *AuthHandler) GetUserToken(writer http.ResponseWriter, r *http.Request) {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		handleError(writer, http.StatusInternalServerError, err)
		return
	}

	var authReq dto.GetTokenRequest
	if err := json.Unmarshal(reqBody, &authReq); err != nil {
		handleError(writer, http.StatusBadRequest, err)
		return
	}

	password, err := s.DB.GetHashedPassword(r.Context(), authReq.Login)
	if err != nil {
		handleError(writer, http.StatusBadRequest, err)
		return
	}

	var token string
	if CheckPasswordHash(authReq.Password, password) {
		if token, err = createToken(authReq.Login); err != nil {
			handleError(writer, http.StatusInternalServerError, err)
		} else {
			sendResponse(writer, dto.GetTokenResponse{Token: token})
		}
		return
	}

	handleError(writer, http.StatusUnauthorized, fmt.Errorf("wrong user password"))
}

// GetUser get user from DB
//
//	@Router		/user	[get]
//	@Summary	Get user
//	@Produce	json
//	@Success	200	{object}	model.User
//	@Failure	500	{object}	httpapi.ErrorResponse
//	@Failure	401	{object}	httpapi.ErrorResponse
//	@Security	ApiKeyAuth
func (s *AuthHandler) GetUser(writer http.ResponseWriter, request *http.Request) {
	var (
		login string
		ok    bool
	)
	if login, ok = request.Context().Value(model.LoginCtxKey).(string); !ok || login == "" {
		handleError(writer, http.StatusInternalServerError, fmt.Errorf("no login in context"))
		return
	}

	user, err := s.DB.GetUserByLogin(request.Context(), login)
	if err != nil {
		handleError(writer, http.StatusInternalServerError, err)
		return
	}

	request = request.WithContext(context.WithValue(request.Context(), model.UserCtxKey, user))
	sendResponse(writer, user)
}
