package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

var (
	cache redis.Conn
)

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

// Credentials ...
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func init() {

	conn, err := redis.DialURL("redis://localhost")

	if err != nil {
		panic(err)
	}

	fmt.Println("> successfully connected to redis cache on :6379")

	cache = conn

}

func main() {

	fmt.Println("> starting up session_auth server on :8000")

	defer cache.Close()

	http.HandleFunc("/signin", SignInHandler)
	http.HandleFunc("/welcome", WelcomeHandler)
	http.HandleFunc("/refresh", RefreshHandler)

	log.Fatal(http.ListenAndServe(":8000", nil))

}

// SignInHandler ...
func SignInHandler(w http.ResponseWriter, r *http.Request) {

	var credentials Credentials

	err := json.NewDecoder(r.Body).Decode(&credentials)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expected, ok := users[credentials.Username]

	if !ok || expected != credentials.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionToken := uuid.New().String()

	_, err = cache.Do("SETEX", sessionToken, "120", credentials.Username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(120 * time.Second),
	})

}

// WelcomeHandler ...
func WelcomeHandler(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie("session_token")

	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionToken := c.Value

	response, err := cache.Do("GET", sessionToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if response == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Write([]byte(fmt.Sprintf("Welcome %s!", response)))

}

// RefreshHandler ...
func RefreshHandler(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	response, err := cache.Do("GET", sessionToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if response == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newSessionToken := uuid.New().String()
	_, err = cache.Do("SETEX", newSessionToken, "120", fmt.Sprintf("%s", response))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = cache.Do("DEL", sessionToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   newSessionToken,
		Expires: time.Now().Add(120 * time.Second),
	})
}
