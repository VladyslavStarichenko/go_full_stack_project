package main

type Director struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Nationality    string `json:"nationality"`
	ContactPhone   string `json:"contact_phone"`
	ContactEmail   string `json:"contact_email"`
	ContactAddress string `json:"contact_address"`
}
