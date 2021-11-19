package backend

import (
	. "eCommerce/internal/database"
)

type PlatformPayment struct {
	PlatformID        int  `json:"-"`
	DeliveryEnabled   bool `json:"delivery_enabled"`
	CreditCardEnabled bool `json:"credit_card_enabled"`
	LinePayEnabled    bool `json:"line_pay_enabled"`
	AtmEnabled        bool `json:"atm_enabled"`
	CvsEnabled        bool `json:"cvs_enabled"`
	BarcodeEnabled    bool `json:"barcode_enabled"`
	TimeDefault
}

func (PlatformPayment) TableName() string {
	return "platform_payment"
}

func (platform *PlatformPayment) Fetch() {
	DB.Where("platform_id = ?", platform.PlatformID).First(&platform)
	return
}
