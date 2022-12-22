package repository

import (
	"encoding/json"
	"fmt"
	"time"
	"util_api/model"
	"util_api/util/kafka"
	"util_api/util/log"
)

type OrderTmpRepository struct{}

func (repository *OrderTmpRepository) GetOrderTempData(orderTmpId int64) map[string]interface{} {
	var (
		orderTmp    model.OrderTmp
		defaultTime time.Time
	)

	// 抓暫存訂單資料
	orderTmp.SetId(orderTmpId).QueryOne()

	if orderTmp.JsonContent == "" {
		errRespMap := map[string]interface{}{
			"Status":  "False",
			"Message": "查無訂單資料",
		}

		return errRespMap
	} else if orderTmp.PaymentAt != defaultTime {
		errRespMap := map[string]interface{}{
			"Status":  "False",
			"Message": "已付款，無法重複付款",
		}

		return errRespMap
	}

	// 解析暫存訂單內容(原內容是用JSON存在)
	var mapResult map[string]interface{}
	if err := json.Unmarshal([]byte(orderTmp.JsonContent), &mapResult); err != nil {
		fmt.Println("JsonToMapDemo err: ", err)
	}

	return mapResult
}

func (repository *OrderTmpRepository) SaveToOrder(mapResult map[string]interface{}, orderTmpId int64) string {
	/* 完成暫存訂單資料寫入訂單程序 */
	var orderTmp model.OrderTmp

	_, err := kafka.Push("OnlineCreateOrder2", mapResult)
	log.Error(err)

	dataList := make(map[string]interface{}, 0)
	for {
		if kafKaRS := <-kafka.KafkaChan["OnlineAppApiCreateOrder"]; kafKaRS.Uuid == mapResult["uuid"].(string) {
			dataList["status"] = kafKaRS.Status
			dataList["message"] = kafKaRS.Message
			dataList["data"] = kafKaRS.Data
			break
		}
	}

	if dataList["status"] == true {
		orderDataMap := dataList["data"].(map[string]interface{})

		orderTmp.
			SetId(orderTmpId).
			SetPaymentAt(time.Now()).
			SetOrderId(int64(orderDataMap["orderId"].(float64))).
			Update([]string{"payment_at", "order_id"})

		orderSn := orderDataMap["orderSn"].(string)

		return orderSn
	} else {
		return "false"
	}
}
