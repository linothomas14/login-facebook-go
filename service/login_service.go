package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"login-meta-jatis/entity"
	"login-meta-jatis/provider"
	"login-meta-jatis/repository"
	"login-meta-jatis/util"
	"net/http"
)

type LoginService interface {
	LoginCore(ctx context.Context, token string, state string) (err error)
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

func (l *LoginImpl) LoginCore(ctx context.Context, token string, state string) (err error) {

	var newToken entity.Token

	// TODO : EXTEND SHORT LIVE TOKEN -> LONG LIVE TOKEN
	longLivedToken, err := extendToken(token)

	if err != nil {
		return
	}

	newToken.TokenMeta = longLivedToken

	// TODO : GENERATE TOKEN JATIS,

	// TODO : STORE TOKEN_JATIS, STATE(CLIENT_ID) DAN ACCESS_TOKEN
	err = l.tokenRepo.Create(ctx, newToken)

	if err != nil {
		return fmt.Errorf("%v: %v", ErrRepository, err)
	}

	// TODO : SEND TOKEN TO WEBHOOK CLIENT

	return
}

func extendToken(shortLivedToken string) (string, error) {

	type AuthToken struct {
		AccessToken string `json:"access_token,omitempty"`
	}

	ctx := context.Background()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://graph.facebook.com/v19.0/oauth/access_token?client_id=%s&client_secret=%s&grant_type=fb_exchange_token&fb_exchange_token=%s", util.Configuration.App.AppID, util.Configuration.App.Secret, shortLivedToken), nil)

	if err != nil {
		// Handle error
		return "", err
	}
	httpClient := &http.Client{}
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		// Handle error
		return "", err
	}
	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		// Handle error
		return "", err
	}

	var longLivedToken AuthToken
	err = json.Unmarshal(body, &longLivedToken)
	if err != nil {

		return "", err
	}
	return longLivedToken.AccessToken, err
}
