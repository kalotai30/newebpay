package config

import (
	"github.com/gin-gonic/gin"
)

type configNewebPay struct {
	MPGUrl        string
	APIUrl        string
	APICLoseUrl   string
	NotifyURL     string
	MerchantID_MS string
	HashKey_MS    string
	HashIV_MS     string
	MerchantID    string
	HashKey       string
	HashIV        string
}

var NewebPayInfo configNewebPay

func initConfigNewebPayInfo() {
	switch gin.Mode() {
	case gin.ReleaseMode:
		NewebPayInfo.MPGUrl = "https://core.newebpay.com/MPG/mpg_gateway"
		NewebPayInfo.APIUrl = "https://core.newebpay.com/API/CreditCard"
		NewebPayInfo.APICLoseUrl = "https://core.newebpay.com/API/CreditCard/Close"
		NewebPayInfo.NotifyURL = ""
		/* 平台商 */
		NewebPayInfo.MerchantID_MS = ""
		NewebPayInfo.HashKey_MS = ""
		NewebPayInfo.HashIV_MS = ""
		/* 合作商店 */
		NewebPayInfo.MerchantID = ""
		NewebPayInfo.HashKey = ""
		NewebPayInfo.HashIV = ""

	case gin.DebugMode:
		NewebPayInfo.MPGUrl = "https://ccore.newebpay.com/MPG/mpg_gateway"
		NewebPayInfo.APIUrl = "https://ccore.newebpay.com/API/CreditCard"
		NewebPayInfo.APICLoseUrl = "https://ccore.newebpay.com/API/CreditCard/Close"
		NewebPayInfo.NotifyURL = ""
		/* 平台商 */
		NewebPayInfo.MerchantID_MS = ""
		NewebPayInfo.HashKey_MS = ""
		NewebPayInfo.HashIV_MS = ""
		/* 合作商店 */
		NewebPayInfo.MerchantID = ""
		NewebPayInfo.HashKey = ""
		NewebPayInfo.HashIV = ""

	case gin.TestMode:
		NewebPayInfo.MPGUrl = "https://ccore.newebpay.com/MPG/mpg_gateway"
		NewebPayInfo.APIUrl = "https://ccore.newebpay.com/API/CreditCard"
		NewebPayInfo.APICLoseUrl = "https://core.newebpay.com/API/CreditCard/Close"
		NewebPayInfo.NotifyURL = ""
		/* 平台商 */
		NewebPayInfo.MerchantID_MS = ""
		NewebPayInfo.HashKey_MS = ""
		NewebPayInfo.HashIV_MS = ""
		/* 合作商店 */
		NewebPayInfo.MerchantID = ""
		NewebPayInfo.HashKey = ""
		NewebPayInfo.HashIV = ""

	}
}
