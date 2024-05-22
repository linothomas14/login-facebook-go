package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"login-meta-jatis/entity"
	"login-meta-jatis/model/response"
	"login-meta-jatis/provider"
	"login-meta-jatis/repository"
	"login-meta-jatis/util"
	"net/http"
	"time"
)

type LoginService interface {
	LoginCore(ctx context.Context, token string, clientID string, session string) (err error)
	FindClientByClientID(ctx context.Context, clientID string) (err error)
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

func (l *LoginImpl) LoginCore(ctx context.Context, token string, clientID string, session string) (err error) {

	tokenMeta, err := extendToken(token)

	if err != nil {
		return
	}

	jatisToken := util.GenerateHexUUID()

	newToken := entity.Token{
		ClientID:  clientID,
		Session:   session,
		TokenMeta: tokenMeta.AccessToken,
		Token:     jatisToken,
		IsActive:  true,
		ExpiredAt: expiresInToDate(tokenMeta.ExpiresIn),
		CreatedAt: time.Now().Add(7 * time.Hour),
	}

	err = l.tokenRepo.Save(ctx, newToken)

	if err != nil {
		return fmt.Errorf("%v: %v", ErrRepository, err)
	}

	return

}

func (l *LoginImpl) FindClientByClientID(ctx context.Context, clientID string) (err error) {
	return l.credRepo.FindClientByClientID(ctx, clientID)
}

func extendToken(shortLivedToken string) (longLivedToken response.TokenMeta, err error) {

	ctx := context.Background()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://graph.facebook.com/v19.0/oauth/access_token?client_id=%s&client_secret=%s&grant_type=fb_exchange_token&fb_exchange_token=%s", util.Configuration.App.AppID, util.Configuration.App.Secret, shortLivedToken), nil)

	if err != nil {
		// Handle error
		return
	}
	httpClient := &http.Client{}
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		// Handle error
		return
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		// Handle error
		return
	}

	err = json.Unmarshal(body, &longLivedToken)
	if err != nil {

		return
	}
	return longLivedToken, err
}

func expiresInToDate(expiresIn int64) time.Time {

	expiresInDuration := time.Second * time.Duration(expiresIn)

	currentTime := time.Now()

	expirationTime := currentTime.Add(expiresInDuration)

	return expirationTime
}
