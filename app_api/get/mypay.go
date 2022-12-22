package get

import (
	"app_api/config"
	"app_api/handler"
	"app_api/util"
	"app_api/util/crypto"
	"app_api/util/log"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func MyPay(c *gin.Context) {

	plaintext, err := json.Marshal(map[string]string{
		"store_uid": "",
		"pfn":       "1",
	})

	log.Error(err)

	encryptedData, _ := crypto.AesCBC256Encrpty("", plaintext)

	c.HTML(http.StatusOK, `mypay.tmpl`, gin.H{
		"encryptedData": encryptedData,
		"appApiHost":    config.ServerInfo.PaymentHost + "api",
	})
}

func Payment(c *gin.Context) {
	var (
		orderHandler handler.OrderHandler
	)

	orderTmpId, _ := strconv.Atoi(c.Query("orderTmpId"))
	status, _ := strconv.Atoi(c.Query("status"))

	if status == 1 {
		utilRsValu := orderHandler.UseGoldFlowCreateOrder(int64(orderTmpId)).(util.RS)

		if utilRsValu.Status {
			orderData := utilRsValu.Data.(map[string]interface{})
			orderId := orderData["orderId"].(float64)

			c.HTML(http.StatusOK, `payment.tmpl`, gin.H{
				"orderId": orderId,
			})
		}
	} else {
		//刷卡失敗的處理方式
		fmt.Println("orderTmpId:" + strconv.Itoa(orderTmpId))
		fmt.Println("status:" + strconv.Itoa(status))
	}

}
