package model

import (
	"time"
	. "util_api/database/mysql"
)

type NewebpayNdnfCommonReturn struct {
	Id              int64     `table:"id"`
	UserId          int64     `table:"user_id"`
	MerchantID      string    `table:"merchantID"`
	Amt             int64     `table:"amt"`
	TradeNo         string    `table:"tradeNo"`
	MerchantOrderNo string    `table:"merchantOrderNo"`
	PaymentType     string    `table:"paymentType"`
	RespondType     string    `table:"respondType"`
	PayTime         time.Time `table:"payTime"`
	IP              string    `table:"iP"`
	EscrowBank      string    `table:"escrowBank"`
}

func (model *NewebpayNdnfCommonReturn) SetUserId(userId int64) *NewebpayNdnfCommonReturn {
	model.UserId = userId
	return model
}

func (model *NewebpayNdnfCommonReturn) SetMerchantID(merchantID string) *NewebpayNdnfCommonReturn {
	model.MerchantID = merchantID
	return model
}

func (model *NewebpayNdnfCommonReturn) SetAmt(amt int64) *NewebpayNdnfCommonReturn {
	model.Amt = amt
	return model
}

func (model *NewebpayNdnfCommonReturn) SetTradeNo(tradeNo string) *NewebpayNdnfCommonReturn {
	model.TradeNo = tradeNo
	return model
}

func (model *NewebpayNdnfCommonReturn) SetMerchantOrderNo(merchantOrderNo string) *NewebpayNdnfCommonReturn {
	model.MerchantOrderNo = merchantOrderNo
	return model
}

func (model *NewebpayNdnfCommonReturn) SetPaymentType(paymentType string) *NewebpayNdnfCommonReturn {
	model.PaymentType = paymentType
	return model
}

func (model *NewebpayNdnfCommonReturn) SetRespondType(respondType string) *NewebpayNdnfCommonReturn {
	model.RespondType = respondType
	return model
}

func (model *NewebpayNdnfCommonReturn) SetPayTime(payTime time.Time) *NewebpayNdnfCommonReturn {
	model.PayTime = payTime
	return model
}

func (model *NewebpayNdnfCommonReturn) SetIP(iP string) *NewebpayNdnfCommonReturn {
	model.IP = iP
	return model
}

func (model *NewebpayNdnfCommonReturn) SetEscrowBank(escrowBank string) *NewebpayNdnfCommonReturn {
	model.EscrowBank = escrowBank
	return model
}

func (model *NewebpayNdnfCommonReturn) Insert() (int64, error) {
	return Model(model).Insert()
}
