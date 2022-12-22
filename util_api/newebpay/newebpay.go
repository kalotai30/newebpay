package newebpay

import (
	"encoding/json"
	"fmt"
	"time"
	"util_api/config"
	"util_api/content"
	"util_api/model"
	"util_api/repository"
	"util_api/util"
	"util_api/util/crypto"
	"util_api/util/log"
	"util_api/util/rpc"

	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type NewebPay content.Handler

type receivedRSData struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Data    struct {
		MerchantID string `json:"MerchantID"`
		PostForm   string `json:"PostForm"`
		TradeInfo  string `json:"TradeInfo"`
		TradeSha   string `json:"TradeSha"`
		Version    string `json:"Version"`
	} `json:"data"`
}

func PAgreementNotifyForPOST(c *gin.Context) {
	var (
		userPayToken             model.UserPayToken
		newebpayPagreementReturn model.NewebpayPagreementReturn
	)

	// 判斷是否有正確接收到回傳資料
	if err := c.Request.ParseForm(); err != nil {
		log.Error(err)
		return
	}

	//將接收到的資料轉換成Map，值強轉為字串
	resultMap := make(map[string]interface{}, 0)
	for key, values := range c.Request.Form {
		resultMap[key] = strings.Join(values, "")
	}

	if resultMap["Status"].(string) == "SUCCESS" {
		// 解密 TradeInfo 交易資料
		decryData, _ := crypto.Aes256CBCDecryptTypeI(resultMap["TradeInfo"].(string), config.NewebPayInfo.HashKey_MS, config.NewebPayInfo.HashIV_MS)
		TradeInfoMap := crypto.HttpParseQuery(decryData)

		// 拆解所需資料
		SetUserId := TradeInfoMap["MerchantOrderNo"].(string)[strings.Index(TradeInfoMap["MerchantOrderNo"].(string), "P")+1 : strings.Index(TradeInfoMap["MerchantOrderNo"].(string), "T")]
		SetUserID64, _ := strconv.ParseInt(SetUserId, 10, 64)

		userPayToken.SetUserId(SetUserID64).SetTokenType(1).QueryOne()
		local, _ := time.LoadLocation("Asia/Taipei") //修改成台北時間
		setTokenLifeDate, _ := time.ParseInLocation("2006-01-02", TradeInfoMap["TokenLife"].(string), local)
		tokenUpdatedAtDatetime, _ := time.ParseInLocation("2006-01-02 15:04:05", TradeInfoMap["PayTime"].(string), local)

		// 更新資料表user_pay_token
		if userPayToken.Token != "" {
			err := userPayToken.
				SetUserId(SetUserID64).
				SetToken(TradeInfoMap["TokenValue"].(string)).
				SetTokenLife(setTokenLifeDate).
				SetTokenUpdatedAt(tokenUpdatedAtDatetime).
				Update([]string{"token", "tokenLife", "token_updated_at"})
			if err != nil {
				log.Error(err)
				return
			}
		} else {
			id, _ := userPayToken.
				SetUserId(SetUserID64).
				SetToken(TradeInfoMap["TokenValue"].(string)).
				SetTokenLife(setTokenLifeDate).
				SetTokenUpdatedAt(tokenUpdatedAtDatetime).
				SetTokenType(1).
				Insert()

			fmt.Println(id)
		}

		// 儲存藍新回覆資料
		amt64, _ := strconv.ParseInt(TradeInfoMap["Amt"].(string), 10, 64)
		instFirst64, _ := strconv.ParseInt(TradeInfoMap["InstFirst"].(string), 10, 64)
		inst64, _ := strconv.ParseInt(TradeInfoMap["Inst"].(string), 10, 64)
		instEach64, _ := strconv.ParseInt(TradeInfoMap["InstEach"].(string), 10, 64)
		tokenUseStatus64, _ := strconv.ParseInt(TradeInfoMap["TokenUseStatus"].(string), 10, 64)
		var itemDescStr string
		if TradeInfoMap["ItemDesc"] == nil {
			itemDescStr = ""
		} else {
			itemDescStr = TradeInfoMap["ItemDesc"].(string)
		}
		id, _ := newebpayPagreementReturn.
			SetUserId(SetUserID64).
			SetMerchantID(TradeInfoMap["MerchantID"].(string)).
			SetAmt(amt64).
			SetTradeNo(TradeInfoMap["TradeNo"].(string)).
			SetMerchantOrderNo(TradeInfoMap["MerchantOrderNo"].(string)).
			SetPaymentType(TradeInfoMap["PaymentType"].(string)).
			SetRespondType(TradeInfoMap["RespondType"].(string)).
			SetPayTime(tokenUpdatedAtDatetime).
			SetIP(TradeInfoMap["IP"].(string)).
			SetEscrowBank(TradeInfoMap["EscrowBank"].(string)).
			SetRespondCode(TradeInfoMap["RespondCode"].(string)).
			SetAuth(TradeInfoMap["Auth"].(string)).
			SetCard6No(TradeInfoMap["Card6No"].(string)).
			SetCard4No(TradeInfoMap["Card4No"].(string)).
			SetExp(TradeInfoMap["Exp"].(string)).
			SetECI(TradeInfoMap["ECI"].(string)).
			SetInst(inst64).
			SetInstFirst(instFirst64).
			SetInstEach(instEach64).
			SetTokenUseStatus(tokenUseStatus64).
			SetTokenValue(TradeInfoMap["TokenValue"].(string)).
			SetTokenLife(setTokenLifeDate).
			SetPaymentMethod(TradeInfoMap["PaymentMethod"].(string)).
			SetItemDesc(itemDescStr).
			Insert()

		fmt.Println(id)
	} else {
		fmt.Println(resultMap["Message"])
		return
	}
}

