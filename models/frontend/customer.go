package frontend

import (
	. "eCommerce/internal/database"
)

type CustomerQuery struct {
	Hostname string
}

type Customer struct {
	ID      int
	Name    string
	LogoUrl string
}

func (query *CustomerQuery) Fetch() (customer Customer) {
	DB.Model(&Customer{}).Select("id, name, logo_url").Where("hostname = ?", query.Hostname).Scan(&customer)
	return
}
