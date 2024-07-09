package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type serve struct{}

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}

var store = sessions.NewCookieStore([]byte("6668fe14-f8cd-4be6-b896-3e04c1065da5"))

var users = []User{}



func Get(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/test":
		fmt.Printf("here?")
		w.Write([]byte(`{"message": "GET /test called"}`))
		return
    case r.URL.Path == "/get/migrations":
        id, name, err := GetMigrations();
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        fmt.Printf("id: %v - name: %v\n", id, name)
	case r.URL.Path == "/users":
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)
		return
	default:
		http.Error(w, "Path Not found", http.StatusNotFound)
	}
}

func Post(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/register":
		register(w, r)
		return
	case r.URL.Path == "/login":
		login(w, r)
		return
	case r.URL.Path == "/logout":
		logout(w, r)
	default:
		http.Error(w, "Path Not found", http.StatusNotFound)
	}
}

func Put(w http.ResponseWriter, r *http.Request) {
	panic("Not Implemented")
}

func Delete(w http.ResponseWriter, r *http.Request) {
	panic("Not Implemented")
}

func (h *serve) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	fmt.Printf("Request URL: %s\n", r.URL.Path)
	switch {
	case r.Method == "GET":
		Get(w, r)
	case r.Method == "POST":
		Post(w, r)
	default:
		http.Error(w, "Method not supported", http.StatusNotFound)
	}
}

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

func isLoggedIn(sess *sessions.Session, u User, w http.ResponseWriter) bool {
	if sess.Values["id"] == u.ID.String() {
		http.Error(w, "Already logged in", http.StatusOK)
		return true
	}

	return false
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

	fmt.Println(session.Values)

	if isLoggedIn(session, user, w) {
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
		session.Values["id"] = user.ID.String()
		session.Values["username"] = user.Username
		session.Save(r, w)
		w.Write([]byte("Login successful"))
	} else {
		http.Error(w, "Invalid password", http.StatusBadRequest)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "login")
	session.Options.MaxAge = -1
	session.Save(r, w)
	w.Write([]byte("Logged out"))
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var p User

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.ID = uuid.New()

	idx := slices.IndexFunc(users, func(u User) bool { return u.Username == p.Username })

	if idx == -1 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusBadRequest)
			return
		} else {
			p.Password = string(hashedPassword)
			saveUser(p)
			return
		}
	}

	w.Write([]byte("User already exists"))
}

func main() {
	conf := LoadConfiguration()
	connect_db(conf.ConnectionString)
	setupDBSchema()

	mux := http.NewServeMux()
	fmt.Println("Server is ready and listening on port 3000")
	mux.Handle("/", &serve{})
	http.ListenAndServe(":3000", mux)
}
