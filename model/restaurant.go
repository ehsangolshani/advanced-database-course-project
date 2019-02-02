package model

type Restaurant struct {
	DocumentId       string  `json:"document_id,omitempty"`
	Name             string  `json:"name,omitempty"`
	City             string  `json:"city,omitempty"`
	Country          string  `json:"country,omitempty"`
	AverageCost      int64   `json:"average_cost,omitempty"`
	PhoneNumber      string  `json:"phone_number,omitempty"`
	CashPayment      bool    `json:"cash_payment,omitempty"`
	CardPayment      bool    `json:"card_payment,omitempty"`
	Address          string  `json:"address,omitempty"`
	Rate             float32 `json:"rate,omitempty"`
	TakeoutAvailable bool    `json:"takeout_available,omitempty"`
	OutdoorSeating   bool    `json:"outdoor_seating,omitempty"`
	Hookah           bool    `json:"hookah,omitempty"`
	SmokingArea      bool    `json:"smoking_area,omitempty"`
	WifiAvailable    bool    `json:"wifi_available,omitempty"`
	Cuisine          string  `json:"cuisine,omitempty"`
}
