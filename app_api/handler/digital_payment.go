package handler

import (
	"app_api/config"
	"app_api/content"
	"app_api/model"
	"app_api/repository"
	"app_api/util"
	"app_api/util/crypto"
	digitalpay "app_api/util/digitalPay"
	"app_api/util/log"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type DigitalPaymentHandler content.Handler

type DigitalPayReceived struct {
	OrderTmpId int64   `json:"orderTmpId"`
	OrderId    int64   `json:"orderId"`
	TokenType  int64   `json:"tokenType"`
	PaymetType int64   `json:"paymentType"`
	ItemDesc   string  `json:"itemDesc"`
	TotalPrice float64 `json:"totalPrice"`
}

type DigitalPayCheckReceived struct {
	OrderId      int64  `json:"orderId"`
	TakeoutId    int64  `json:"takeoutId"`
	IsAccepted   bool   `json:"isAccepted"`
	CancelReason string `json:"cancelReason"`
}

func (handler *DigitalPaymentHandler) PayByCreditCard() interface{} {
	var (
		data               DigitalPayReceived
		orderData          model.Orders
		orderTmpRepository repository.OrderTmpRepository
		tradeInfoMap       map[string]interface{}
	)

	if err := json.Unmarshal([]byte(handler.Parameter), &data); err != nil {
		return util.RS{Message: "", Status: false}
	} else {
		if data.OrderId > 0 {
			// 抓訂單編號
			orderData.QueryOneById(data.OrderId)

			// 產出呼叫藍新金流api需帶的參數
			tradeInfoMap = map[string]interface{}{
				"MerchantOrderNo": orderData.OrderSn,
				"Amt":             data.TotalPrice,
				"ItemDesc":        data.ItemDesc,
				"CREDIT":          1,
			}
		} else {
			// 抓暫存訂單資料
			mapResult := orderTmpRepository.GetOrderTempData(data.OrderTmpId)

			// 產出呼叫藍新金流api需帶的參數
			tradeInfoMap = map[string]interface{}{
				"MerchantOrderNo": data.OrderTmpId,
				"Amt":             int64(mapResult["totalPrice"].(float64)),
				"ItemDesc":        data.ItemDesc,
				"CREDIT":          1,
			}
		}
		returnData := digitalpay.NDNF101(tradeInfoMap)

		return util.RS{Message: "信用卡付款_幕前", Status: true, Data: returnData}
	}
}

func (handler *DigitalPaymentHandler) PayByCreditCardCompleted() interface{} {
	var (
		data         DigitalPayReceived
		orderTmp     model.OrderTmp
		takeoutOrder model.TakeoutOrder
	)

	if err := json.Unmarshal([]byte(handler.Parameter), &data); err != nil {
		return util.RS{Message: "", Status: false}
	} else {
		orderTmp.SetId(int64(data.OrderTmpId)).GetOrderData()
		takeoutOrder.SetOrderId(orderTmp.OrderId).QueryOne()

		returnData := map[string]interface{}{
			"orderTmpId": orderTmp.Id,
			"orderId":    orderTmp.OrderId,
			"orderSn":    orderTmp.OrderSn,
			"takeoutId":  takeoutOrder.Id,
		}

		return util.RS{Message: "信用卡付款完成狀態回傳", Status: true, Data: returnData}
	}
}

func (handler *DigitalPaymentHandler) AgreedCreditCardAuthorization() interface{} {
	if handler.UserId <= 0 {
		return util.RS{Message: "無登入會員", Status: false}
	}

	var users model.Users

	userInfo := users.GetUsersInfoWithId(handler.UserId)

	// 產出呼叫藍新金流api需帶的參數
	tradeInfoMap := map[string]interface{}{
		"MerchantOrderNo": "P" + strconv.FormatInt(handler.UserId, 10) + "T" + strconv.FormatInt(time.Now().Unix(), 10),
		"Email":           userInfo["email"],
		"TokenTerm":       handler.UserId,
	}
	returnData := digitalpay.PAgreementAuthorize(tradeInfoMap)

	return util.RS{Message: "約定信用卡付款授權_幕前", Status: true, Data: returnData}
}

func (handler *DigitalPaymentHandler) PayByAgreedCreditCard() interface{} {
	if handler.UserId <= 0 {
		return util.RS{Message: "無登入會員", Status: false}
	}

	var (
		data                DigitalPayReceived
		users               model.Users
		userPayToken        model.UserPayToken
		newebpayNdnpaReturn model.NewebpayNdnpaReturn
		orderTmpRepository  repository.OrderTmpRepository
	)

	if err := json.Unmarshal([]byte(handler.Parameter), &data); err != nil {
		return util.RS{Message: "", Status: false}
	} else {
		// 抓暫存訂單資料
		mapResult := orderTmpRepository.GetOrderTempData(data.OrderTmpId)

		userInfo := users.GetUsersInfoWithId(handler.UserId)
		userPayToken.SetUserId(handler.UserId).SetTokenType(1).QueryOne()

		// 呼叫藍新金流api需帶的參數
		tradeInfoMap := map[string]interface{}{
			"TimeStamp":       strconv.FormatInt(time.Now().Unix(), 10),
			"Version":         "1.6",
			"MerchantOrderNo": data.OrderTmpId,
			"Amt":             int64(mapResult["totalPrice"].(float64)),
			"ProdDesc":        "綁定授權後付款測試",
			"PayerEmail":      userInfo["email"],
			"TokenValue":      userPayToken.Token,
			"TokenTerm":       handler.UserId,
			"TokenSwitch":     "on",
		}

		// 接收藍新金流api處理完成資訊
		returnData := digitalpay.NdnpaB10(tradeInfoMap, config.NewebPayInfo.MerchantID_MS, config.NewebPayInfo.HashKey_MS, config.NewebPayInfo.HashIV_MS)

		if returnData["Status"] == "SUCCESS" {
			/* 寫入藍新回傳的資訊 */
			_, err := newebpayNdnpaReturn.
				SetUserId(int64(mapResult["userId"].(float64))).
				SetMerchantID(digitalpay.ReturnValueStr(returnData["MerchantID"])).
				SetAmt(digitalpay.ReturnValueInt64(returnData["Amt"])).
				SetTradeNo(digitalpay.ReturnValueStr(returnData["TradeNo"])).
				SetMerchantOrderNo(digitalpay.ReturnValueStr(returnData["MerchantOrderNo"])).
				SetRespondCode(digitalpay.ReturnValueStr(returnData["RespondCode"])).
				SetAuthBank(digitalpay.ReturnValueStr(returnData["AuthBank"])).
				SetAuth(digitalpay.ReturnValueStr(returnData["Auth"])).
				SetAuthDate(digitalpay.ReturnValueStr(returnData["AuthDate"])).
				SetAuthTime(digitalpay.ReturnValueStr(returnData["AuthTime"])).
				SetCard6No(digitalpay.ReturnValueStr(returnData["Card6No"])).
				SetCard4No(digitalpay.ReturnValueStr(returnData["Card4No"])).
				SetExp(digitalpay.ReturnValueStr(returnData["Exp"])).
				SetInst(digitalpay.ReturnValueInt64(returnData["Inst"])).
				SetInstFirst(digitalpay.ReturnValueInt64(returnData["InstFirst"])).
				SetInstEach(digitalpay.ReturnValueInt64(returnData["InstEach"])).
				SetECI(digitalpay.ReturnValueStr(returnData["ECI"])).
				SetRedAmt(digitalpay.ReturnValueInt64(returnData["RedAmt"])).
				SetPaymentMethod(digitalpay.ReturnValueStr(returnData["PaymentMethod"])).
				SetIP(digitalpay.ReturnValueStr(returnData["IP"])).
				SetEscrowBank(digitalpay.ReturnValueStr(returnData["EscrowBank"])).
				Insert()
			if err != nil {
				return util.RS{Message: "綁定授權付款寫入藍新回傳資訊出現問題，請連絡系統管理員", Status: false}
			}

			/* 完成寫入訂單程序 */
			orderSn := orderTmpRepository.SaveToOrder(mapResult, data.OrderTmpId)
			if orderSn != "false" {
				return util.RS{Message: "綁定授權後付款_幕後", Status: true, Data: orderSn}
			} else {
				return util.RS{Message: "綁定授權付款寫入訂單程序出現問題，請連絡系統管理員", Status: false}
			}
		} else {
			return util.RS{Message: "綁定授權付款讀取藍新回傳資料出現問題，請連絡系統管理員", Status: false, Data: returnData}
		}
	}
}

func (handler *DigitalPaymentHandler) PayByInstallment() interface{} {
	var (
		data               DigitalPayReceived
		orderData          model.Orders
		orderTmpRepository repository.OrderTmpRepository
		tradeInfoMap       map[string]interface{}
	)

	if err := json.Unmarshal([]byte(handler.Parameter), &data); err != nil {
		return util.RS{Message: "", Status: false}
	} else {
		if data.OrderId > 0 {
			// 抓訂單編號
			orderData.QueryOneById(data.OrderId)

			// 產出呼叫藍新金流api需帶的參數
			tradeInfoMap = map[string]interface{}{
				"MerchantOrderNo": orderData.OrderSn,
				"Amt":             data.TotalPrice,
				"ItemDesc":        data.ItemDesc,
				"InstFlag":        "3,6,12",
			}
		} else {
			// 抓暫存訂單資料
			mapResult := orderTmpRepository.GetOrderTempData(data.OrderTmpId)

			// 產出呼叫藍新金流api需帶的參數
			tradeInfoMap = map[string]interface{}{
				"MerchantOrderNo": data.OrderTmpId,
				"Amt":             int64(mapResult["totalPrice"].(float64)),
				"ItemDesc":        data.ItemDesc,
				"InstFlag":        "3,6,12",
			}
		}
		returnData := digitalpay.NDNF101(tradeInfoMap)

		return util.RS{Message: "信用卡分期付款_幕前", Status: true, Data: returnData}
	}
}

func (handler *DigitalPaymentHandler) UserPayTokenByType() interface{} {
	var (
		data                     DigitalPayReceived
		userPayToken             model.UserPayToken
		newebpayPagreementReturn model.NewebpayPagreementReturn
	)

	if err := json.Unmarshal([]byte(handler.Parameter), &data); err != nil {
		return util.RS{Message: "", Status: false}
	} else {
		userPayToken.SetUserId(handler.UserId).SetTokenType(data.TokenType).QueryOne()

		if userPayToken.Token != "" {
			creditCardNumber := ""
			expired := true

			//判斷是否過期
			if time.Until(userPayToken.TokenLife) > 0 {
				expired = false
			} else {
				expired = true
			}

			if userPayToken.TokenType == 1 {
				newebpayPagreementReturn.SetUserId(handler.UserId).SetTokenValue(userPayToken.Token).QueryOne()
				creditCardNumber = newebpayPagreementReturn.Card6No + "******" + newebpayPagreementReturn.Card4No
			}

			returnData := map[string]interface{}{
				"tokenLife":        userPayToken.TokenLife.Format("2006-01-02"),
				"tokenType":        userPayToken.TokenType,
				"creditCardNumber": creditCardNumber,
				"expired":          expired,
			}

			return util.RS{Message: "金流token", Status: true, Data: returnData}
		} else {
			return util.RS{Message: "查無金流token", Status: false}
		}
	}
}

func (handler *DigitalPaymentHandler) RefundByCreditCard() interface{} {
	var (
		data                     DigitalPayReceived
		orderTmp                 model.OrderTmp
		newebpayNdnpaReturn      model.NewebpayNdnpaReturn
		newebpayNdnfCommonReturn model.NewebpayNdnfCommonReturn
		tradeInfoMap             map[string]interface{}
		returnData               map[string]interface{}
	)

	if err := json.Unmarshal([]byte(handler.Parameter), &data); err != nil {
		return util.RS{Message: "", Status: false}
	} else if orderTmpId := orderTmp.SetOrderId(data.OrderId).GetTmpIdByOrderId(); orderTmpId <= 0 {
		return util.RS{Message: "信用卡退款_查無訂單編號", Status: false}
	} else if newebpayNdnpaReturn.SetMerchantOrderNo(strconv.FormatInt(orderTmpId, 10)).GetTradeNoByMerchantOrderNo(); newebpayNdnpaReturn.TradeNo != "" {
		// 產出呼叫藍新金流api需帶的參數
		tradeInfoMap = map[string]interface{}{
			"Amt":             newebpayNdnpaReturn.Amt,
			"MerchantOrderNo": strconv.FormatInt(orderTmpId, 10),
			"TradeNo":         newebpayNdnpaReturn.TradeNo,
		}
	} else if newebpayNdnfCommonReturn.SetMerchantOrderNo(strconv.FormatInt(orderTmpId, 10)).GetTradeNoByMerchantOrderNo(); newebpayNdnfCommonReturn.TradeNo != "" {
		// 產出呼叫藍新金流api需帶的參數
		tradeInfoMap = map[string]interface{}{
			"Amt":             newebpayNdnfCommonReturn.Amt,
			"MerchantOrderNo": strconv.FormatInt(orderTmpId, 10),
			"TradeNo":         newebpayNdnfCommonReturn.TradeNo,
		}
	} else {
		return util.RS{Message: "信用卡退款", Status: false}
	}

	returnData = digitalpay.NDNF101NPAB032(tradeInfoMap)
	/* 呼叫藍新金流 api */
	resp, err := http.PostForm(config.NewebPayInfo.APICLoseUrl,
		url.Values{"MerchantID_": {returnData["MerchantID"].(string)},
			"PostData_": {returnData["PostData"].(string)}})
	if err != nil {
		return util.RS{Message: "藍新金流_退款幕前_取回回傳值失敗", Status: false}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return util.RS{Message: "藍新金流_退款幕前_讀取body失敗", Status: false}
	}

	// 解析 非3D 回覆值
	respBody := string(body)
	respBodyMap := crypto.HttpParseQuery(respBody)

	return util.RS{Message: "信用卡退款", Status: true, Data: respBodyMap}
}

func (handler *DigitalPaymentHandler) CreditCardPaymentCheckByOrder() interface{} {
	var (
		checkData          DigitalPayCheckReceived
		takeoutOrder       model.TakeoutOrder
		orderTmpRepository repository.OrderTmpRepository
	)

	if err := json.Unmarshal([]byte(handler.Parameter), &checkData); err != nil {
		log.Error(err)
		return util.RS{Message: "", Status: false}
	} else if checkData.OrderId == 0 && checkData.TakeoutId == 0 {
		log.Error(errors.New("data.OrderId == 0 && data.TakeoutId == 0"))
		return util.RS{Message: "", Status: false}
	} else if checkData.CancelReason == "" {
		return util.RS{Message: "請選擇取消訂單原因", Status: false}
	} else if takeoutOrder.SetId(checkData.TakeoutId).SetOrderId(checkData.OrderId).SetUserId(handler.UserId).QueryOneJoinOrder(); takeoutOrder.Name == "" {
		return util.RS{Message: "查無此訂單", Status: false}
	} else if takeoutOrder.Status == 4 {
		return util.RS{Message: "訂單已取消，無法重複取消訂單", Status: false}
	} else if takeoutOrder.Status == 5 {
		return util.RS{Message: "訂單已完成，無法取消訂單", Status: false}
	} else if takeoutOrder.Status != 0 && checkData.IsAccepted == false {
		return util.RS{Message: "取消失敗，店家已接單，若欲取消請重新操作", Status: false}
	} else if mapResult := orderTmpRepository.GetOrderTempPayment(takeoutOrder.OrderId); mapResult["digitalPayment"].(bool) {
		return util.RS{Message: "有使用線上金流付款", Status: true, Data: mapResult}
	} else {
		return util.RS{Message: "", Status: false}
	}
}
