package line

import (
	models "eCommerce/models/frontend"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func SendOrderNotify(order models.OrderCreateRequest) {
	bot, err := linebot.New("282ebd89bcdf588ca0213228a616a134", "uDUSCi+C01hkDLjM/ry2dwzOIYn77JTWltHRToBLrak/+XauTTqkc4QkhuimFS9s2xZCaxAy+uctOEbU37jnN1dXyyBn4xqOvcmwDeC+ijgsqLl3o7cfNjjJEAsHUQVjTQ/jpu2/+z0izZmCWZ7ipQdB04t89/1O/w1cDnyilFU=")

	carouselContainer := &linebot.CarouselContainer{
		Type: linebot.FlexContainerTypeCarousel,
		Contents: []*linebot.BubbleContainer{
			{
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
							Text:   "往右滑看訂單商品，或到巨粽後台觀看訂單明細，以便趕快安排出貨喔~~",
							Weight: "bold",
							Size:   linebot.FlexTextSizeTypeMd,
							Wrap:   true,
							Color:  "#9d4edd",
						},
						&linebot.BoxComponent{
							Type:    linebot.FlexComponentTypeBox,
							Layout:  linebot.FlexBoxLayoutTypeHorizontal,
							Margin:  "lg",
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
									Text:  time.Unix(int64(order.CreatedAt), 0).Format("2006/01/02 15:04:05"),
									Color: "#666666",
									Size:  linebot.FlexTextSizeTypeSm,
									Flex:  linebot.IntPtr(2),
								},
							},
						},
						&linebot.BoxComponent{
							Type:    linebot.FlexComponentTypeBox,
							Layout:  linebot.FlexBoxLayoutTypeHorizontal,
							Margin:  "lg",
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
									Text:  fmt.Sprintf("$%d", int(order.Total)),
									Color: "#FF0000",
									Size:  linebot.FlexTextSizeTypeLg,
									Flex:  linebot.IntPtr(2),
								},
							},
						},
					},
				},
			},
		},
	}

	for _, product := range order.Products {
		for _, style := range product.Styles {
			carouselPage := linebot.BubbleContainer{
				Type: linebot.FlexContainerTypeBubble,
				Hero: &linebot.ImageComponent{
					Type:        linebot.FlexComponentTypeImage,
					URL:         style.Photo,
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
							Text:   style.Title,
							Weight: "bold",
							Size:   linebot.FlexTextSizeTypeXl,
						},
						&linebot.BoxComponent{
							Type:    linebot.FlexComponentTypeBox,
							Layout:  linebot.FlexBoxLayoutTypeHorizontal,
							Margin:  "lg",
							Spacing: "sm",
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
									Text:  style.StyleTitle,
									Wrap:  true,
									Color: "#666666",
									Size:  linebot.FlexTextSizeTypeSm,
									Flex:  linebot.IntPtr(5),
								},
							},
						},
						&linebot.BoxComponent{
							Type:    linebot.FlexComponentTypeBox,
							Layout:  linebot.FlexBoxLayoutTypeHorizontal,
							Margin:  "lg",
							Spacing: "sm",
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
									Text:  fmt.Sprintf("%d", int(style.DiscountedPrice)),
									Wrap:  true,
									Color: "#666666",
									Size:  linebot.FlexTextSizeTypeSm,
									Flex:  linebot.IntPtr(5),
								},
							},
						},
						&linebot.BoxComponent{
							Type:    linebot.FlexComponentTypeBox,
							Layout:  linebot.FlexBoxLayoutTypeHorizontal,
							Margin:  "lg",
							Spacing: "sm",
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
									Text:  fmt.Sprintf("%d", style.Qty),
									Wrap:  true,
									Color: "#666666",
									Size:  linebot.FlexTextSizeTypeSm,
									Flex:  linebot.IntPtr(5),
								},
							},
						},
					},
				},
			}
			carouselContainer.Contents = append(carouselContainer.Contents, &carouselPage)
		}
	}

	containerJson, _ := json.Marshal(carouselContainer)
	// fmt.Println(string(containerJson))
	container, err := linebot.UnmarshalFlexMessageJSON(containerJson)
	message := linebot.NewFlexMessage("您有一筆新訂單", container)

	_, err = bot.Multicast([]string{"U14fff345bfc700aa44170a860d851c23", "Ub79e88993077ecc98abc2a53711a5c9f"}, message).Do()
	if err != nil {
		fmt.Println(err)
		// Do something when some bad happened
	}
}
