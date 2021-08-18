package backend

import (
	. "eCommerce/internal/database"

	"gorm.io/plugin/soft_delete"
)

type LoginReq struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type AdminResetPassword struct {
	AdminID   int
	Token     string
	ExpiredAt int
	CreatedAt int
	DeletedAt soft_delete.DeletedAt
}

func (AdminResetPassword) TableName() string {
	return "admin_reset_password"
}

func (reset *AdminResetPassword) CancelOldToken() {
	DB.Debug().Where("admin_id = ?", reset.AdminID).Delete(&AdminResetPassword{})
}

func (req *PageReq) Login() (pagesRowset []Pages, err error) {
	err = DB.Table("rel_customer_pages as rel").Select("rel.*, pages.name").Joins("inner join pages on rel.page_id = pages.id").
		Where("rel.customer_id = ?", req.CustomerID).
		Scan(&pagesRowset).Error
	return
}
