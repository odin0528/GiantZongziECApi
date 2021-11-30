package line

import (
	backend "eCommerce/models/backend"
	models "eCommerce/models/frontend"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func SendOrderNotify(container linebot.FlexContainer) {
	bot, err := linebot.New(os.Getenv("LINE_MESSAGE_API_SECRET"), os.Getenv("LINE_MESSAGE_API_TOKEN"))
	message := linebot.NewFlexMessage("您有一筆新訂單", container)
	_, err = bot.Multicast([]string{"U14fff345bfc700aa44170a860d851c23", "Ub79e88993077ecc98abc2a53711a5c9f"}, message).Do()
	// _, err = bot.Multicast([]string{"U14fff345bfc700aa44170a860d851c23"}, message).Do()
	if err != nil {
		fmt.Println(err)
		// Do something when some bad happened
	}
}

func SendOrderNotifyByOrderCreateRequest(order models.OrderCreateRequest) {
	carouselContainer := GenerateOrderNotifyContainer(int(order.Total), order.CreatedAt, order.Memo)

	for _, product := range order.Products {
		for _, style := range product.Styles {
			carouselContainer.Contents = append(carouselContainer.Contents, GenerateOrderNotifyOrderItemBubble(
				style.Title,
				style.StyleTitle,
				style.Photo,
				style.DiscountedPrice,
				style.Qty,
			))
		}
	}

	containerJson, _ := json.Marshal(carouselContainer)
	fmt.Println(string(containerJson))
	container, _ := linebot.UnmarshalFlexMessageJSON(containerJson)
	SendOrderNotify(container)
}

func SendOrderNotifyByOrder(order backend.Orders) {
	carouselContainer := GenerateOrderNotifyContainer(int(order.Total), order.CreatedAt, order.Memo)

	for _, product := range order.Products {
		carouselContainer.Contents = append(carouselContainer.Contents, GenerateOrderNotifyOrderItemBubble(
			product.Title,
			product.StyleTitle,
			product.Photo,
			product.DiscountedPrice,
			product.Qty,
		))
	}

	containerJson, _ := json.Marshal(carouselContainer)
	fmt.Println(string(containerJson))
	container, _ := linebot.UnmarshalFlexMessageJSON(containerJson)
	SendOrderNotify(container)
}

func GenerateOrderNotifyContainer(total int, createdAt int, memo string) *linebot.CarouselContainer {
	carouselContainer := &linebot.CarouselContainer{
		Type: linebot.FlexContainerTypeCarousel,
		Contents: []*linebot.BubbleContainer{
			GenerateOrderNotifyFirstBubble(total, createdAt, memo),
		},
	}
	return carouselContainer
}

func GenerateOrderNotifyFirstBubble(total int, createdAt int, memo string) *linebot.BubbleContainer {
	p := message.NewPrinter(language.English)
	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Hero: &linebot.ImageComponent{
			Type:        linebot.FlexComponentTypeImage,
			URL:         "https://www.giantzongzi.com/img/projects/project-home-1.png",
			Size:        "full",
			AspectRatio: "1:1",
			AspectMode:  "cover",
			Action:      &linebot.URIAction{URI: os.Getenv("EC_ADMIN_URL") + "/orders"},
		},
		Body: &linebot.BoxComponent{
			Type:    linebot.FlexComponentTypeBox,
			Layout:  linebot.FlexBoxLayoutTypeVertical,
			Margin:  "lg",
			Spacing: "sm",
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "往左滑看訂單商品，以便趕快回到巨粽後台安排出貨喔~~",
					Weight: "bold",
					Size:   linebot.FlexTextSizeTypeMd,
					Wrap:   true,
					Color:  "#9d4edd",
				},
				&linebot.BoxComponent{
					Type:    linebot.FlexComponentTypeBox,
					Layout:  linebot.FlexBoxLayoutTypeHorizontal,
					Margin:  "sm",
					Spacing: "sm",
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  "訂單日期",
							Color: "#aaaaaa",
							Size:  linebot.FlexTextSizeTypeSm,
							Flex:  linebot.IntPtr(1),
						},
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  time.Unix(int64(createdAt), 0).Format("2006/01/02 15:04:05"),
							Color: "#666666",
							Size:  linebot.FlexTextSizeTypeSm,
							Flex:  linebot.IntPtr(2),
						},
					},
				},
				&linebot.BoxComponent{
					Type:    linebot.FlexComponentTypeBox,
					Layout:  linebot.FlexBoxLayoutTypeHorizontal,
					Margin:  "sm",
					Spacing: "sm",
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  "訂單金額",
							Color: "#aaaaaa",
							Size:  linebot.FlexTextSizeTypeSm,
							Flex:  linebot.IntPtr(1),
						},
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  p.Sprintf("$%d", int(total)),
							Color: "#FF0000",
							Size:  linebot.FlexTextSizeTypeLg,
							Flex:  linebot.IntPtr(2),
						},
					},
				},
				&linebot.BoxComponent{
					Type:    linebot.FlexComponentTypeBox,
					Layout:  linebot.FlexBoxLayoutTypeHorizontal,
					Margin:  "sm",
					Spacing: "sm",
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  "備註",
							Color: "#aaaaaa",
							Size:  linebot.FlexTextSizeTypeSm,
							Flex:  linebot.IntPtr(1),
						},
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  memo,
							Color: "#666666",
							Size:  linebot.FlexTextSizeTypeLg,
							Flex:  linebot.IntPtr(2),
						},
					},
				},
			},
		},
	}
}

