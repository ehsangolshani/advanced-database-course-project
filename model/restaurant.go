package model

type Restaurant struct {
	DocumentId       string   `json:"document_id"`
	Name             string   `json:"name"`
	City             int      `json:"city"`
	Country          int      `json:"Country"`
	AverageCost      int64    `json:"average_cost"`
	PhoneNumber      string   `json:"phone_number"`
	CashPayment      bool     `json:"cash_payment"`
	CardPayment      bool     `json:"card_payment"`
	Address          string   `json:"address"`
	Rate             float32  `json:"rate"`
	TakeoutAvailable bool     `json:"takeout_available"`
	OutdoorSeating   bool     `json:"outdoor_seating"`
	Hookah           bool     `json:"hookah"`
	SmokingArea      bool     `json:"smoking_area"`
	WifiAvailable    bool     `json:"wifi_available"`
	Cuisines         []string `json:"cuisines"`
	Tags             []string `json:"tags"`
}
