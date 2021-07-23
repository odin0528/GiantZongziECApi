package backends

import (
	. "cms/pkg/db"

	"github.com/jinzhu/gorm"
)

type DraftComponentReq struct {
	DmoID int    `json:"dmo_id" uri:"id"`
	Type  string `json:"type" uri:"type"`
}

type DraftComponent struct {
	ID            int                  `json:"id" gorm:"primaryKey;autoIncrement"`
	DmoID         int                  `json:"dmo_id"`
	Sort          int                  `json:"sort"`
	ComponentName string               `json:"componentName"`
	Title         string               `json:"title"`
	Type          string               `json:"type"`
	List          []DraftComponentData `json:"list" gorm:"-"`
}

type DraftComponentChangeReq struct {
	DmoID     int `json:"dmo_id"`
	Position1 int `json:"position_1"`
	Position2 int `json:"position_2"`
}

func (req *DraftComponentReq) DmoComponentFetch() (component Component) {
	DB.Debug().Model(&component).Where("dmo_id = ? and type = ?", req.DmoID, req.Type).Scan(&component)
	return
}

func (component *DraftComponent) FetchByDmoID() (components []DraftComponent) {
	DB.Model(&component).Where("dmo_id = ?", component.DmoID).Order("sort asc").Scan(&components)
	return
}

func (component *DraftComponent) Save() (err error) {
	component.ID = 0
	err = DB.Model(&DraftComponent{}).Where("dmo_id = ? AND sort >= ?", component.DmoID, component.Sort).Update("sort", gorm.Expr("sort + 1")).Error
	err = DB.Create(&component).Error
	return
}

func (component *DraftComponent) FetchBySort() {
	DB.Model(component).Where("dmo_id = ? AND sort = ?", component.DmoID, component.Sort).Scan(&component)
}

func (component *DraftComponent) Delete() (err error) {
	DB.Delete(DraftComponent{}, "id = ?", component.ID)
	err = DB.Model(&DraftComponent{}).Where("dmo_id = ? AND sort > ?", component.DmoID, component.Sort).Update("sort", gorm.Expr("sort - 1")).Error
	return
}

func (component *DraftComponent) DeleteChildren() (err error) {
	DB.Delete(DraftComponentData{}, "com_id = ?", component.ID)
	DB.Delete(DraftComponentList{}, "com_id = ?", component.ID)
	return
}

func (component *DraftComponent) Update() (err error) {
	err = DB.Save(&component).Error
	return
}
