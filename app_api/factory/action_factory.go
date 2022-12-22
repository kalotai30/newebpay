package factory

import (
	"app_api/handler"
)

var (
	myPayHandler          handler.MyPayHandler
	ocpHandler            handler.OcpHandler
	digitalPaymentHandler handler.DigitalPaymentHandler
)

var ActionFactoryAuth = map[string]interface{}{}

var ActionFactory = map[string]interface{}{
	//======================================================================
	//							OCP Electronic Wallet Api (高鉅電子錢包)
	//======================================================================
	"GetForwardingUrl": &ocpHandler,

	//======================================================================
	//							MyPayHandler Api(高鉅金流)
	//======================================================================
	"MyPayment":   &myPayHandler,
	"OCPEPayment": &myPayHandler,

	//======================================================================
	//							DigitalPaymentHandler Api(電子支付)
	//======================================================================
	"AgreedCreditCardAuthorization": &digitalPaymentHandler,
	"PayByAgreedCreditCard":         &digitalPaymentHandler,
	"PayByCreditCard":               &digitalPaymentHandler,
	"PayByCreditCardCompleted":      &digitalPaymentHandler,
	"UserPayTokenByType":            &digitalPaymentHandler,
	"RefundByCreditCard":            &digitalPaymentHandler,
	"CreditCardPaymentCheckByOrder": &digitalPaymentHandler,
}
