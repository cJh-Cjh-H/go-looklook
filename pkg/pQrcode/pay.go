package pQrcode

import (
	"fmt"
	qrcode "github.com/skip2/go-qrcode"
)

type Pay struct {
	UserId  int64
	AuthKey string
	Sn      string
}

func PayQrcode(pay *Pay) ([]byte, error) {
	// 2. 生成支付回调URL
	// 注意：这里必须是实际可访问的URL
	callbackUrl := fmt.Sprintf("%s/thirdPayment/thirdPaymentWxPayCallback/%s", "http://localhost:9993/payment/v1", pay.Sn)
	qrBytes, err := qrcode.Encode(callbackUrl, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	return qrBytes, nil
}