func Ndnf101NotifyForPOST(c *gin.Context) {
	// 判斷是否有接收到回傳資料
	if err := c.Request.ParseForm(); err != nil {
		log.Error(err)
		return
	}

	var (
		orderTmpRepository       repository.OrderTmpRepository
		newebpayNdnfCommonReturn model.NewebpayNdnfCommonReturn
		newebpayNdnfCreditReturn model.NewebpayNdnfCreditReturn
	)

	//將接收到的資料轉換成Map，值強轉為字串
	resultMap := make(map[string]interface{}, 0)
	for key, values := range c.Request.Form {
		resultMap[key] = strings.Join(values, "")
	}

	if resultMap["Status"] == "SUCCESS" {
		// 解密 TradeInfo 交易資料
		decryData, _ := crypto.Aes256CBCDecryptTypeI(resultMap["TradeInfo"].(string), config.NewebPayInfo.HashKey_MS, config.NewebPayInfo.HashIV_MS)
		TradeInfoMap := crypto.HttpParseQuery(decryData)

		// 抓暫存訂單資料
		mapResult := orderTmpRepository.GetOrderTempData(ReturnValueInt64(TradeInfoMap["MerchantOrderNo"]))

		// 儲存藍新回覆資料(newebpay_ndnf_common_return)
		ndnfId, _ := newebpayNdnfCommonReturn.
			SetUserId(int64(mapResult["userId"].(float64))).
			SetMerchantID(ReturnValueStr(TradeInfoMap["MerchantID"])).
			SetAmt(ReturnValueInt64(TradeInfoMap["Amt"])).
			SetTradeNo(ReturnValueStr(TradeInfoMap["TradeNo"])).
			SetMerchantOrderNo(ReturnValueStr(TradeInfoMap["MerchantOrderNo"])).
			SetPaymentType(ReturnValueStr(TradeInfoMap["PaymentType"])).
			SetRespondType(ReturnValueStr(TradeInfoMap["RespondType"])).
			SetPayTime(ReturnValueTime(TradeInfoMap["PayTime"])).
			SetIP(ReturnValueStr(TradeInfoMap["IP"])).
			SetEscrowBank(ReturnValueStr(TradeInfoMap["EscrowBank"])).
			Insert()

		// 儲存藍新回覆資料(newebpay_ndnf_credit_return)
		ndnfCreditId, _ := newebpayNdnfCreditReturn.
			SetNdnfId(ndnfId).
			SetAuthBank(ReturnValueStr(TradeInfoMap["AuthBank"])).
			SetRespondCode(ReturnValueStr(TradeInfoMap["RespondCode"])).
			SetAuth(ReturnValueStr(TradeInfoMap["Auth"])).
			SetCard6No(ReturnValueStr(TradeInfoMap["Card6No"])).
			SetCard4No(ReturnValueStr(TradeInfoMap["Card4No"])).
			SetExp(ReturnValueStr(TradeInfoMap["Exp"])).
			SetInst(ReturnValueInt64(TradeInfoMap["Inst"])).
			SetInstFirst(ReturnValueInt64(TradeInfoMap["InstFirst"])).
			SetInstEach(ReturnValueInt64(TradeInfoMap["InstEach"])).
			SetECI(ReturnValueStr(TradeInfoMap["ECI"])).
			SetTokenUseStatus(ReturnValueInt64(TradeInfoMap["TokenUseStatus"])).
			SetRedAmt(ReturnValueInt64(TradeInfoMap["RedAmt"])).
			SetPaymentMethod(ReturnValueStr(TradeInfoMap["PaymentMethod"])).
			Insert()

		/* 完成寫入訂單程序 */
		orderSn := orderTmpRepository.SaveToOrder(mapResult, ReturnValueInt64(TradeInfoMap["MerchantOrderNo"]))
		if ndnfId >= 1 && ndnfCreditId >= 1 && orderSn != "false" {
			fmt.Println(orderSn)
		}

		return
	} else {
		fmt.Println(resultMap["Message"])
		return
	}
}

