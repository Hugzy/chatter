package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	ConnectionString string
}

func LoadConfiguration() Configuration {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}

	return configuration
}
