package digitalpay

import (
	"app_api/config"
	"app_api/util/crypto"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/*
 * NotifyURL 設定 管理後台 補觸發按鈕的網址，沒設定就不會出現補觸發
 * ReturnURL 設定 當下交易完成後，回傳資料的網址。非3D交易無需設定
 */

func NdnpaB10(tradeInfoMap map[string]interface{}, merchantID string, hashKey string, hashIV string) map[string]interface{} {
	tradeInfoStr := crypto.HttpBuildQuery(tradeInfoMap)
	encryData, _ := crypto.Aes256CBCEncryptTypeI(tradeInfoStr, hashKey, hashIV)

	/*
	 * 呼叫藍新金流 api
	 * 回傳格式需設定為 String，接收的資料才會有 key
	 */
	resp, err := http.PostForm(config.NewebPayInfo.APIUrl,
		url.Values{"MerchantID_": {merchantID},
			"PostData_": {encryData},
			"Pos_":      {"String"}})
	if err != nil {
		errRespBodyMap := map[string]interface{}{
			"Status":  "False",
			"Message": "藍新金流_幕後_取回回傳值失敗",
		}

		return errRespBodyMap
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errRespBodyMap := map[string]interface{}{
			"Status":  "False",
			"Message": "藍新金流_幕後_讀取body失敗",
		}

		return errRespBodyMap
	}

	// 解析 非3D 回覆值
	respBody := string(body)
	respBodyMap := crypto.HttpParseQuery(respBody)
	checkCodeMap := map[string]interface{}{
		"Amt":             respBodyMap["Amt"],
		"MerchantID":      respBodyMap["MerchantID"],
		"MerchantOrderNo": respBodyMap["MerchantOrderNo"],
		"TradeNo":         respBodyMap["TradeNo"],
	}
	checkCodeBool := crypto.CheckcodeChcek(checkCodeMap, fmt.Sprint(respBodyMap["CheckCode"]), hashKey, hashIV)
	if respBodyMap["Status"] == "SUCCESS" && checkCodeBool {
		return respBodyMap
	} else {
		errRespBodyMap := map[string]interface{}{
			"Status":  respBodyMap["Status"],
			"Message": respBodyMap["Message"],
		}

		return errRespBodyMap
	}
}

func PAgreementAuthorize(data map[string]interface{}) map[string]interface{} {

	//排列參數並串聯後，進行資料加密
	tradeInfoBasicMap := map[string]interface{}{
		"MerchantID":      config.NewebPayInfo.MerchantID_MS,
		"RespondType":     "String",
		"TimeStamp":       strconv.FormatInt(time.Now().Unix(), 10),
		"Version":         "1.7",
		"Amt":             1,
		"ItemDesc":        "綁定授權",
		"NotifyURL":       config.NewebPayInfo.NotifyURL + "/pagreementNotify",
		"LoginType":       0,
		"CREDITAGREEMENT": 1,
		"OrderComment":    "約定事項",
	}
	tradeInfoMap := MergeMaps(tradeInfoBasicMap, data)

	tradeInfoStr := crypto.HttpBuildQuery(tradeInfoMap)
	encryData, _ := crypto.Aes256CBCEncryptTypeI(tradeInfoStr, config.NewebPayInfo.HashKey_MS, config.NewebPayInfo.HashIV_MS)

	//加密字串前後加上商店專屬的HashKey及HashIV
	strSlices := []string{
		"HashKey=" + config.NewebPayInfo.HashKey_MS,
		encryData,
		"HashIV=" + config.NewebPayInfo.HashIV_MS,
	}
	paramsStr := strings.Join(strSlices, "&")

	//使用SHA256壓碼過後並轉大寫
	hash := sha256.Sum256([]byte(paramsStr))
	hashStr := strings.ToUpper(fmt.Sprintf("%x", hash))

	returnData := map[string]interface{}{
		"PostForm":   config.NewebPayInfo.MPGUrl,
		"MerchantID": config.NewebPayInfo.MerchantID_MS,
		"TradeInfo":  encryData,
		"TradeSha":   hashStr,
		"Version":    "1.7",
	}

	return returnData
}

func NDNF101(data map[string]interface{}) map[string]interface{} {
	//排列參數並串聯後，進行資料加密
	tradeInfoBasicMap := map[string]interface{}{
		"MerchantID":   config.NewebPayInfo.MerchantID_MS,
		"RespondType":  "String",
		"TimeStamp":    strconv.FormatInt(time.Now().Unix(), 10),
		"Version":      "2.0",
		"NotifyURL":    config.NewebPayInfo.NotifyURL + "/ndnf101Notify",
		"LoginType":    0,
		"OrderComment": "約定事項",
	}
	tradeInfoMap := MergeMaps(tradeInfoBasicMap, data)

	tradeInfoStr := crypto.HttpBuildQuery(tradeInfoMap)
	encryData, _ := crypto.Aes256CBCEncryptTypeI(tradeInfoStr, config.NewebPayInfo.HashKey_MS, config.NewebPayInfo.HashIV_MS)

	//加密字串前後加上商店專屬的HashKey及HashIV
	strSlices := []string{
		"HashKey=" + config.NewebPayInfo.HashKey_MS,
		encryData,
		"HashIV=" + config.NewebPayInfo.HashIV_MS,
	}
	paramsStr := strings.Join(strSlices, "&")

	//使用SHA256壓碼過後並轉大寫
	hash := sha256.Sum256([]byte(paramsStr))
	hashStr := strings.ToUpper(fmt.Sprintf("%x", hash))

	returnData := map[string]interface{}{
		"PostForm":   config.NewebPayInfo.MPGUrl,
		"MerchantID": config.NewebPayInfo.MerchantID_MS,
		"TradeInfo":  encryData,
		"TradeSha":   hashStr,
		"Version":    "2.0",
	}

	return returnData
}

func NDNF101NPAB032(data map[string]interface{}) map[string]interface{} {
	//排列參數並串聯後，進行資料加密
	tradeInfoBasicMap := map[string]interface{}{
		"RespondType": "String",
		"Version":     "1.1",
		"TimeStamp":   strconv.FormatInt(time.Now().Unix(), 10),
		"IndexType":   2,
		"CloseType":   2,
	}
	tradeInfoMap := MergeMaps(tradeInfoBasicMap, data)

	tradeInfoStr := crypto.HttpBuildQuery(tradeInfoMap)
	encryData, _ := crypto.Aes256CBCEncryptTypeI(tradeInfoStr, config.NewebPayInfo.HashKey_MS, config.NewebPayInfo.HashIV_MS)

	returnData := map[string]interface{}{
		"PostForm":   config.NewebPayInfo.APICLoseUrl,
		"MerchantID": config.NewebPayInfo.MerchantID_MS,
		"PostData":   encryData,
	}

	return returnData
}

func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	/* 合併多個map */
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
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
