package main

type ProductionCompany struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Country        string `json:"country"`
	ContactPhone   string `json:"contact_phone"`
	ContactEmail   string `json:"contact_email"`
	ContactAddress string `json:"contact_address"`
}
