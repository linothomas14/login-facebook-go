package service

import (
	"context"
	"fmt"
	"login-meta-jatis/entity"
	"login-meta-jatis/provider"
	"login-meta-jatis/repository"
)

type LoginService interface {
	Login(ctx context.Context, code string) (token string, err error)
}

type LoginImpl struct {
	tokenRepo repository.TokenRepository
	credRepo  repository.CredentialRepository
	log       provider.ILogger
}

func NewLoginImpl(
	tokenRepo repository.TokenRepository,
	credRepo repository.CredentialRepository,
	log provider.ILogger,
) *LoginImpl {
	return &LoginImpl{
		tokenRepo: tokenRepo,
		credRepo:  credRepo,
		log:       log,
	}
}

func (l *LoginImpl) Login(ctx context.Context, code string) (token string, err error) {

	var newToken entity.Token

	// TODO : TUKER CODE -> ACCESS TOKEN

	// TODO : AMBIL INFO DARI ACCESS TOKEN

	// TODO : GET TOKEN INFO BY HIT ENDPOINT : https://graph.facebook.com/v19.0/debug_token?input_token={YOUR_TOKEN}

	// TODO : GET CLIENT BY USERID

	// TODO : GENERATE TOKEN JATIS,

	// TODO : STORE BESERTA CLIENT_ID DAN ACCESS TOKEN
	err = l.tokenRepo.Create(ctx, newToken)

	if err != nil {
		return token, fmt.Errorf("%v: %v", ErrRepository, err)
	}

	// TODO : SEND TOKEN TO WEBHOOK CLIENT

	return
}
