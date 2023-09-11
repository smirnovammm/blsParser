package service

import (
	"blsParser/src"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"time"
)

type Service struct {
	imapClient *src.ImapClient
	tgClient   *src.TgClient
}

func NewService(imapClient *src.ImapClient, tgClient *src.TgClient) *Service {
	return &Service{
		imapClient: imapClient,
		tgClient:   tgClient,
	}
}

func (s Service) ParseSlots(params src.AppointmentParameters) error {
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

	blsUi := src.NewBlsUi(&page, &params)

	if err = blsUi.FillFirstPage(); err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	code, _, err := s.imapClient.GetLastOtp()
	if err != nil {
		return err
	}

	if err = blsUi.FillOtp(*code); err != nil {
		return err
	}

	dates, err := blsUi.CheckSlots()
	if err != nil {
		return err
	}
	if err = s.tgClient.SendToTg(fmt.Sprintf("Free slots: %s", dates)); err != nil {
		return err
	}

	return nil
}
