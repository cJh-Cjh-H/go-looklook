package thirdPayment

import (
	"net/http"

	"go-zero-looklook/app/payment/api/internal/logic/thirdPayment"
	"go-zero-looklook/app/payment/api/internal/svc"
	"go-zero-looklook/app/payment/api/internal/types"
	"go-zero-looklook/pkg/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ThirdPaymentwxPayHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ThirdPaymentWxPayReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := thirdPayment.NewThirdPaymentwxPayLogic(r.Context(), ctx)
		resp, err := l.ThirdPaymentwxPay(req)
		if err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Write(resp)
		//result.HttpResult(r, w, resp, err)
	}
}
