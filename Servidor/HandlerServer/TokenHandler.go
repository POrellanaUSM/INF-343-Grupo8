package Handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sync"
	"os"
)

var (
	mu    sync.Mutex 
	token string     // Variable  para almacenar token
)

// obtener token al iniciar el servidor
func FetchToken() error {

	 
	clientID := os.Getenv("CLIENT_ID")
	secretID := os.Getenv("SECRET_ID")
	tokenURL := "https://test.api.amadeus.com/v1/security/oauth2/token"

	tokenData := url.Values{}
	tokenData.Set("grant_type", "client_credentials")
	tokenData.Set("client_id", clientID)
	tokenData.Set("client_secret", secretID)

	tokenResponse, err := http.PostForm(tokenURL, tokenData)
	if err != nil {
		return err
	}
	defer tokenResponse.Body.Close()

	var tokenResponseData struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(tokenResponse.Body).Decode(&tokenResponseData); err != nil {
		return err
	}

	mu.Lock()
	token = tokenResponseData.AccessToken
	mu.Unlock()

	

	return nil
}

// devuelve valor del token
func GetToken() string {
	mu.Lock()
	defer mu.Unlock()
	return token
}

