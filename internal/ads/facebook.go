package ads

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type FbConversionParams struct {
	EventName      string                 `json:"event_name"`
	EventTime      int64                  `json:"event_time"`
	EventSourceUrl string                 `json:"event_source_url"`
	ActionSource   string                 `json:"action_source"`
	UserData       FbConversionUserData   `json:"user_data"`
	CustomData     FbConversionCustomData `json:"custom_data"`
}

type FbConversionUserData struct {
	EM              string `json:"em"`
	ExternalID      string `json:"external_id"`
	ClientIPAddress string `json:"client_ip_address"`
	ClientUserAgent string `json:"client_user_agent"`
}

type FbConversionCustomData struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
}

func Send(params FbConversionParams, pixelID string, pixelToken string) {
	str, _ := json.Marshal([]FbConversionParams{params})
	resp, _ := http.Post(
		fmt.Sprintf("https://graph.facebook.com/v12.0/%s/events?access_token=%s", pixelID, pixelToken),
		"application/x-www-form-urlencoded",
		strings.NewReader("data="+string(str)+"&test_event_code=TEST53275"),
	)

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
}