func NewebpayPost(c *gin.Context) {
	var (
		utilRsValu  util.RS
		returnData  map[string]interface{}
		pageContent string
	)

	authToken := c.Query("token")
	paymentType, _ := strconv.ParseInt(c.Query("pay"), 10, 64)

	if paymentType == 2 {
		//藍新信用卡頁面
		queryDataJsonString, _ := crypto.KeyDecrypt(authToken)

		if response, err := rpc.ApiCall(config.ServerInfo.OinApiHost, "POST", map[string]interface{}{
			"action":     "PayByCreditCard",
			"parameters": queryDataJsonString,
		}); err != nil {
			log.Error(err)
			return
		} else if err = json.Unmarshal([]byte(response), &utilRsValu); err != nil {
			log.Error(err)
			return
		}

		returnData = utilRsValu.Data.(map[string]interface{})
		pageContent = "線上金流頁面轉向中..."
	} else if paymentType == 9 {
		//藍新信用卡分期頁面
		queryDataJsonString, _ := crypto.KeyDecrypt(authToken)

		if response, err := rpc.ApiCall(config.ServerInfo.OinApiHost, "POST", map[string]interface{}{
			"action":     "PayByInstallment",
			"parameters": queryDataJsonString,
		}); err != nil {
			log.Error(err)
			return
		} else if err = json.Unmarshal([]byte(response), &utilRsValu); err != nil {
			log.Error(err)
			return
		}

		returnData = utilRsValu.Data.(map[string]interface{})
		pageContent = "線上金流頁面轉向中..."

	} else {
		//藍新綁卡頁面
		if response, err := rpc.AuthTokeApiCall(config.ServerInfo.OinApiHost, "POST", authToken, map[string]interface{}{
			"action": "AgreedCreditCardAuthorization",
		}); err != nil {
			log.Error(err)
			return
		} else if err = json.Unmarshal([]byte(response), &utilRsValu); err != nil {
			log.Error(err)
			return
		}

		returnData = utilRsValu.Data.(map[string]interface{})
		pageContent = "綁定頁面轉向中..."
	}

	c.HTML(http.StatusOK, `newebpayPost.tmpl`, gin.H{
		"pageContent": pageContent,
		"action":      returnData["PostForm"].(string),
		"merchantID":  returnData["MerchantID"].(string),
		"tradeInfo":   returnData["TradeInfo"].(string),
		"tradeSha":    returnData["TradeSha"].(string),
		"version":     returnData["Version"].(string),
	})
}

func ReturnValueInt64(returnValue interface{}) int64 {
	var result int64

	if returnValue != nil {
		result, _ = strconv.ParseInt(returnValue.(string), 10, 64)
	} else {
		result = 0
	}

	return result
}

func ReturnValueStr(returnValue interface{}) string {
	var result string

	if returnValue != nil {
		result = returnValue.(string)
	} else {
		result = ""
	}

	return result
}

func ReturnValueTime(returnValue interface{}) time.Time {
	var result time.Time
	local, _ := time.LoadLocation("Asia/Taipei") //修改成台北時間

	if returnValue != nil {
		result, _ = time.ParseInLocation("2006-01-02 15:04:05", returnValue.(string), local)
	} else {
		result, _ = time.ParseInLocation("2006-01-02 15:04:05", "1000-01-01 00:00:00", local)
	}

	return result
}
