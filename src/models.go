package src

type Email struct {
	Address  string `json:"address"`
	Password string `json:"password"`
}

type Config struct {
	Email   Email  `json:"email"`
	TgToken string `json:"tg_token"`
}

type AppointmentParameters struct {
	City     string
	Category string
	Phone    string
	Email    Email
}
