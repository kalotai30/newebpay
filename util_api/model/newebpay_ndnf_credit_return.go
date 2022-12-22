package model

import (
	. "util_api/database/mysql"
)

type NewebpayNdnfCreditReturn struct {
	Id             int64  `table:"id"`
	NdnfId         int64  `table:"ndnf_id"`
	AuthBank       string `table:"authBank"`
	RespondCode    string `table:"respondCode"`
	Auth           string `table:"auth"`
	Card6No        string `table:"card6No"`
	Card4No        string `table:"card4No"`
	Exp            string `table:"exp"`
	Inst           int64  `table:"inst"`
	InstFirst      int64  `table:"instFirst"`
	InstEach       int64  `table:"instEach"`
	ECI            string `table:"eCI"`
	TokenUseStatus int64  `table:"tokenUseStatus"`
	RedAmt         int64  `table:"redAmt"`
	PaymentMethod  string `table:"paymentMethod"`
}

func (model *NewebpayNdnfCreditReturn) SetNdnfId(ndnfId int64) *NewebpayNdnfCreditReturn {
	model.NdnfId = ndnfId
	return model
}

func (model *NewebpayNdnfCreditReturn) SetAuthBank(authBank string) *NewebpayNdnfCreditReturn {
	model.AuthBank = authBank
	return model
}

func (model *NewebpayNdnfCreditReturn) SetRespondCode(respondCode string) *NewebpayNdnfCreditReturn {
	model.RespondCode = respondCode
	return model
}

func (model *NewebpayNdnfCreditReturn) SetAuth(auth string) *NewebpayNdnfCreditReturn {
	model.Auth = auth
	return model
}

func (model *NewebpayNdnfCreditReturn) SetCard6No(card6No string) *NewebpayNdnfCreditReturn {
	model.Card6No = card6No
	return model
}

func (model *NewebpayNdnfCreditReturn) SetCard4No(card4No string) *NewebpayNdnfCreditReturn {
	model.Card4No = card4No
	return model
}

func (model *NewebpayNdnfCreditReturn) SetExp(exp string) *NewebpayNdnfCreditReturn {
	model.Exp = exp
	return model
}

func (model *NewebpayNdnfCreditReturn) SetInst(inst int64) *NewebpayNdnfCreditReturn {
	model.Inst = inst
	return model
}

func (model *NewebpayNdnfCreditReturn) SetInstFirst(instFirst int64) *NewebpayNdnfCreditReturn {
	model.InstFirst = instFirst
	return model
}

func (model *NewebpayNdnfCreditReturn) SetInstEach(instEach int64) *NewebpayNdnfCreditReturn {
	model.InstEach = instEach
	return model
}

func (model *NewebpayNdnfCreditReturn) SetECI(eCI string) *NewebpayNdnfCreditReturn {
	model.ECI = eCI
	return model
}

func (model *NewebpayNdnfCreditReturn) SetTokenUseStatus(tokenUseStatus int64) *NewebpayNdnfCreditReturn {
	model.TokenUseStatus = tokenUseStatus
	return model
}

func (model *NewebpayNdnfCreditReturn) SetRedAmt(redAmt int64) *NewebpayNdnfCreditReturn {
	model.RedAmt = redAmt
	return model
}

func (model *NewebpayNdnfCreditReturn) SetPaymentMethod(paymentMethod string) *NewebpayNdnfCreditReturn {
	model.PaymentMethod = paymentMethod
	return model
}

func (model *NewebpayNdnfCreditReturn) Insert() (int64, error) {
	return Model(model).Insert()
}
