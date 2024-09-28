package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type GetUserResponse struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}

// SimpleHandler make dummy response
//
//	@Router			/user/{id}	[get]
//	@Summary		Make dummy response
//	@Description	Send user with specified ID
//	@Produce		json
//	@Param			id	path		string	true	"ID of user"
//	@Success		200	{object}	httpapi.GetUserResponse
//	@Failure		500
func SimpleHandler(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")

	newUser := GetUserResponse{
		Id:   id,
		Name: "Someone",
		Age:  33,
	}

	body, err := json.Marshal(newUser)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}

	_, err = writer.Write(body)
	if err != nil {
		slog.Error("got write error", "error", err)
	}
}
