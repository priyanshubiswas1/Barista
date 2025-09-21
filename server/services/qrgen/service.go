package internal

import (
	"fmt"

	"github.com/skip2/go-qrcode"
)

func GenerateQR(url string) []byte {
	qrCode, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return []byte(fmt.Sprintf("error:%s", err.Error()))
	}
	return qrCode
}
