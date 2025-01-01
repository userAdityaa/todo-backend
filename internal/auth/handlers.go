package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var user struct {
	Profile string `json:"profile"`
	Name    string `json:"name"`
}

func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := GetGoogleAuthURL()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	fmt.Println(url)
}

func GoogleCallBackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("coming here...")
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "code not found", http.StatusBadRequest)
		return
	}

	userInfo, err := HandleGoogleCallBack(code)

	if err != nil {
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal([]byte(userInfo), &user); err != nil {
		http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
		return
	}

	redirectURL := fmt.Sprintf("http://localhost:3000/home?profile=%s&name=%s", url.QueryEscape(user.Profile), url.QueryEscape(user.Name))

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
