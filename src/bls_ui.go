package src

import (
	"github.com/playwright-community/playwright-go"
	"log"
	"time"
)

const (
	CenterDropdown   = "//select[@name='centre']"
	CategoryDropdown = "//select[@name='category']"
	PhoneInput       = "//input[@name='phone']"
	EmailInput       = "//input[@name='email']"
	SendOtpText      = " Запросить проверочный код"
	OtpInput         = "//input[@name='otp']"
	ContinueButton   = ".btn-ctnue"
	AgreeButton      = ".primary-btn"
)

func GetSlots(page playwright.Page, params AppointmentParameters) error {
	var err error
	if _, err = page.Goto("https://blsrussiaportugal.com/russian/appointment.php"); err != nil {
		return err
	}

	cityOption := []string{params.City}
	if _, err = page.Locator(CenterDropdown).SelectOption(playwright.SelectOptionValues{
		ValuesOrLabels: &cityOption,
	}); err != nil {
		return err
	}

	ctgrOption := []string{params.Category}
	if _, err = page.Locator(CategoryDropdown).SelectOption(playwright.SelectOptionValues{
		ValuesOrLabels: &ctgrOption,
	}); err != nil {
		return err
	}

	time.Sleep(1 * time.Second)
	if err = page.Locator(PhoneInput).Fill(params.Phone); err != nil {
		return err
	}
	if err = page.Locator(EmailInput).Fill(params.Email.Address); err != nil {
		return err
	}
	//if err = page.GetByText(SendOtpText).Click(); err != nil {
	//	return err
	//}

	time.Sleep(2 * time.Second)

	mailClient, err := NewImapClient(params.Email)
	if err != nil {
		return err
	}
	defer mailClient.Logout()

	code, _, err := mailClient.GetLastOtp()
	if err != nil {
		return err
	}

	if err = page.Locator(OtpInput).Fill(*code); err != nil {
		return err
	}

	if err = page.Locator(ContinueButton).Click(); err != nil {
		return err
	}

	if err = page.Locator(AgreeButton).Click(); err != nil {
		return err
	}

	if err = page.Locator(".app_date.validate").Click(); err != nil {
		return err
	}
	//if err = page.Locator(".day.disabled.fullcap").WaitFor(playwright.LocatorWaitForOptions{
	//	State: playwright.WaitForSelectorStateVisible,
	//}); err != nil {
	//	return err
	//}

	var dts []string
	dates, _ := page.Locator(".day.disabled.fullcap").All()
	for _, d := range dates {
		dat, _ := d.TextContent()
		dts = append(dts, dat)
	}

	log.Println(dts)

	return nil
}
