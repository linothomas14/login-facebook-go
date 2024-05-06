package api

import (
	"fmt"
	"net/http"

	"login-meta-jatis/provider"
	"login-meta-jatis/service"

	"login-meta-jatis/util"

	"github.com/gin-gonic/gin"
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
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (a *App) sendMessages(ctx *gin.Context) {

}
