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
	HomeChargeFee   float32 `json:"home_charge_fee"`
	UniChargeFee    float32 `json:"uni_charge_fee"`
	FamilyChargeFee float32 `json:"family_charge_fee"`
	HilifeChargeFee float32 `json:"hilife_charge_fee"`
	OkChargeFee     float32 `json:"ok_charge_fee"`
	TimeDefault
}

func (PlatformLogistics) TableName() string {
	return "platform_logistics"
}

func (platform *PlatformLogistics) Fetch() {
	DB.Where("platform_id = ?", platform.PlatformID).First(&platform)
	return
}
