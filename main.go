package main

import (
	"blsParser/service"
	"blsParser/src"
	"blsParser/src/config"
	"log"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
	}

	tgClient := src.NewTgClient(cfg.TgToken)
	imapClient, err := src.NewImapClient(cfg.Email.Address, cfg.Email.Password)
	if err != nil {
	}
	defer imapClient.Logout()

	blsUi, err := src.NewBlsUi(&src.AppointmentParameters{
		City:     "Москва",
		Category: "Обычная подача",
		Phone:    "9108934422",
		Email:    cfg.Email.Address,
	})
	if err != nil {
	}

	srv := service.NewService(imapClient, tgClient, blsUi)

	for i := 1; i < 30; i++ {
		if err := srv.ParseSlots(); err != nil {
			log.Printf("Failed to parse slots, retrying in 1 minute %s", err)
			time.Sleep(1 + time.Minute)
			continue
		}
		time.Sleep(5 * time.Minute)
	}
}
