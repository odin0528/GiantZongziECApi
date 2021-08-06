package uploader

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	env "github.com/joho/godotenv"
	"github.com/nfnt/resize"
)

var Uploader *s3manager.Uploader

func init() {
	env.Load()
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), ""),
	})

	Uploader = s3manager.NewUploader(sess)
}

func Thumbnail(filename string, base64String string, maxSize uint) string {
	i := strings.Index(base64String, ",") //有找到base64的編碼關鍵字
	blob, _ := base64.StdEncoding.DecodeString(base64String[i+1:])
	var img image.Image
	buff := new(bytes.Buffer)

	switch strings.TrimSuffix(base64String[5:i], ";base64") {
	case "image/png":
		img, _ = png.Decode(bytes.NewReader(blob))
		thumbnail := resize.Resize(maxSize, 0, img, resize.Lanczos3)
		png.Encode(buff, thumbnail)
	case "image/jpeg":
		img, _ = jpeg.Decode(bytes.NewReader(blob))
		thumbnail := resize.Resize(maxSize, 0, img, resize.Lanczos3)
		jpeg.Encode(buff, thumbnail, nil)
	}

	result, _ := Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(buff.Bytes()),
	})

	return result.Location
}
