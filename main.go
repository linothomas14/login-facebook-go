package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	appID = "242067665649305"
	// redirectURL = "https://webhook.site/dbe0257a-e729-4fc6-aef6-b3372e3472ab"
	redirectURL = "https://8fd1-122-50-6-195.ngrok-free.app/callback"
	secret      = "fb20deb1c5cb8e17d41b40c49f9f33c9"
	configID    = "2062786767487557"
)

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)
	http.ListenAndServe(":8080", nil)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Login with Facebook</h1><a href=\"/login\">Login with Facebook</a>")
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	loginURL := fmt.Sprintf("https://www.facebook.com/v19.0/dialog/oauth?client_id=%s&redirect_uri=%s&scope=email&config_id=%s", appID, redirectURL, configID)
	http.Redirect(w, r, loginURL, http.StatusSeeOther)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {

	code := r.URL.Query().Get("code")
	fmt.Println(code)
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	accessToken, err := getAccessToken(code)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Access Token: %s", accessToken)
}

func getAccessToken(code string) (string, error) {
	accessTokenURL := fmt.Sprintf("http://localhost:5001/v19.0/oauth/access_token?client_id=%s&redirect_uri=%s&client_secret=%s&code=%s", appID, redirectURL, secret, code)

	resp, err := http.Get(accessTokenURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		Error       *struct {
			Message string `json:"message"`
			Type    string `json:"type"`
			Code    int    `json:"code"`
		} `json:"error"`
	}

	fmt.Println(resp.Body)
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Error != nil {
		return "", fmt.Errorf("Facebook API error: %s", result.Error)
	}

	return result.AccessToken, nil
}
