package main

import (
	"blsParser/service"
	"blsParser/src"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

const configFile = "./config.json"

func main() {
	var config src.Config
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to read config %v", err))
	}
	if err := json.Unmarshal(file, &config); err != nil {
		panic(fmt.Sprintf("Failed to parse config %v", err))
	}

	tgClient := src.NewTgClient(config.TgToken)
	imapClient, err := src.NewImapClient(config.Email)
	defer imapClient.Logout()

	srv := service.NewService(imapClient, tgClient)

	for i := 1; i < 30; i++ {
		if err := srv.ParseSlots(src.AppointmentParameters{
			City:     "Москва",
			Category: "Обычная подача",
			Phone:    "9108934422",
			Email:    config.Email,
		}); err != nil {
			panic(fmt.Sprintf("Failed to parse slots %v", err))
		}
		time.Sleep(5 * time.Minute)
	}
}
