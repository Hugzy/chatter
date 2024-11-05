package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *serve) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	fmt.Printf("Request URL: %s\n", r.URL.Path)
	switch r.Method {
	case "GET":
		Get(w, r)
	case "POST":
		Post(w, r)
	default:
		http.Error(w, "Method not supported", http.StatusNotFound)
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/weather":
		forecasts, err := getWeather(w, r)
		encoder := json.NewEncoder(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		encoder.Encode(forecasts)
		return
	case "/test":
		w.Write([]byte(`{"message": "GET /test called"}`))
		return
	case "/migrations":
		GetMigrations()
	case "/users":
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		users, err := getAllUsers()
		if err != nil {
			panic(err)
		}
		json.NewEncoder(w).Encode(users)
		return
	case "/test/afterlogin":
		isLoggedIn(w, r)
		// testLoggedIn(w, r)
	default:
		http.Error(w, "Path Not found", http.StatusNotFound)
	}
}

func Post(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/register":
		register(w, r)
		return
	case "/login":
		login(w, r)
		return
	case "/logout":
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
