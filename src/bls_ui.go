package src

import (
	"github.com/playwright-community/playwright-go"
	"time"
)

const (
	blsURL           = "https://blsrussiaportugal.com/russian/appointment.php"
	CenterDropdown   = "//select[@name='centre']"
	CategoryDropdown = "//select[@name='category']"
	PhoneInput       = "//input[@name='phone']"
	EmailInput       = "//input[@name='email']"
	SendOtpText      = " Запросить проверочный код"
	OtpInput         = "//input[@name='otp']"
	ContinueButton   = ".btn-ctnue"
	AgreeButton      = ".primary-btn"
	datesTable       = ".app_date.validate"
	activeDays       = ".day.active.activeClass"
)

type Browser struct {
	page   playwright.Page
	params AppointmentParameters
}

func NewBlsUi(page *playwright.Page, params *AppointmentParameters) *Browser {
	return &Browser{
		page:   *page,
		params: *params,
	}
}

func (b Browser) FillFirstPage() error {
	var err error
	if _, err = b.page.Goto(blsURL); err != nil {
		return err
	}

	cityOption := []string{b.params.City}
	if _, err = b.page.Locator(CenterDropdown).SelectOption(playwright.SelectOptionValues{
		ValuesOrLabels: &cityOption,
	}); err != nil {
		return err
	}

	ctgrOption := []string{b.params.Category}
	if _, err = b.page.Locator(CategoryDropdown).SelectOption(playwright.SelectOptionValues{
		ValuesOrLabels: &ctgrOption,
	}); err != nil {
		return err
	}

	time.Sleep(1 * time.Second)
	if err = b.page.Locator(PhoneInput).Fill(b.params.Phone); err != nil {
		return err
	}
	if err = b.page.Locator(EmailInput).Fill(b.params.Email.Address); err != nil {
		return err
	}
	if err = b.page.GetByText(SendOtpText).Click(); err != nil {
		return err
	}

	time.Sleep(4 * time.Second)
	return nil
}

func (b Browser) FillOtp(otp string) error {
	if err := b.page.Locator(OtpInput).Fill(otp); err != nil {
		return err
	}

	if err := b.page.Locator(ContinueButton).Click(); err != nil {
		return err
	}
	return nil
}

func (b Browser) CheckSlots() ([]string, error) {

	if err := b.page.Locator(AgreeButton).Click(); err != nil {
		return nil, err
	}

	if err := b.page.Locator(datesTable).Click(); err != nil {
		return nil, err
	}

	var dts []string

	dt, err := b.page.Locator(activeDays).All()
	for _, d := range dt {
		dat, _ := d.TextContent()
		dts = append(dts, dat)
	}
	if err != nil {
	}

	return dts, nil
}
