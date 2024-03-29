package frontend

import (
	. "eCommerce/internal/database"
)

type PlatformQuery struct {
	Hostname string
}

type Platform struct {
	ID                 int    `json:"-"`
	Title              string `json:"title"`
	Description        string `json:"description"`
	LogoUrl            string `json:"logo_url"`
	IconUrl            string `json:"icon_url"`
	Code               string `json:"code"`
	FBPageID           string `json:"fb_page_id"`
	FBMessengerEnabled bool   `json:"fb_messenger_enabled"`
	FBPixel            string `json:"fb_pixel"`
}

func (Platform) TableName() string {
	return "platform"
}

func (query *PlatformQuery) Fetch() (platform Platform) {
	DB.Model(&Platform{}).Where("hostname = ?", query.Hostname).Scan(&platform)
	return
}

func (platform *Platform) GetMenu() (menus []Menus) {
	DB.Model(&Menus{}).Where("platform_id = ? and is_enabled = 1", platform.ID).Order("sort ASC").Scan(&menus)
	return
}

func (platform *Platform) GetPromotions() (promotions []Promotions) {
	DB.Model(&Promotions{}).Where("platform_id = ? AND is_enabled = 1 AND start_timestamp <= UNIX_TIMESTAMP() AND end_timestamp > UNIX_TIMESTAMP()", platform.ID).Scan(&promotions)
	return
}

func (platform *Platform) GetPayments() (payment PlatformPayment) {
	DB.Model(&PlatformPayment{}).Where("platform_id = ?", platform.ID).Scan(&payment)
	return
}

func (platform *Platform) GetLogistics() (logistics PlatformLogistics) {
	DB.Model(&PlatformLogistics{}).Where("platform_id = ?", platform.ID).Scan(&logistics)
	return
}
