package frontend

import (
	. "eCommerce/internal/database"
	"eCommerce/internal/rdb"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type PlatformQuery struct {
	Hostname string
}

type Platform struct {
	ID                 int    `json:"-"`
	Title              string `json:"title"`
	Description        string `json:"description"`
	LogoUrl            string `json:"logo_url"`
	MobileLogoUrl      string `json:"mobile_logo_url"`
	IconUrl            string `json:"icon_url"`
	Hostname           string `json:"hostname"`
	Code               string `json:"code"`
	FBAppID            string `json:"fb_app_id"`
	FBPageID           string `json:"fb_page_id"`
	FBMessengerEnabled bool   `json:"fb_messenger_enabled"`
	FBPixel            string `json:"fb_pixel"`
	FBPixelToken       string `json:"-"`
	LineChannelID      string `json:"line_channel_id"`
}

func (Platform) TableName() string {
	return "platform"
}

func (query *PlatformQuery) Fetch() (platform Platform) {
	DB.Model(&Platform{}).Where("hostname = ?", query.Hostname).Scan(&platform)
	return
}

func (platform *Platform) GetMenu() (menus []Menus) {
	key := fmt.Sprintf("platform_%s_menu", platform.Code)
	err := rdb.Get(key, &menus)
	if err == redis.Nil {
		DB.Model(&Menus{}).Where("platform_id = ? and is_enabled = 1", platform.ID).Order("sort ASC").Scan(&menus)
		rdb.Set(key, menus)
	}
	return
}

func (platform *Platform) GetPromotions() (promotions []Promotions) {
	key := fmt.Sprintf("platform_%s_promotion", platform.Code)
	err := rdb.Get(key, &promotions)
	if err == redis.Nil {
		DB.Model(&Promotions{}).Where("platform_id = ? AND is_enabled = 1 AND start_timestamp <= UNIX_TIMESTAMP() AND end_timestamp > UNIX_TIMESTAMP()", platform.ID).Scan(&promotions)
		rdb.Set(key, promotions)
	}
	return
}

func (platform *Platform) GetPayments() (payment PlatformPayment) {
	key := fmt.Sprintf("platform_%s_payment", platform.Code)
	err := rdb.Get(key, &payment)
	if err == redis.Nil {
		DB.Model(&PlatformPayment{}).Where("platform_id = ?", platform.ID).Scan(&payment)
		rdb.Set(key, payment)
	}
	return
}
func (platform *Platform) GetCategory() (sortCategory []Category) {

	key := fmt.Sprintf("platform_%s_category", platform.Code)
	err := rdb.Get(key, &sortCategory)
	if err == redis.Nil {
		categories := []Category{}
		DB.Model(&Category{}).Where("platform_id = ?", platform.ID).Order("layer ASC, sort ASC").Scan(&categories)
		for _, category := range categories {
			if category.ParentID == -1 {
				sortCategory = append(sortCategory, category)
			} else {
				FindParentCategory(&sortCategory, category)
			}
		}
		rdb.Set(key, sortCategory)
	}

	return
}

func FindParentCategory(parentCategory *[]Category, child Category) {
	for index := range *parentCategory {
		if (*parentCategory)[index].ID == child.ParentID {
			(*parentCategory)[index].Child = append((*parentCategory)[index].Child, child)
			return
		}
		FindParentCategory(&(*parentCategory)[index].Child, child)
	}
}

func (platform *Platform) GetLogistics() (logistics PlatformLogistics) {
	key := fmt.Sprintf("platform_%s_logistics", platform.Code)
	err := rdb.Get(key, &logistics)
	if err == redis.Nil {
		DB.Model(&PlatformLogistics{}).Where("platform_id = ?", platform.ID).Scan(&logistics)
		rdb.Set(key, logistics)
	} else {
		println("logistics is cached")
	}
	return
}
