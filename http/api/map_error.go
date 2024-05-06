package api

import (
	"errors"
	"net/http"

	"login-meta-jatis/model/response"
	"login-meta-jatis/service"
)

func mapError(err error) (int, response.ErrorResponse) {
	var message string
	var code int
	var title string

	if errors.Is(err, service.ErrEmptyRequestBody) {
		message = "Request body cannot be empty"
		title = "InvalidPayload"
		code = http.StatusBadRequest
	} else if errors.Is(err, service.ErrEmptyMessagingProduct) {
		message = "parameter messaging_product is required"
		title = "InvalidPayload"
		code = http.StatusBadRequest
	} else if errors.Is(err, service.ErrEmptyTo) {
		message = "parameter 'to' is required"
		title = "InvalidPayload"
		code = http.StatusBadRequest
	} else if errors.Is(err, service.ErrUnauthorized) {
		message = "invalid token"
		title = "Unauthorized"
		code = http.StatusUnauthorized
	} else if errors.Is(err, service.ErrPrefixNotAllowed) {
		message = "prefix not allowed"
		title = "InvalidPrefix"
		code = http.StatusUnauthorized
	} else {
		title = "InternalServerError"
		code = http.StatusInternalServerError
		message = http.StatusText(code)
	}

	return code, response.ErrorResponse{
		Error: response.Error{
			Code:    code,
			Title:   title,
			Details: message,
		},
	}
}
