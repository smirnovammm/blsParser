package main

import (
	"blsParser/src"
	"encoding/json"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"io/ioutil"
)

const configFile = "./config.json"

func main() {
	var mailDate src.Email
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to read config %v", err))
	}
	if err := json.Unmarshal(file, &mailDate); err != nil {
		panic(fmt.Sprintf("Failed to parse config %v", err))
	}

	pw, err := playwright.Run()
	if err != nil {
		panic(fmt.Sprintf("could not start playwright: %v", err))
	}
	defer pw.Stop()

	headless := false

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: &headless,
	})
	defer browser.Close()
	if err != nil {
		panic(fmt.Sprintf("could not launch browser: %v", err))
	}

	page, err := browser.NewPage()
	if err != nil {
		panic(fmt.Sprintf("could not create page: %v", err))
	}

	if err = src.GetSlots(page, src.AppointmentParameters{
		City:     "Москва",
		Category: "Обычная подача",
		Phone:    "9108934422",
		Email:    mailDate,
	}); err != nil {
		panic(fmt.Sprintf("Unable to get avalible slots %v", err))
	}
}
