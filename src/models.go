package src

type Email struct {
	Address  string `json:"address"`
	Password string `json:"password"`
}

type AppointmentParameters struct {
	City     string
	Category string
	Phone    string
	Email    Email
}
