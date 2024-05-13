package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"login-meta-jatis/provider"
	"login-meta-jatis/util"
	"net/http"
	"strings"
	"text/template"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Token struct {
	ClientID   string    `json:"client_id" bson:"client_id"`
	Session    string    `json:"session_id" bson:"session"`
	TokenJatis string    `json:"token_jatis" bson:"token_jatis"`
	TokenMeta  string    `json:"token_meta" bson:"token_meta"`
	ExpiredAt  time.Time `json:"expired_at" bson:"expired_at"`
}

type AuthToken struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int64  `json:"expires_in"`
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
	RedirectURL = util.Configuration.App.HostURLCallback + "/callback?"
	Secret = util.Configuration.App.Secret
	ConfigID = util.Configuration.App.ConfigID
	fmt.Println(AppID)
	fmt.Println(RedirectURL)
	fmt.Println(Secret)
	fmt.Println(ConfigID)
	fmt.Println(" ")
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", corsMiddleware(handleHome))
	http.HandleFunc("/success-login", corsMiddleware(handleSuccessLogin))
	http.HandleFunc("/get-access-token", corsMiddleware(handleGetToken))

	http.HandleFunc("/login-bento", corsMiddleware(handleLoginBento))
	http.HandleFunc("/callback", corsMiddleware(handleCallbackBento))

	http.HandleFunc("/callback/core", corsMiddleware(handleCallbackBentoCore))

	http.HandleFunc("/logout", corsMiddleware(handleLogout))
	log.Println("Server starting on http://localhost:8080...")
	http.ListenAndServe(":8080", nil)
}

// corsMiddleware adds CORS headers to the response
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any domain
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
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

func handleSuccessLogin(w http.ResponseWriter, r *http.Request) {
	tokenJatis := r.URL.Query().Get("tokenJatis")

	// Create a data structure to pass to the template
	data := struct {
		TokenJatis string
	}{
		TokenJatis: tokenJatis,
	}

	// Parse the HTML template file
	tmpl, err := template.ParseFiles("templates/success_login.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the template with the token data
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Redirect to login page FB
func handleLoginBento(w http.ResponseWriter, r *http.Request) {

	session := r.URL.Query().Get("session")
	// fmt.Println(session)
	clientID := r.URL.Query().Get("client_id")
	// fmt.Println(clientID)
	finalRedirectURL := fmt.Sprintf("%sstate=%s__%s", RedirectURL, clientID, session)

	// fmt.Println("URL = ", finalRedirectURL)

	loginURL := fmt.Sprintf("https://www.facebook.com/dialog/oauth?client_id=%s&display=page&redirect_uri=%s&response_type=token&scope=email,read_insights,pages_manage_cta,pages_manage_instant_articles,pages_show_list,read_page_mailboxes,ads_management,ads_read,business_management,page_events,pages_read_engagement,pages_manage_metadata,pages_read_user_content,pages_manage_ads,pages_manage_posts,pages_manage_engagement,whatsapp_business_messaging,public_profile", AppID, finalRedirectURL)

	http.Redirect(w, r, loginURL, http.StatusSeeOther)
}

func handleCallbackBento(w http.ResponseWriter, r *http.Request) {

	state := r.URL.Query().Get("state")

	callbackUri := fmt.Sprintf("%s/callback/core?state=%s", util.Configuration.App.HostURLCallback, state)
	fmt.Println(callbackUri)
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
	fmt.Println("done ygy")
}

func handleCallbackBentoCore(w http.ResponseWriter, r *http.Request) {

	access_token := r.URL.Query().Get("access_token")
	fmt.Println("access_token = ", access_token)
	state := r.URL.Query().Get("state")

	fmt.Println("state = ", state)

	parts := strings.Split(state, "__")

	// Memastikan ada dua bagian yang dipisahkan
	if len(parts) != 2 {
		fmt.Println("Format string tidak valid")
		return
	}

	// Mendapatkan client ID dan session ID
	clientID := parts[0]
	session := parts[1]

	fmt.Println(clientID)
	fmt.Println(session)

	ctx := context.Background()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://graph.facebook.com/v19.0/oauth/access_token?client_id=%s&client_secret=%s&grant_type=fb_exchange_token&fb_exchange_token=%s", util.Configuration.App.AppID, util.Configuration.App.Secret, access_token), nil)

	if err != nil {
		fmt.Println("Error on create request", err)
		return
	}
	httpClient := &http.Client{}
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		fmt.Println("Error on hit endpoint meta", err)
		return
	}
	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		fmt.Println("Error on read body res", err)
		return
	}

	// fmt.Println("Response Body:", string(body))

	// Decode JSON jika perlu
	var longLivedToken AuthToken
	err = json.Unmarshal(body, &longLivedToken)
	if err != nil {
		fmt.Println("Error on Unmarshal", err)
		return
	}

	tokenJatis, err := generateToken()

	if err != nil {
		fmt.Println("Error on generateToken", err)
		return
	}

	fmt.Println("expires in = ", longLivedToken.ExpiresIn)

	token := Token{
		ClientID:   clientID,
		Session:    session,
		TokenMeta:  longLivedToken.AccessToken,
		TokenJatis: tokenJatis,
		ExpiredAt:  expiresInToDate(longLivedToken.ExpiresIn),
	}

	err = insertToken(token)

	if err != nil {
		fmt.Println("Error on insertToken", err)
		return
	}

	loginURL := fmt.Sprintf("%s/success-login?tokenJatis=%s", util.Configuration.App.HostURLCallback, token.TokenJatis)
	http.Redirect(w, r, loginURL, http.StatusSeeOther)
	fmt.Println("\ndone")
	fmt.Println("Long Lived Token:", longLivedToken)
}

