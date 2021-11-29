package email

import (
	"bytes"
	models "eCommerce/models/frontend"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func SendOrderNotify(order models.OrderCreateRequest) {
	// NewEmail返回一個email結構體的指標
	e := email.NewEmail()
	// 發件人
	e.From = "巨粽數位購物平台<hugelark813@gmail.com>"
	// 收件人(可以有多個)
	e.To = []string{"hugelark99@gmail.com"}
	// 郵件主題
	e.Subject = "訂單通知信件 - 巨粽數位購物平台"
	// 解析html模板
	t, err := template.ParseFiles("./internal/email/order_notify.html")
	if err != nil {
		fmt.Println(err)
	}
	body := new(bytes.Buffer)
	t.Execute(body, order)
	e.HTML = body.Bytes()
	e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "hugelark813@gmail.com", "gsgmwcmpginniapl", "smtp.gmail.com"))
}
