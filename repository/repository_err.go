package repository

import "errors"

var (
	ErrNotFound = errors.New("repository: the requested data not found")
	ErrMongo    = errors.New("repository: something went wrong with mongo db")
)
