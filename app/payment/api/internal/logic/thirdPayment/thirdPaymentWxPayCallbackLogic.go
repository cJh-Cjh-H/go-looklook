package thirdPayment

import (
	"context"
	"go-zero-looklook/app/order/rpc/order"
	"net/http"

	"go-zero-looklook/app/payment/api/internal/svc"
	"go-zero-looklook/app/payment/api/internal/types"
	"go-zero-looklook/app/payment/model"
	"go-zero-looklook/pkg/xerr"

	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrWxPayCallbackError = xerr.NewErrMsg("wechat pay callback fail")

type Payment struct {
	Sn             string
	Payer          string
	SuccessTime    string
	TradeState     string
	TradeStateDesc string
	TradeType      string
	TransactionId  string
}
type ThirdPaymentcallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type ThirdPaymentWxPayCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewThirdPaymentWxPayCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) ThirdPaymentWxPayCallbackLogic {
	return ThirdPaymentWxPayCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ThirdPaymentWxPayCallbackLogic) ThirdPaymentWxPayCallback(rw http.ResponseWriter, req *http.Request, sn string) (*types.ThirdPaymentWxPayCallbackResp, error) {

	p := &Payment{
		Sn:         sn,
		TradeState: SUCCESS,
	}

	returnCode := "SUCCESS"
	err := l.verifyAndUpdateState(p)
	if err != nil {
		returnCode = "FAIL"
	}

	return &types.ThirdPaymentWxPayCallbackResp{
		ReturnCode: returnCode,
	}, err
}

// Verify and update relevant flow data
func (l *ThirdPaymentWxPayCallbackLogic) verifyAndUpdateState(pay *Payment) error {

	//paymentResp, err := l.svcCtx.PaymentRpc.GetPaymentBySn(l.ctx, &payment.GetPaymentBySnReq{
	//	Sn: pay.Sn,
	//})
	//fmt.Printf("paymentResp:%v", paymentResp)
	//if err != nil || paymentResp.PaymentDetail.Id == 0 {
	//	return errors.Wrapf(ErrWxPayCallbackError, "Failed to get payment flow record err:%v ,notifyTrasaction:%+v ", err, pay)
	//}

	//// Judgment status
	//payStatus := l.getPayStatusByWXPayTradeState(pay.TradeState)
	//if payStatus == model.ThirdPaymentPayTradeStateSuccess {
	//	//付款通知Payment Notification.
	//	if paymentResp.PaymentDetail.PayStatus != model.ThirdPaymentPayTradeStateWait {
	//		return nil
	//	}
	//
	//	// Update the flow status.
	//	if _, err = l.svcCtx.PaymentRpc.UpdateTradeState(l.ctx, &payment.UpdateTradeStateReq{
	//		Sn:             pay.Sn,
	//		TradeState:     "已支付",
	//		TransactionId:  pay.TransactionId,
	//		TradeType:      "微信支付",
	//		TradeStateDesc: pay.TradeStateDesc,
	//		PayStatus:      l.getPayStatusByWXPayTradeState(pay.TradeState),
	//	}); err != nil {
	//		return errors.Wrapf(ErrWxPayCallbackError, "更新流水状态失败  err:%v , notifyTrasaction:%v ", err, pay)
	//	}
	//
	//} else if payStatus == model.ThirdPaymentPayTradeStateWait {
	//	//退款通知。Refund notification @todo to be done later, not needed at this time
	//
	//}
	orderResp, err := l.svcCtx.OrderRpc.HomestayOrderDetail(l.ctx, &order.HomestayOrderDetailReq{
		Sn: pay.Sn,
	})
	if orderResp.HomestayOrder.TradeState != 0 {
		return errors.New("只有未支付的订单才可以支付")
	}
	_, err = l.svcCtx.OrderRpc.UpdateHomestayOrderTradeState(l.ctx, &order.UpdateHomestayOrderTradeStateReq{
		Sn:         pay.Sn,
		TradeState: 1,
	})
	if err != nil {
		return errors.Wrapf(err, "Payment将订单状态改为1时错误")
	}

	return nil

}

const (
	SUCCESS    = "SUCCESS"    //支付成功
	REFUND     = "REFUND"     //转入退款
	NOTPAY     = "NOTPAY"     //未支付
	CLOSED     = "CLOSED"     //已关闭
	REVOKED    = "REVOKED"    //已撤销（付款码支付）
	USERPAYING = "USERPAYING" //用户支付中（付款码支付）
	PAYERROR   = "PAYERROR"   //支付失败(其他原因，如银行返回失败)
)

func (l *ThirdPaymentWxPayCallbackLogic) getPayStatusByWXPayTradeState(wxPayTradeState string) int64 {
	switch wxPayTradeState {
	case SUCCESS: //支付成功
		return model.ThirdPaymentPayTradeStateSuccess
	case USERPAYING: //支付中
		return model.ThirdPaymentPayTradeStateWait
	case REFUND: //已退款
		return model.ThirdPaymentPayTradeStateWait
	default:
		return model.ThirdPaymentPayTradeStateFAIL
	}
}
