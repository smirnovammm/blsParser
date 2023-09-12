package service

import (
	"blsParser/src"
	"context"
	"fmt"
	"time"
)

type Service struct {
	imapClient *src.ImapClient
	tgClient   *src.TgClient
	browser    *src.Browser
}

func NewService(imapClient *src.ImapClient, tgClient *src.TgClient, browser *src.Browser) *Service {
	return &Service{
		imapClient: imapClient,
		tgClient:   tgClient,
		browser:    browser,
	}
}

func (s Service) ParseSlots() error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	otpSendTime := time.Now()
	if err := s.browser.FillFirstPage(); err != nil {
		return err
	}
	time.Sleep(5 * time.Second)

	code, _, err := s.imapClient.GetLastOtp(&otpSendTime)
	if err != nil {
		return err
	}

	if err = s.browser.FillOtp(*code); err != nil {
		return err
	}

	dates, err := s.browser.CheckSlots()
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)

	if err = s.tgClient.SendToTg(fmt.Sprintf("Free slots: %s", dates)); err != nil {
		return err
	}

	return nil
}