func GenerateOrderNotifyOrderItemBubble(title string, styleTitle string, photo string, price float32, qty int) *linebot.BubbleContainer {
	p := message.NewPrinter(language.English)
	if photo[:8] != "https://" {
		photo = "https://giant-zongzi-ec.s3.ap-northeast-1.amazonaws.com" + photo
	}
	carouselPage := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Hero: &linebot.ImageComponent{
			Type:        linebot.FlexComponentTypeImage,
			URL:         photo,
			Size:        "full",
			AspectRatio: "1:1",
			AspectMode:  "cover",
			Action:      &linebot.URIAction{URI: os.Getenv("EC_ADMIN_URL") + "/orders"},
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   title,
					Weight: "bold",
					Size:   linebot.FlexTextSizeTypeMd,
				},
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeHorizontal,
					Margin: "sm",
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  "規格",
							Color: "#aaaaaa",
							Size:  linebot.FlexTextSizeTypeSm,
							Flex:  linebot.IntPtr(1),
						},
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  styleTitle,
							Wrap:  true,
							Color: "#666666",
							Size:  linebot.FlexTextSizeTypeSm,
							Flex:  linebot.IntPtr(5),
						},
					},
				},
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeHorizontal,
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  "單價",
							Color: "#aaaaaa",
							Size:  linebot.FlexTextSizeTypeSm,
							Flex:  linebot.IntPtr(1),
						},
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  p.Sprintf("$%d", int(price)),
							Wrap:  true,
							Color: "#666666",
							Size:  linebot.FlexTextSizeTypeSm,
							Flex:  linebot.IntPtr(5),
						},
					},
				},
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeHorizontal,
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  "數量",
							Color: "#aaaaaa",
							Size:  linebot.FlexTextSizeTypeSm,
							Flex:  linebot.IntPtr(1),
						},
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  fmt.Sprintf("%d", qty),
							Wrap:  true,
							Color: "#666666",
							Size:  linebot.FlexTextSizeTypeSm,
							Flex:  linebot.IntPtr(5),
						},
					},
				},
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeHorizontal,
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  "小計",
							Color: "#aaaaaa",
							Size:  linebot.FlexTextSizeTypeSm,
							Flex:  linebot.IntPtr(1),
						},
						&linebot.TextComponent{
							Type:   linebot.FlexComponentTypeText,
							Text:   p.Sprintf("$%d", qty*int(price)),
							Wrap:   true,
							Color:  "#ff0000",
							Size:   linebot.FlexTextSizeTypeMd,
							Flex:   linebot.IntPtr(5),
							Weight: "bold",
						},
					},
				},
			},
		},
	}

	return carouselPage
}
