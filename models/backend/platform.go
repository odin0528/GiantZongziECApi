package backend

import (
	. "eCommerce/internal/database"
)

type Platform struct {
	ID                 int    `json:"-"`
	Title              string `json:"title"`
	Description        string `json:"description"`
	Code               string `json:"code"`
	LogoUrl            string `json:"logo_url"`
	MobileLogoUrl      string `json:"mobile_logo_url"`
	IconUrl            string `json:"icon_url"`
	Hostname           string `json:"hostname"`
	FBPixel            string `json:"fb_pixel"`
	FBPageID           string `json:"fb_page_id"`
	FBMessengerEnabled bool   `json:"fb_messenger_enabled"`
	TimeDefault
}

func (Platform) TableName() string {
	return "platform"
}

func (platform *Platform) Fetch() {
	DB.First(&platform)
	return
}
