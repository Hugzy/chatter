package main

import "net/http"

func getWeather(w http.ResponseWriter, r *http.Request) (*string, error) {
	if !isLoggedIn(w, r) {
		return nil, &NotAuthenticated{}
	}
	l := ""
	return &l, nil
}
