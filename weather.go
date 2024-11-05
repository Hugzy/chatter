package main

import (
	"math/rand"
	"net/http"
	"time"
)

type Forecast struct {
	Date         string `json:"date"`
	Summary      string `json:"summary"`
	TemperatureC int    `json:"temperaturec"`
}

var summaries = []string{"Freezing", "Bracing", "Chilly", "Cool", "Mild", "Warm", "Balmy", "Hot", "Sweltering", "Scorching"}

func getWeather(w http.ResponseWriter, r *http.Request) ([]Forecast, error) {
	if !isLoggedIn(w, r) {
		return nil, &NotAuthenticated{}
	}

	forecasts := []Forecast{}

	min := -20
	max := 55

	for i := 1; i <= 5; i++ {
		d := time.Now().AddDate(0, 0, i)
		formattedD := d.Format("02-Jan-2006")
		temp := rand.Intn((max - min) + min)
		summary := summaries[rand.Intn(len(summaries))]

		forecasts = append(forecasts, Forecast{
			formattedD,
			summary,
			temp,
		})
	}

	return forecasts, nil
}
