package backends

import (
	"cms/pkg/common"
	. "cms/pkg/db"
	"cms/pkg/setting"
)

type DraftComponentData struct {
	ID    int                  `json:"id"`
	DmoID int                  `json:"dmo_id"`
	ComID int                  `json:"com_id"`
	Title string               `json:"title"`
	List  []DraftComponentList `json:"list" gorm:"-"`
	Img   string               `json:"img"`
	Link  string               `json:"link"`
	Text  string               `json:"text"`
	Lang  string               `json:"lang"`
}

func (comData *DraftComponentData) FetchByComID() {
	DB.Model(&comData).Where("com_id = ? AND lang = ?", comData.ComID, comData.Lang).Scan(&comData)
}

func (comData *DraftComponentData) Save() (err error) {
	transObj := map[string]map[string]string{
		"zh-Hant": {
			"title": comData.Title,
			"text":  comData.Text,
		},
	}

	transObj = common.Trans(transObj)

	for _, lang := range setting.Lang {
		comData.ID = 0
		comData.Lang = lang
		if lang != "zh-Hant" {
			comData.Title = transObj[lang]["title"]
			comData.Text = transObj[lang]["text"]
		}
		err = DB.Create(&comData).Error
	}

	return
}
