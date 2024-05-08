package main

// import (
// 	"bytes"
// 	"dummy-login-meta/util"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"text/template"
// )

// var (
// 	AppID       string
// 	RedirectURL string
// 	Secret      string
// 	ConfigID    string
// )

// func init() {
// 	if err := util.LoadConfig("."); err != nil {
// 		log.Fatal(err)
// 	}
// 	AppID = util.Configuration.App.AppID
// 	RedirectURL = util.Configuration.App.HostURLCallback + "/callback"
// 	Secret = util.Configuration.App.Secret
// 	ConfigID = util.Configuration.App.ConfigID
// }

// func main() {
// 	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
// 	http.HandleFunc("/", handleHome)
// 	http.HandleFunc("/login", handleLogin)
// 	http.HandleFunc("/login-sdk", handleLoginSDK) // 
// 	http.HandleFunc("/callback", handleCallback)
// 	http.HandleFunc("/validate-token", handleTokenValidity)
// 	http.HandleFunc("/logout", handleLogout)
// 	log.Println("Server starting on http://localhost:8080...")
// 	http.ListenAndServe(":8080", nil)
// }
// func handleTokenValidity(w http.ResponseWriter, r *http.Request) {
// 	token := r.Header.Get("Token")
// 	url := fmt.Sprintf("https://graph.facebook.com/me?access_token=%s", token)
// 	resp, err := http.Get(url)
// 	fmt.Println(resp)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer resp.Body.Close()
// }

// func handleLoginSDK(w http.ResponseWriter, r *http.Request) {
// 	// Baca isi file index.html
// 	tmpl, err := template.ParseFiles("templates/index_facebook_SDK.html")
// 	if err != nil {
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}

// 	logger.Infof(provider.AppLog, "Successfully connected to MongoDB.")

// 	go func(c *mongo.Client, logger provider.ILogger) {
// 		var credRepo repository.CredentialRepository = repository.NewCredRepositoryImpl(c, logger)
// 		var tokenRepo repository.TokenRepository = repository.NewTokenRepositoryImpl(c, logger)
// 		var loginService service.LoginService = service.NewLoginImpl(tokenRepo, credRepo, logger)
// 		app := api.NewApp(loginService, logger)
// 		addr := fmt.Sprintf(":%v", util.Configuration.Server.Port)
// 		server, err := app.CreateServer(addr)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		logger.Infof(provider.AppLog, "Server running at: %s", addr)
// 		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			logger.Errorf(provider.AppLog, "Server error: %v", err)
// 		}

// 	}(mongoClient, logger)

// 	shutdownCh := make(chan os.Signal, 1)
// 	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

// 	sig := <-shutdownCh
// 	logger.Infof(provider.AppLog, "Receiving signal: %s", sig)

// 	func(c *mongo.Client) {
// 		if err := c.Disconnect(context.Background()); err != nil {
// 			log.Fatal(err)
// 		}

// 		logger.Infof(provider.AppLog, "Successfully disconnected from MongoDB.")

// 	}(mongoClient)
// }
