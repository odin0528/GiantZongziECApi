package backend

import (
	. "eCommerce/internal/database"
)

type PlatformLogistics struct {
	PlatformID      int  `json:"-"`
	HomeEnabled     bool `json:"home_enabled"`
	UniEnabled      bool `json:"uni_enabled"`
	FamilyEnabled   bool `json:"family_enabled"`
	HilifeEnabled   bool `json:"hilife_enabled"`
	OkEnabled       bool `json:"ok_enabled"`
	HomeChargeFee   int  `json:"home_charge_fee"`
	UniChargeFee    int  `json:"uni_charge_fee"`
	FamilyChargeFee int  `json:"family_charge_fee"`
	HilifeChargeFee int  `json:"hilife_charge_fee"`
	OkChargeFee     int  `json:"ok_charge_fee"`
	TimeDefault
}

func (PlatformLogistics) TableName() string {
	return "platform_logistics"
}

func (platform *PlatformLogistics) Fetch() {
	DB.Where("platform_id = ?", platform.PlatformID).First(&platform)
	return
}