func handleGetToken(w http.ResponseWriter, r *http.Request) {

	session := r.URL.Query().Get("session")
	// fmt.Println(session)
	clientID := r.URL.Query().Get("client_id")
	// fmt.Println(clientID)

	token, err := getTokenByClientIDAndSession(clientID, session)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			response := map[string]interface{}{"error": "Token not found"}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		fmt.Println("Error on getTokenByClientIDAndSession", err)
		return
	}

	response := map[string]interface{}{"token": token.TokenJatis, "expired_at": token.ExpiredAt}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	loginURL := fmt.Sprintf("%s/success-login?tokenJatis=%s", util.Configuration.App.HostURLCallback, token.TokenJatis)
	http.Redirect(w, r, loginURL, http.StatusSeeOther)

}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	session := r.URL.Query().Get("session")
	// fmt.Println(session)
	clientID := r.URL.Query().Get("client_id")
	// fmt.Println(clientID)

	token, err := getTokenByClientIDAndSession(clientID, session)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			response := map[string]interface{}{"error": "Token not found"}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		fmt.Println("Error on getTokenByClientIDAndSession", err)
		return
	}

	if err := revokeFacebookToken(token.TokenMeta); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = deleteTokenByClientIDAndSession(clientID, session)

	if err != nil {
		fmt.Println("Error on getTokenByClientIDAndSession", err)
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

func insertToken(token Token) error {
	ctx := context.Background()
	mongoclient, err := provider.NewMongoDBClient()

	if err != nil {
		fmt.Printf("Connect to DB failed: %s", err)
		return err
	}

	db := mongoclient.Database(util.Configuration.MongoDB.Database)

	coll := db.Collection(util.Configuration.MongoDB.Collection.Token)
	result, err := coll.InsertOne(ctx, token)
	if err != nil {
		fmt.Printf("creating chat history in MongoDB failed: %s", err)
		return err
	}

	_, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		fmt.Printf("error asserting InsertedID to ObjectID")
		return err
	}
	return nil
}

func generateToken() (string, error) {
	// Menghitung jumlah byte yang diperlukan untuk panjang token yang diinginkan
	length := 35

	byteLength := length / 4 * 3
	if length%4 != 0 {
		byteLength++
	}

	// Membuat slice untuk menyimpan byte acak
	randomBytes := make([]byte, byteLength)

	// Mengisi slice dengan byte acak
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Mengonversi byte menjadi string base64
	token := base64.URLEncoding.EncodeToString(randomBytes)

	// Memastikan panjang token sesuai dengan yang diinginkan
	if len(token) > length {
		token = token[:length]
	}

	return token, nil
}

func expiresInToDate(expiresIn int64) time.Time {
	// Membuat durasi dari expiresIn dalam detik
	expiresInDuration := time.Second * time.Duration(expiresIn)

	// Menghitung waktu sekarang
	currentTime := time.Now()

	// Menambahkan expiresInDuration ke currentTime untuk mendapatkan waktu kadaluarsa
	expirationTime := currentTime.Add(expiresInDuration)

	return expirationTime
}

func getTokenByClientIDAndSession(clientID string, session string) (token Token, err error) {

	ctx := context.Background()
	mongoclient, err := provider.NewMongoDBClient()

	if err != nil {
		fmt.Printf("Connect to DB failed: %s", err)
		return token, err
	}

	db := mongoclient.Database(util.Configuration.MongoDB.Database)

	coll := db.Collection(util.Configuration.MongoDB.Collection.Token)
	filter := bson.M{"client_id": clientID, "session": session}

	// Melakukan pencarian dalam koleksi
	err = coll.FindOne(ctx, filter).Decode(&token)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf("Error Token not found: %s", err)
			return token, err
		}
		fmt.Printf("Error on FindOne Data : %s", err)
		return token, err
	}

	return token, nil

}

func deleteTokenByClientIDAndSession(clientID string, session string) error {
	ctx := context.Background()
	mongoclient, err := provider.NewMongoDBClient()
	if err != nil {
		fmt.Printf("Connect to DB failed: %s", err)
		return err
	}

	db := mongoclient.Database(util.Configuration.MongoDB.Database)
	coll := db.Collection(util.Configuration.MongoDB.Collection.Token)
	filter := bson.M{"client_id": clientID, "session": session}

	// Menghapus dokumen yang sesuai dengan filter
	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Printf("Error deleting token: %s", err)
		return err
	}

	// Jika tidak ada dokumen yang dihapus, kembalikan error
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
