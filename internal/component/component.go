package component

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Component struct {
	TTextBlockA       ComponentConfig
	TTextBlockB       ComponentConfig
	TTextBlockC       ComponentConfig
	TTextImageBlockA  ComponentConfig
	TTextImageBlockB  ComponentConfig
	TBannerBlockA     ComponentConfig
	TBannerBlockB     ComponentConfig
	TBannerBlockC     ComponentConfig
	TPlaceBlockA      ComponentConfig
	TPlaceBlockB      ComponentConfig
	TPlaceBlockC      ComponentConfig
	TEventBlockA      ComponentConfig
	TEventBlockB      ComponentConfig
	TEventBlockC      ComponentConfig
	TAttractionBlockA ComponentConfig
	TAttractionBlockB ComponentConfig
	TAttractionBlockC ComponentConfig
	TStoreBlockA      ComponentConfig
	TStoreBlockB      ComponentConfig
	TStoreBlockC      ComponentConfig
	TFoodBlockA       ComponentConfig
	TFoodBlockB       ComponentConfig
	TFoodBlockC       ComponentConfig
	TFoodBlockD       ComponentConfig
	TProductBlockA    ComponentConfig
	TProductBlockB    ComponentConfig
	TPlayBlockA       ComponentConfig
	TPlayBlockB       ComponentConfig
	TPlayBlockC       ComponentConfig
	TVideoBlockA      ComponentConfig
	TVideoBlockB      ComponentConfig
	TMapBlockA        ComponentConfig
	TMessageBlockA    ComponentConfig
	TNewsBlockA       ComponentConfig
	TNewsBlockB       ComponentConfig
}

type ComponentConfig struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	MaxNumber int    `json:"maxNumber"`
}

var Config *Component

func init() {
	jsonFile, _ := os.Open("internal/component/componentList.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &Config)
}
