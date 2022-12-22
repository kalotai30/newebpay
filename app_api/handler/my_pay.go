package handler

import (
	"app_api/config"
	"app_api/content"
	"app_api/model"
	"app_api/util"
	"app_api/util/crypto"
	"app_api/util/log"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type MyPayHandler content.Handler

type MyPayReceived struct {
	TradeToken string `json:"tradeToken"`
	OrderTmpId int    `json:"orderTmpId"`
}

type MyPayEPaymentResult struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Url  string `json:"url"`
}

func (handler *MyPayHandler) MyPayment() interface{} {
	var (
		data        MyPayReceived
		orderTmp    model.OrderTmp
		defaultTime time.Time
	)

	if err := json.Unmarshal([]byte(handler.Parameter), &data); err != nil {
		return util.RS{Message: "", Status: false}
	} else {
		agentUid := ""
		storeUid := ""

		//抓暫存訂單資料
		orderTmp.SetId(int64(data.OrderTmpId)).QueryOne()

		if orderTmp.JsonContent == "" {
			return util.RS{Message: "查無訂單資料", Status: false}
		} else if orderTmp.PaymentAt != defaultTime {
			return util.RS{Message: "已付款，無法重複付款", Status: false}
		}

		var mapResult map[string]interface{}
		if err := json.Unmarshal([]byte(orderTmp.JsonContent), &mapResult); err != nil {
			fmt.Println("JsonToMapDemo err: ", err)
		}

		userCellPhone, _ := strconv.ParseInt(mapResult["takeoutInfo"].(map[string]interface{})["clientPhone"].(string), 10, 16)
		userData := map[string]interface{}{
			"user_id":        fmt.Sprint(mapResult["userId"].(float64)),
			"ip":             "127.0.0.1",
			"user_name":      mapResult["takeoutInfo"].(map[string]interface{})["clientName"].(string),
			"user_real_name": mapResult["takeoutInfo"].(map[string]interface{})["clientName"].(string),
			"user_address":   "台南市",
			"user_cellphone": userCellPhone,
			"user_email":     "test@test.com",
		}

		mealChoiceListData := mapResult["mealChoiceList"].([]interface{})
		items := make([]map[string]interface{}, 0)
		for _, m := range mealChoiceListData {
			mealChoice := m.(map[string]interface{})
			item := map[string]interface{}{
				"id":     fmt.Sprint(mealChoice["id"].(float64)),
				"name":   mealChoice["name"].(string),
				"cost":   fmt.Sprint(mealChoice["choicePrice"].(float64)),
				"amount": fmt.Sprint(mealChoice["amount"].(float64)),
				"total":  fmt.Sprint(mealChoice["choicePrice"].(float64) * mealChoice["amount"].(float64)),
			}
			items = append(items, item)
		}

		successReturl := config.ServerInfo.PaymentHost + "payment?status=1&orderTmpId=" + strconv.Itoa(data.OrderTmpId)
		failureReturl := config.ServerInfo.PaymentHost + "payment?status=0&orderTmpId=" + strconv.Itoa(data.OrderTmpId)
		payment := map[string]interface{}{
			"store_uid":      storeUid,
			"items":          items,
			"cost":           int(mapResult["totalPrice"].(float64)),
			"currency":       "TWD",
			"order_id":       strconv.Itoa(data.OrderTmpId),
			"user_data":      userData,
			"success_returl": successReturl,
			"failure_returl": failureReturl,
			"trade_token":    data.TradeToken,
		}

		paymentJsonText, err := json.Marshal(payment)
		log.Error(err)

		encryData, err := crypto.AesCBC256Encrpty("grgeoEhGH0PnImEEb9YLGuggD7nyNZAm", paymentJsonText)
		log.Error(err)

		serviceData := map[string]interface{}{
			"service_name": "api",
			"cmd":          "api/iaptransaction",
		}

		serviceJsonText, err := json.Marshal(serviceData)
		log.Error(err)

		service, err := crypto.AesCBC256Encrpty("grgeoEhGH0PnImEEb9YLGuggD7nyNZAm", serviceJsonText)
		log.Error(err)

		postData := map[string]interface{}{
			"agent_uid":  agentUid,
			"service":    service,
			"encry_data": encryData,
		}

		fmt.Println(postData)

		body, err := json.Marshal(postData)

		log.Error(err)

		formData := url.Values{
			"agent_uid":  {agentUid},
			"service":    {service},
			"encry_data": {encryData},
		}

		res, err := http.NewRequest("POST",
			"https://pay.usecase.cc/api/agent",
			strings.NewReader(formData.Encode()))

		log.Error(err)

		res.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		res.Header.Set("Accept-Charset", "UTF-8")

		client := &http.Client{}

		resp, err := client.Do(res)
		log.Error(err)

		defer closeRespBody(resp.Body)

		body, err = ioutil.ReadAll(resp.Body)
		log.Error(err)

		fmt.Println("response Body:", string(body))
		return util.RS{}
	}
}

func closeRespBody(body io.ReadCloser) {
	fmt.Println(body.Close())
}

func (handler *MyPayHandler) OCPEPayment() interface{} {
	/*
		可全移除掉

		var (
			responseData OcpResultFormat
		)

		//高鉅OCP電子錢包
		if paymentJsonText, err := json.Marshal(map[string]interface{}{ //指定交易資訊內容並轉成json字串
			"user_id":    strconv.FormatInt(handler.UserId, 10),
			"type":       "bind",
			"return_url": "",
		}); err != nil {
			log.Error(err) //"建立傳輸資料錯誤"
		} else if encryData, err := crypto.AesCBC256Encrpty(config.MyPayInfo.MypayKey, paymentJsonText); err != nil { //將交易資訊加密
			log.Error(err) //"將交易資訊加密錯誤"
		} else if serviceJsonText, err := json.Marshal(map[string]interface{}{ //服務資訊轉乘json字串
			"service_name": "ocpap",
			"cmd":          "api/ewallet",
		}); err != nil {
			log.Error(err) //"建立服務資訊錯誤"
		} else if service, err := crypto.AesCBC256Encrpty(config.MyPayInfo.MypayKey, serviceJsonText); err != nil { //服務資訊加密
			log.Error(err) //"服務資訊加密錯誤"
		} else if response, err := rpc.ApiCall(config.MyPayInfo.MypayInnerPayUrl, "POST", map[string]interface{}{ //呼叫金流api
			"agent_uid":  config.MyPayInfo.AgentUid,
			"service":    service,
			"encry_data": encryData,
		}); err != nil {
			log.Error(err) //"呼叫金流api錯誤"
		} else if err := json.Unmarshal([]byte(response), &responseData); err != nil { //將回傳json字串轉成物件
			log.Error(err) //"將回傳json字串轉成物件錯誤"
		}

		data := map[string]interface{}{
			"userId": strconv.FormatInt(handler.UserId, 10),
			"code":   responseData.Code,
			"msg":    responseData.Msg,
			"url":    responseData.Url,
		}
	*/
	return util.RS{Message: "OCP電子錢包綁定", Status: false}
}
