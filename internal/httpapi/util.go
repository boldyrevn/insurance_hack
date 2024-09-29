package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type ErrorResponse struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

func handleError(writer http.ResponseWriter, status int, err error) {
	respBody, _ := json.Marshal(ErrorResponse{
		Status:  strconv.Itoa(status),
		Message: err.Error(),
	})
	writer.WriteHeader(status)
	_, _ = writer.Write(respBody)
}

func sendResponse(writer http.ResponseWriter, respStruct any) {
	respBody, err := json.Marshal(respStruct)
	if err != nil {
		slog.Error("failed to marshal response struct", "respStruct", respStruct)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = writer.Write(respBody); err != nil {
		slog.Error("failed to send response body", "respBody", string(respBody))
	}
}
