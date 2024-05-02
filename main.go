package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// TokenInfo stores information from Facebook's token debug endpoint

const (
	appID       = "1121571312516703"
	redirectURL = "https://b753-101-128-100-252.ngrok-free.app/callback"
	secret      = "23c5f0ef44a3eb3f2fcc9c21948f928b"
	configID    = "427409016599790"
)

func main() {
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
	fmt.Println("   ")
	fmt.Println(resp)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
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
	accessTokenURL := fmt.Sprintf("https://graph.facebook.com/v19.0/oauth/access_token?client_id=%s&redirect_uri=%s&client_secret=%s&code=%s", appID, redirectURL, secret, code)

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

func handleLogout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Token")
	if token == "" {
		http.Error(w, "Token is missing", http.StatusBadRequest)
		return
	}

	if err := invalidateFacebookToken(token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "You have been logged out successfully.")
}

func invalidateFacebookToken(token string) error {
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
		return fmt.Errorf("Failed to revoke token")
	}

	return nil
}
