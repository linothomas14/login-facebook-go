package service

import "errors"

var (
	ErrRepository            = errors.New("service: repository error happened")
	ErrEmptyRequestBody      = errors.New("service: request body is empty")
	ErrEmptyMessagingProduct = errors.New("service: \"messaging_product\" is empty")
	ErrEmptyTo               = errors.New("service: \"to\" is empty")
	ErrUnauthorized          = errors.New("service: invalid token")
	ErrPrefixNotAllowed      = errors.New("service: prefix not allowed")
)
