package api

import (
	"fmt"
	"net/http"
	"strings"

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

	r.LoadHTMLGlob("templates/*")

	r.Static("/static", "./static")

	r.GET("/", a.handleHTML)

	r.POST("/ping", a.checkConnectivity)

	r.GET("/login", a.Login)
	r.GET("/callback", a.handleCallback)
	r.GET("/callback/core", a.handleCallbackCore)

	server := &http.Server{
		Addr:    address,
		Handler: r,
	}

	return server, nil
}

func (a *App) handleHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func (a *App) checkConnectivity(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (a *App) Login(ctx *gin.Context) {

	session := ctx.Query("session")
	clientID := ctx.Query("client_id")

	if session == "" || clientID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "client_id and session are required"})
		return
	}

	err := a.loginService.FindClientByClientID(ctx, clientID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	redirectURL := fmt.Sprintf("%s/callback?state=%s__%s", util.Configuration.App.HostURLCallback, clientID, session)

	loginURL := fmt.Sprintf("https://www.facebook.com/dialog/oauth?client_id=%s&display=page&redirect_uri=%s&response_type=token&scope=email,read_insights,pages_manage_cta,pages_manage_instant_articles,pages_show_list,read_page_mailboxes,ads_management,ads_read,business_management,page_events,pages_read_engagement,pages_manage_metadata,pages_read_user_content,pages_manage_ads,pages_manage_posts,pages_manage_engagement,whatsapp_business_messaging,public_profile", util.Configuration.App.AppID, redirectURL)

	ctx.Redirect(http.StatusSeeOther, loginURL)
}

func (a *App) handleCallback(ctx *gin.Context) {

	state := ctx.Query("state")
	callbackURI := fmt.Sprintf("%s/callback/core?state=%s", util.Configuration.App.HostURLCallback, state)

	ctx.Header("Content-Type", "text/html")
	ctx.String(http.StatusOK, fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<body>
		<script type="text/javascript">
			const parts = window.location.href.split("#");
			window.location.href = "%s&" + parts[1];
		</script>
		</body>
		</html>
	`, callbackURI))
}

func (a *App) handleCallbackCore(ctx *gin.Context) {

	access_token := ctx.Query("access_token")

	state := ctx.Query("state")

	parts := strings.Split(state, "__")

	// Memastikan ada dua bagian yang dipisahkan
	if len(parts) != 2 {
		fmt.Println("Format string tidak valid")
		return
	}

	// Mendapatkan client ID dan session ID
	clientID := parts[0]
	session := parts[1]

	err := a.loginService.LoginCore(ctx, access_token, clientID, session)
	if err != nil {
		a.log.WithFields(provider.AppLog, logrus.Fields{"SESSION": session}).Errorf("send message error, reason: %s", err)
		// 	code, errorResp := mapError(err)

		// 	ctx.JSON(code, errorResp)
		ctx.Abort()
		return
	}

	ctx.HTML(http.StatusOK, "success_login.html", nil)

}
