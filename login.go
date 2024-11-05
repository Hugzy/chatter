package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func getUserByName(username string) (*User, error) {
	p := User{}
	err := db.Get(&p, "SELECT * FROM users where username = $1 LIMIT 1", username)
	if err != nil {
		fmt.Println(err)
		return nil, &UserNotFound{}
	}

	return &p, nil
}

func isLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	fmt.Println("checking login status")

	sess, err := store.Get(r, "login")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(sess)

	return sess.Values["accesskey"] != nil
}

func register(w http.ResponseWriter, r *http.Request) error {
	// Get the user from the signup request
	tx := db.MustBegin()

	var p User
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	hashsedPassword, err := bcrypt.GenerateFromPassword([]byte(p.Password), 0)
	if err != nil {
		return err
	}

	tx.MustExec("INSERT INTO users (username, password) values $1, $2", p.Username, hashsedPassword)
	err = tx.Commit()
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			panic(err)
		}
		return err
	}

	return nil
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "login")

	for k := range session.Values {
		fmt.Print("deleting: ")
		fmt.Println(k)
		delete(session.Values, k)
	}

	session.Save(r, w)

	fmt.Println(session)

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

	user, error := getUserByName(p.Username)
	if error != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	session, _ := store.Get(r, "login")

	// User was already logged in
	if isLoggedIn(w, r) {
		w.Write([]byte("Already logged in"))
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
		session.Values["accesskey"] = user.ID.String()
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// session.Values["accesskey"] =
		w.Write([]byte("Login successful"))
	}
}
