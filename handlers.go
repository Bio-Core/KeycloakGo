package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
)

const emptyString = ""

var oauthStateString string //randomly generated state string
var token *oauth2.Token     //token for keycloak
var login bool

//HandleLogin is the keycloak login funtion
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	token = nil
	//create a random string for oath2 verification
	oauthStateString = randSeq(20)
	//Uses random gnerated string to verify keyclock security
	url := oauth2Config.AuthCodeURL(oauthStateString)
	//redirects to loginCallback
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

//HandleLoginCallback is a fuction that verifies login success and forwards to index
func HandleLoginCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	//Checks that the strings are in a consistent state
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	//Gets the code from keycloak
	code := r.FormValue("code")
	//Exchanges code for token
	token, err = oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("Code exchange failed with '%v'\n", err)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	client := &http.Client{}
	url := keycloakserver + "/auth/realms/" + realm + "/protocol/openid-connect/userinfo"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	//Sends the token to get user info
	response, err := client.Do(req)

	//Checks if token and authentication were successful
	if err != nil || response.Status != "200 OK" {
		//forwards back to login if not successful
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	} else {
		body, _ := ioutil.ReadAll(response.Body)
		var f interface{}
		json.Unmarshal(body, &f)
		m := f.(map[string]interface{})
		username := m["preferred_username"].(string)
		//forwards to index if login sucessful
		logAction(username, actionLogin, emptyString)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	return
}

//AuthMiddleware is a middlefuntion that verifies authentication before each redirect
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//If running unit tests skip authentication (temp)
		client := &http.Client{}
		url := keycloakserver + "/auth/realms/" + realm + "/protocol/openid-connect/userinfo"
		req, _ := http.NewRequest("GET", url, nil)
		if token == nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		req.Header.Set("Authorization", "Bearer "+token.AccessToken)
		//Check if token is still valid
		response, err := client.Do(req)
		if err != nil || response.Status != "200 OK" {
			//Go to login if token is no longer valid
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		} else {
			body, _ := ioutil.ReadAll(response.Body)
			var f interface{}
			json.Unmarshal(body, &f)
			m := f.(map[string]interface{})
			username := m["preferred_username"].(string)
			//Go to redirect if token is still valid
			logAction(username, actionPageAccess, r.RequestURI)
			next.ServeHTTP(w, r)
		}
	})
	//return function for page handling
	return handler
}

//AuthMiddlewareHandler is a middlefuntion that verifies authentication before each redirect
func AuthMiddlewareHandler(next http.Handler) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//If running unit tests skip authentication (temp)
		client := &http.Client{}
		url := keycloakserver + "/auth/realms/" + realm + "/protocol/openid-connect/userinfo"
		req, _ := http.NewRequest("GET", url, nil)
		if token == nil {
			if oauthStateString == "" {
				HandleLogin(w, r)
				return
			}
			result := getToken(r.FormValue("state"), r.FormValue("code"))
			if !result {
				HandleLogin(w, r)
				return
			}
		}
		req.Header.Set("Authorization", "Bearer "+token.AccessToken)
		//Check if token is still valid
		response, err := client.Do(req)
		if err != nil || response.Status != "200 OK" {
			//Go to login if token is no longer valid
			HandleLogin(w, r)
			return
		}
		body, _ := ioutil.ReadAll(response.Body)
		var f interface{}
		json.Unmarshal(body, &f)
		m := f.(map[string]interface{})
		username := m["preferred_username"].(string)
		if login {
			loginLog(username)
			http.Redirect(w, r, "/Jtree/metadata/0.1.0/sample/search", http.StatusTemporaryRedirect)
		} else {
			//Go to redirect if token is still valid
			logAction(username, actionPageAccess, r.RequestURI)
		}
		if r.RequestURI == "/Jtree/metadata/0.1.0/logout" {
			Logout(w, r)
		}
		next.ServeHTTP(w, r)

	})
	//return function for page handling
	return handler
}

//Logout logs the user out
func Logout(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	url := keycloakserver + "/auth/realms/" + realm + "/protocol/openid-connect/userinfo"
	req, _ := http.NewRequest("GET", url, nil)
	if token == nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	//Check if token is still valid
	response, err := client.Do(req)
	if response.Status == "200 OK" && err == nil {
		body, _ := ioutil.ReadAll(response.Body)
		var f interface{}
		json.Unmarshal(body, &f)
		m := f.(map[string]interface{})
		username := m["preferred_username"].(string)
		//Go to redirect if token is still valid
		logAction(username, actionLogout, emptyString)
	}
	//Makes the logout page redirect to login page
	URI := server + "/login"
	//Logout using endpoint and redirect to login page
	http.Redirect(w, r, keycloakserver+"/auth/realms/"+realm+"/protocol/openid-connect/logout?redirect_uri="+URI, http.StatusTemporaryRedirect)

}

func getToken(state, code string) bool {
	if state != "" && oauthStateString != "" {
		//Checks that the strings are in a consistent state
		if state != oauthStateString {
			return false
		}
		//Gets the code from keycloak
		//Exchanges code for token
		token, err = oauth2Config.Exchange(context.Background(), code)
		if err != nil {
			return false
		}
		login = true
		return true
	}
	return false
}

func loginLog(username string) {
	logAction(username, actionLogin, emptyString)
	login = false
}
