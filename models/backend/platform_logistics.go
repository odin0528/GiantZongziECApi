package backend

import (
	. "eCommerce/internal/database"
)

type PlatformLogistics struct {
	PlatformID      int     `json:"-"`
	HomeEnabled     bool    `json:"home_enabled"`
	UniEnabled      bool    `json:"uni_enabled"`
	FamilyEnabled   bool    `json:"family_enabled"`
	HilifeEnabled   bool    `json:"hilife_enabled"`
	OkEnabled       bool    `json:"ok_enabled"`
	HomeChargeFee   float64 `json:"home_charge_fee"`
	UniChargeFee    float64 `json:"uni_charge_fee"`
	FamilyChargeFee float64 `json:"family_charge_fee"`
	HilifeChargeFee float64 `json:"hilife_charge_fee"`
	OkChargeFee     float64 `json:"ok_charge_fee"`
	TimeDefault
}

func (PlatformLogistics) TableName() string {
	return "platform_logistics"
}

func (platform *PlatformLogistics) Fetch() {
	DB.Where("platform_id = ?", platform.PlatformID).First(&platform)
	return
}
