package main

import (
	"fmt"
	"gator/internal/config"
)


func main(){
	username := "natedawg"
	configFilePath := config.GetConfigFilePath()

	cfg, err := config.Read(configFilePath)
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	err = cfg.SetUser(username, configFilePath)
	if err != nil {
		fmt.Printf("Error setting user: %v\n", err)
		return
	}

	cfg, err = config.Read(configFilePath)
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	fmt.Printf("Current user: %s\n", cfg.CurrentUserName)
	fmt.Printf("Postgres URL: %s\n", cfg.DBUrl)
}