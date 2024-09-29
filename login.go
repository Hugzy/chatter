package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var users = []User{}

func saveUser(u User) {
	users = append(users, u)
}

func getUserById(id uuid.UUID) User {
	idx := slices.IndexFunc(users, func(u User) bool { return u.ID == id })
	if idx == -1 {
		return User{}
	}
	return users[idx]
}

func getUserByName(username string) User {
	idx := slices.IndexFunc(users, func(u User) bool { return u.Username == username })
	if idx == -1 {
		return User{}
	}
	return users[idx]
}

func isLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	fmt.Println("checking login status")

	sess, _ := store.Get(r, "login")

	return sess == nil
}

func register(w http.ResponseWriter, r *http.Request) error {
	// Get the user from the signup request
	tx, err := db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	var p User
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	_, err = tx.Exec(context.Background(), "INSERT INTO User")
	if err != nil {
		return err
	}

	return nil
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "login")

	for k := range session.Values {
		delete(session.Values, k)
	}

	w.Write([]byte("Logout successful"))
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var p User
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := getUserByName(p.Username)
	if (User{}) == user {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	session, _ := store.Get(r, "login")

	// User was already logged in
	if !session.IsNew {
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(p.Password))

	if err == nil {
		// login successful
		session.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   21600,
			HttpOnly: true,
		}
		// TODO could add and store an accesskey here with an expiration
		session.Values["id"] = user.ID.String()
		session.Values["username"] = user.Username
		// session.Values["accesskey"] =
		w.Write([]byte("Login successful"))
	}
}
