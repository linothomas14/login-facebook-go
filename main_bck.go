package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"login-meta-jatis/util"
	"net/http"
	"text/template"
)

type AuthToken struct {
	AccessToken string `json:"access_token,omitempty"`
}

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
	RedirectURL = util.Configuration.App.HostURLCallback + "%2Fcallback%3Fstate%3DcobaState123"
	Secret = util.Configuration.App.Secret
	ConfigID = util.Configuration.App.ConfigID
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/login-sdk", handleLoginSDK)
	http.HandleFunc("/login-bento", handleLoginBento)
	http.HandleFunc("/callback", handleCallbackBento)
	http.HandleFunc("/callback/core", handleCallbackBentoCore)
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

func handleLoginSDK(w http.ResponseWriter, r *http.Request) {
	// Baca isi file index.html
	tmpl, err := template.ParseFiles("templates/index_facebook_SDK.html")
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

	loginURL := fmt.Sprintf("https://www.facebook.com/dialog/oauth?client_id=%s&redirect_uri=%s&state=coba&config_id=%s", AppID, RedirectURL, ConfigID)
	http.Redirect(w, r, loginURL, http.StatusSeeOther)
}

func handleLogin1(w http.ResponseWriter, r *http.Request) {

	loginURL := fmt.Sprintf("https://www.facebook.com/dialog/oauth?client_id=%s&redirect_uri=%s&config_id=%s", AppID, "https://d28f-180-252-88-237.ngrok-free.app/callback-client1?", ConfigID)
	http.Redirect(w, r, loginURL, http.StatusSeeOther)
}

func handleLogin2(w http.ResponseWriter, r *http.Request) {

	loginURL := fmt.Sprintf("https://www.facebook.com/dialog/oauth?client_id=%s&redirect_uri=%s&config_id=%s", AppID, "https://d28f-180-252-88-237.ngrok-free.app/callback-client2", ConfigID)
	http.Redirect(w, r, loginURL, http.StatusSeeOther)
}

func handleLoginBento(w http.ResponseWriter, r *http.Request) {

	loginURL := fmt.Sprintf("https://www.facebook.com/dialog/oauth?client_id=%s&display=page&redirect_uri=%s&response_type=token&scope=pages_read_engagement,pages_manage_metadata,instagram_basic,instagram_manage_messages,public_profile", AppID, RedirectURL)

	http.Redirect(w, r, loginURL, http.StatusSeeOther)
}

func handleCallbackBento(w http.ResponseWriter, r *http.Request) {

	state := r.URL.Query().Get("state")

	callbackUri := fmt.Sprintf("%s/callback/core?state=%s", util.Configuration.App.HostURLCallback, state)
	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<body>
		<script type="text/javascript">
			const part = window.location.href.split("#")
			window.location.href = "%s&" + part[1]
		</script>
		</body>
		</html>
	`, callbackUri)

	w.Header().Set("Content-Type", "text/html")

	// Mengirimkan respons HTML ke client
	_, err := fmt.Fprintf(w, "%s", html)
	if err != nil {
		// Tangani kesalahan jika terjadi saat mengirim respons
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func handleCallbackBentoCore(w http.ResponseWriter, r *http.Request) {

	access_token := r.URL.Query().Get("access_token")
	fmt.Println("access_token = ", access_token)
	state := r.URL.Query().Get("state")
	fmt.Println("state = ", state)

	ctx := context.Background()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://graph.facebook.com/v19.0/oauth/access_token?client_id=%s&client_secret=%s&grant_type=fb_exchange_token&fb_exchange_token=%s", util.Configuration.App.AppID, util.Configuration.App.Secret, access_token), nil)

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

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		// Handle error
		return
	}

	fmt.Println("Response Body:", string(body))

	// Decode JSON jika perlu
	var longLivedToken AuthToken
	err = json.Unmarshal(body, &longLivedToken)
	if err != nil {

		return
	}
	payload := map[string]string{"access_token": longLivedToken.AccessToken}
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

	// fmt.Fprintf(w, "Access Token was sent to %s\nToken : %s", webhookURL, accessToken)
	// REDIRECT TO WEB CLIENT DENGAN MEMBAWA TOKEN JATIS
	loginURL := fmt.Sprintf("https://d28f-180-252-88-237.ngrok-free.app/")
	http.Redirect(w, r, loginURL, http.StatusSeeOther)
	fmt.Println("done")
	fmt.Println("Long Lived Token:", longLivedToken)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {

	code := r.URL.Query().Get("code")
	fmt.Println("code=", code)

	client_id := r.URL.Query().Get("client_id")
	fmt.Println("client_id=", client_id)

	state := r.URL.Query().Get("state")
	fmt.Println("state=", state)

	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}
	fmt.Println("CODE = ", code)
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

	// fmt.Fprintf(w, "Access Token was sent to %s\nToken : %s", webhookURL, accessToken)
	// REDIRECT TO WEB CLIENT DENGAN MEMBAWA TOKEN JATIS
	loginURL := fmt.Sprintf("https://d28f-180-252-88-237.ngrok-free.app/")
	http.Redirect(w, r, loginURL, http.StatusSeeOther)
	fmt.Println("done")
}

func getAccessToken(code string) (string, error) {
	accessTokenURL := fmt.Sprintf("https://graph.facebook.com/v19.0/oauth/access_token?client_id=%s&redirect_uri=%s&client_secret=%s&code=%s", AppID, RedirectURL, Secret, code)

	fmt.Println("access token URL = ", accessTokenURL)

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
