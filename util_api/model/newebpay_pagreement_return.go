package model

import (
	"time"
	. "util_api/database/mysql"
)

type NewebpayPagreementReturn struct {
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
	RespondCode     string    `table:"respondCode"`
	Auth            string    `table:"auth"`
	Card6No         string    `table:"card6No"`
	Card4No         string    `table:"card4No"`
	Exp             string    `table:"exp"`
	ECI             string    `table:"eCI"`
	Inst            int64     `table:"inst"`
	InstFirst       int64     `table:"instFirst"`
	InstEach        int64     `table:"instEach"`
	TokenUseStatus  int64     `table:"tokenUseStatus"`
	TokenValue      string    `table:"tokenValue"`
	TokenLife       time.Time `table:"tokenLife"`
	PaymentMethod   string    `table:"paymentMethod"`
	ItemDesc        string    `table:"itemDesc"`
}

func (model *NewebpayPagreementReturn) SetUserId(userId int64) *NewebpayPagreementReturn {
	model.UserId = userId
	return model
}

func (model *NewebpayPagreementReturn) SetMerchantID(merchantID string) *NewebpayPagreementReturn {
	model.MerchantID = merchantID
	return model
}

func (model *NewebpayPagreementReturn) SetAmt(amt int64) *NewebpayPagreementReturn {
	model.Amt = amt
	return model
}

func (model *NewebpayPagreementReturn) SetTradeNo(tradeNo string) *NewebpayPagreementReturn {
	model.TradeNo = tradeNo
	return model
}

func (model *NewebpayPagreementReturn) SetMerchantOrderNo(merchantOrderNo string) *NewebpayPagreementReturn {
	model.MerchantOrderNo = merchantOrderNo
	return model
}

func (model *NewebpayPagreementReturn) SetPaymentType(paymentType string) *NewebpayPagreementReturn {
	model.PaymentType = paymentType
	return model
}

func (model *NewebpayPagreementReturn) SetRespondType(respondType string) *NewebpayPagreementReturn {
	model.RespondType = respondType
	return model
}

func (model *NewebpayPagreementReturn) SetPayTime(payTime time.Time) *NewebpayPagreementReturn {
	model.PayTime = payTime
	return model
}

func (model *NewebpayPagreementReturn) SetIP(iP string) *NewebpayPagreementReturn {
	model.IP = iP
	return model
}

func (model *NewebpayPagreementReturn) SetEscrowBank(escrowBank string) *NewebpayPagreementReturn {
	model.EscrowBank = escrowBank
	return model
}

func (model *NewebpayPagreementReturn) SetRespondCode(respondCode string) *NewebpayPagreementReturn {
	model.RespondCode = respondCode
	return model
}

func (model *NewebpayPagreementReturn) SetAuth(auth string) *NewebpayPagreementReturn {
	model.Auth = auth
	return model
}

func (model *NewebpayPagreementReturn) SetCard6No(card6No string) *NewebpayPagreementReturn {
	model.Card6No = card6No
	return model
}

func (model *NewebpayPagreementReturn) SetCard4No(card4No string) *NewebpayPagreementReturn {
	model.Card4No = card4No
	return model
}

func (model *NewebpayPagreementReturn) SetExp(exp string) *NewebpayPagreementReturn {
	model.Exp = exp
	return model
}

func (model *NewebpayPagreementReturn) SetECI(eCI string) *NewebpayPagreementReturn {
	model.ECI = eCI
	return model
}

func (model *NewebpayPagreementReturn) SetInst(inst int64) *NewebpayPagreementReturn {
	model.Inst = inst
	return model
}

func (model *NewebpayPagreementReturn) SetInstFirst(instFirst int64) *NewebpayPagreementReturn {
	model.InstFirst = instFirst
	return model
}

func (model *NewebpayPagreementReturn) SetInstEach(instEach int64) *NewebpayPagreementReturn {
	model.InstEach = instEach
	return model
}

func (model *NewebpayPagreementReturn) SetTokenUseStatus(tokenUseStatus int64) *NewebpayPagreementReturn {
	model.TokenUseStatus = tokenUseStatus
	return model
}

func (model *NewebpayPagreementReturn) SetTokenValue(tokenValue string) *NewebpayPagreementReturn {
	model.TokenValue = tokenValue
	return model
}

func (model *NewebpayPagreementReturn) SetTokenLife(tokenLife time.Time) *NewebpayPagreementReturn {
	model.TokenLife = tokenLife
	return model
}

func (model *NewebpayPagreementReturn) SetPaymentMethod(paymentMethod string) *NewebpayPagreementReturn {
	model.PaymentMethod = paymentMethod
	return model
}

func (model *NewebpayPagreementReturn) SetItemDesc(itemDesc string) *NewebpayPagreementReturn {
	model.ItemDesc = itemDesc
	return model
}

func (model *NewebpayPagreementReturn) Insert() (int64, error) {
	return Model(model).Insert()
}
