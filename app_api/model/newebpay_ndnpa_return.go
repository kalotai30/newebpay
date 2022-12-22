package model

import (
	. "app_api/database/mysql"
	"app_api/util/log"
)

type NewebpayNdnpaReturn struct {
	Id              int64  `table:"id"`
	UserId          int64  `table:"user_id"`
	MerchantID      string `table:"merchantID"`
	Amt             int64  `table:"amt"`
	TradeNo         string `table:"tradeNo"`
	MerchantOrderNo string `table:"merchantOrderNo"`
	RespondCode     string `table:"respondCode"`
	AuthBank        string `table:"authBank"`
	Auth            string `table:"auth"`
	AuthDate        string `table:"authDate"`
	AuthTime        string `table:"authTime"`
	Card6No         string `table:"card6No"`
	Card4No         string `table:"card4No"`
	Exp             string `table:"exp"`
	Inst            int64  `table:"inst"`
	InstFirst       int64  `table:"instFirst"`
	InstEach        int64  `table:"instEach"`
	ECI             string `table:"eCI"`
	RedAmt          int64  `table:"redAmt"`
	PaymentMethod   string `table:"paymentMethod"`
	IP              string `table:"iP"`
	EscrowBank      string `table:"escrowBank"`
}

func (model *NewebpayNdnpaReturn) SetUserId(userId int64) *NewebpayNdnpaReturn {
	model.UserId = userId
	return model
}

func (model *NewebpayNdnpaReturn) SetMerchantID(merchantID string) *NewebpayNdnpaReturn {
	model.MerchantID = merchantID
	return model
}

func (model *NewebpayNdnpaReturn) SetAmt(amt int64) *NewebpayNdnpaReturn {
	model.Amt = amt
	return model
}

func (model *NewebpayNdnpaReturn) SetTradeNo(tradeNo string) *NewebpayNdnpaReturn {
	model.TradeNo = tradeNo
	return model
}

func (model *NewebpayNdnpaReturn) SetMerchantOrderNo(merchantOrderNo string) *NewebpayNdnpaReturn {
	model.MerchantOrderNo = merchantOrderNo
	return model
}

func (model *NewebpayNdnpaReturn) SetRespondCode(respondCode string) *NewebpayNdnpaReturn {
	model.RespondCode = respondCode
	return model
}

func (model *NewebpayNdnpaReturn) SetAuthBank(authBank string) *NewebpayNdnpaReturn {
	model.AuthBank = authBank
	return model
}

func (model *NewebpayNdnpaReturn) SetAuth(auth string) *NewebpayNdnpaReturn {
	model.Auth = auth
	return model
}

func (model *NewebpayNdnpaReturn) SetAuthDate(authDate string) *NewebpayNdnpaReturn {
	model.AuthDate = authDate
	return model
}

func (model *NewebpayNdnpaReturn) SetAuthTime(authTime string) *NewebpayNdnpaReturn {
	model.AuthTime = authTime
	return model
}

func (model *NewebpayNdnpaReturn) SetCard6No(card6No string) *NewebpayNdnpaReturn {
	model.Card6No = card6No
	return model
}

func (model *NewebpayNdnpaReturn) SetCard4No(card4No string) *NewebpayNdnpaReturn {
	model.Card4No = card4No
	return model
}

func (model *NewebpayNdnpaReturn) SetExp(exp string) *NewebpayNdnpaReturn {
	model.Exp = exp
	return model
}

func (model *NewebpayNdnpaReturn) SetInst(inst int64) *NewebpayNdnpaReturn {
	model.Inst = inst
	return model
}

func (model *NewebpayNdnpaReturn) SetInstFirst(instFirst int64) *NewebpayNdnpaReturn {
	model.InstFirst = instFirst
	return model
}

func (model *NewebpayNdnpaReturn) SetInstEach(instEach int64) *NewebpayNdnpaReturn {
	model.InstEach = instEach
	return model
}

func (model *NewebpayNdnpaReturn) SetECI(eCI string) *NewebpayNdnpaReturn {
	model.ECI = eCI
	return model
}

func (model *NewebpayNdnpaReturn) SetRedAmt(redAmt int64) *NewebpayNdnpaReturn {
	model.RedAmt = redAmt
	return model
}

func (model *NewebpayNdnpaReturn) SetPaymentMethod(paymentMethod string) *NewebpayNdnpaReturn {
	model.PaymentMethod = paymentMethod
	return model
}

func (model *NewebpayNdnpaReturn) SetIP(iP string) *NewebpayNdnpaReturn {
	model.IP = iP
	return model
}

func (model *NewebpayNdnpaReturn) SetEscrowBank(escrowBank string) *NewebpayNdnpaReturn {
	model.EscrowBank = escrowBank
	return model
}

func (model *NewebpayNdnpaReturn) Insert() (int64, error) {
	return Model(model).Insert()
}

func (model *NewebpayNdnpaReturn) GetTradeNoByMerchantOrderNo() *NewebpayNdnpaReturn {
	table := Model(model)

	log.Error(table.Select([]string{"tradeNo", "amt"}).
		Where("merchantOrderNo", "=", model.MerchantOrderNo).
		Find().Scan(&model.TradeNo, &model.Amt))

	return model
}
