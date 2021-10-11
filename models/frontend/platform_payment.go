package frontend

import (
	. "eCommerce/internal/database"
)

type PlatformPayment struct {
	PlatformID      int    `json:"-"`
	TransferEnabled bool   `json:"transfer_enabled"`
	TransferBank    string `json:"transfer_bank"`
	TransferAccount string `json:"transfer_account"`
	DeliveryEnabled bool   `json:"delivery_enabled"`
	Delivery711     bool   `json:"delivery_711"`
	DeliveryFamily  bool   `json:"delivery_family"`
	DeliveryHilife  bool   `json:"delivery_hilife"`
	DeliveryOK      bool   `json:"delivery_ok"`
	LinePayEnabled  bool   `json:"line_pay_enabled"`
	TimeDefault
}

func (PlatformPayment) TableName() string {
	return "platform_payment"
}

func (platform *PlatformPayment) Fetch() {
	DB.Debug().Where("platform_id = ?", platform.PlatformID).First(&platform)
	return
}
