package service

import (
	"context"
)

type LoginService interface {
	Login(ctx context.Context) (refID string, err error)
}
