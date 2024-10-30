package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

type serve struct{}

type User struct {
	Username string    `json:"username"`
	Password string    `json:"password"`
	ID       uuid.UUID `json:"id"`
}

var store = sessions.NewCookieStore([]byte("6668fe14-f8cd-4be6-b896-3e04c1065da5"))

func main() {
	conf := LoadConfiguration()
	connect_db(conf.ConnectionString)
	setupDBSchema()
	seed_users()
	mux := http.NewServeMux()
	fmt.Println("Server is ready and listening on port 3000")
	mux.Handle("/", &serve{})
	http.ListenAndServe(":3000", mux)
}
