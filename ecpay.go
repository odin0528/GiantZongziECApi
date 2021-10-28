import (
	"github.com/Laysi/go-ecpay-sdk"
	"github.com/Laysi/go-ecpay-sdk/base"
)

func main() {
	// client := ecpay.NewClient("2000132", "5294y06JbISpM5x9", "v77hoKGq4kWxNNIS", "<RETURN_URL>")
	client := ecpay.NewStageClient(ecpay.WithReturnURL("https://ec.giantzongzi.com/ec/return"), ecpay.WithDebug)
}