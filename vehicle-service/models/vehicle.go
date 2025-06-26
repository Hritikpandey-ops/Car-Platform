package models

type Vehicle struct {
	ID                 int    `json:"id"`
	Brand              string `json:"brand"`
	Model              string `json:"model"`
	Year               int    `json:"year"`
	Color              string `json:"color"`
	RegistrationNumber string `json:"registration_number"`
}
