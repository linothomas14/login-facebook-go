package main

import (
	"bytes"
	"dummy-login-meta/util"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

var (
	AppID       string
	RedirectURL string
	Secret      string
	ConfigID    string
)

func init() {
	if err := util.LoadConfig("."); err != nil {
		log.Fatal(err)
	}
	AppID = util.Configuration.App.AppID
	RedirectURL = util.Configuration.App.HostURLCallback + "/callback"
	Secret = util.Configuration.App.Secret
	ConfigID = util.Configuration.App.ConfigID
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)
	http.HandleFunc("/validate-token", handleTokenValidity)
	http.HandleFunc("/logout", handleLogout)
	log.Println("Server starting on http://localhost:8080...")
	http.ListenAndServe(":8080", nil)
}
func handleTokenValidity(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Token")
	url := fmt.Sprintf("https://graph.facebook.com/me?access_token=%s", token)
	resp, err := http.Get(url)
	fmt.Println(resp)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	// Baca isi file index.html
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute template with no data (since this is a simple example)
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {

	loginURL := fmt.Sprintf("https://www.facebook.com/v19.0/dialog/oauth?client_id=%s&redirect_uri=%s&scope=email&config_id=%s", AppID, RedirectURL, ConfigID)
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

	payload := map[string]string{"access_token": accessToken}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	webhookURL := util.Configuration.App.HostClientCallback

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, "Failed to send webhook request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		http.Error(w, "Failed to send webhook request", resp.StatusCode)
		return
	}

	fmt.Fprintf(w, "Access Token was sent to %s\nToken : %s", webhookURL, accessToken)
}

func getAccessToken(code string) (string, error) {
	accessTokenURL := fmt.Sprintf("https://graph.facebook.com/v19.0/oauth/access_token?client_id=%s&redirect_uri=%s&client_secret=%s&code=%s", AppID, RedirectURL, Secret, code)

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

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Error != nil {
		return "", fmt.Errorf("facebook API error: %v", result.Error)
	}

	return result.AccessToken, nil
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Token")
	if token == "" {
		http.Error(w, "Token is missing", http.StatusBadRequest)
		return
	}

	if err := revokeFacebookToken(token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "You have been logged out successfully.")
}

func revokeFacebookToken(token string) error {
	revokeURL := fmt.Sprintf("https://graph.facebook.com/me/permissions?access_token=%s", token)
	req, err := http.NewRequest(http.MethodDelete, revokeURL, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("failed to revoke token")
	}

	return nil
}
