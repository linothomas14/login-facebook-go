package api

import (
	"fmt"
	"net/http"

	"login-meta-jatis/provider"
	"login-meta-jatis/service"

	"login-meta-jatis/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type App struct {
	loginService service.LoginService
	log          provider.ILogger
}

func NewApp(srv service.LoginService, log provider.ILogger) *App {
	return &App{loginService: srv, log: log}
}

func (a *App) CreateServer(address string) (*http.Server, error) {
	gin.SetMode(util.Configuration.Server.Mode)

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(loggingMiddleware(a.log))

	r.GET("/ping", a.checkConnectivity)

	r.POST("/login", a.handleLogin)
	r.POST("/callback", a.handleCallback)

	server := &http.Server{
		Addr:    address,
		Handler: r,
	}

	return server, nil
}

func (a *App) checkConnectivity(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (a *App) handleLogin(ctx *gin.Context) {
	loginURL := fmt.Sprintf("https://www.facebook.com/v19.0/dialog/oauth?client_id=%s&redirect_uri=%s&scope=email&config_id=%s", util.Configuration.App.AppID, util.Configuration.App.HostURLCallback, util.Configuration.App.ConfigID)
	ctx.Redirect(http.StatusSeeOther, loginURL)
}

func (a *App) handleCallback(ctx *gin.Context) {

	var reqID string
	val, ok := ctx.Get("req-id")
	if ok {
		reqID = val.(string)
	}

	code := ctx.Request.URL.Query().Get("code")

	if code == "" {
		ctx.String(http.StatusBadRequest, "Code not found")
		return
	}
	token, err := a.loginService.Login(ctx, code)
	if err != nil {
		a.log.WithFields(provider.AppLog, logrus.Fields{"REQUEST_ID": reqID}).Errorf("send message error, reason: %s", err)
		code, errorResp := mapError(err)

		ctx.JSON(code, errorResp)
		ctx.Abort()
		return
	}
	a.log.WithFields(provider.AppLog, logrus.Fields{"REQUEST_ID": reqID}).Infof("token generated : %#v", token)

	// Redirect to web client with the token
	loginURL := fmt.Sprintf("https://8d00-180-252-93-189.ngrok-free.app/?token=%s", token)
	ctx.Redirect(http.StatusSeeOther, loginURL)
}
