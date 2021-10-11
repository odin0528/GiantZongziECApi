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
	FBMessengerEnabled string `json:"fb_messenger_enabled"`
	FBPixel            string `json:"fb_pixel"`
}

func (Platform) TableName() string {
	return "platform"
}

func (query *PlatformQuery) Fetch() (platform Platform) {
	DB.Model(&Platform{}).Where("hostname = ?", query.Hostname).Scan(&platform)
	return
}

func (platform *Platform) GetMenu() (pages []Pages) {
	DB.Model(&Pages{}).Where("platform_id = ? AND is_menu = 1 AND is_enabled = 1 AND released_at > 0", platform.ID).Scan(&pages)
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
