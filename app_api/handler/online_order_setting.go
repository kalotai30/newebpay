package handler

import (
	"app_api/content"
	"app_api/model"
	"app_api/repository"
	"app_api/util"
	"app_api/util/log"
	"encoding/json"
	"errors"
)

type OnlineOrderSettingHandler content.Handler

type OnlineOrderSettingReceived struct {
	StoreId int64 `json:"storeId"`
}

func (handler *OnlineOrderSettingHandler) GetOnlineOrderAndOtherSetting() interface{} {
	var (
		data                                OnlineOrderSettingReceived
		storeRepository                     repository.StoreRepository
		onlineOrderDeliveryConditionSetting model.OnlineOrderDeliveryConditionSetting
		posPaymentMethod                    model.PosPaymentMethod
		userAddress                         model.UserAddress
		onlineOrderDeliveryTimeSetting      model.OnlineOrderDeliveryTimeSetting
		userPayToken                        model.UserPayToken
	)

	if err := json.Unmarshal([]byte(handler.Parameter), &data); err != nil {
		log.Error(err)
		return util.RS{Message: "", Status: false}
	} else if data.StoreId == 0 {
		log.Error(errors.New("data.StoreId == 0"))
		return util.RS{Message: "", Status: false}
	} else {

		storeInfo := storeRepository.SetId(data.StoreId).GetStoreOrderSetting()

		deliveryConditionList := make([]map[string]interface{}, 0)
		maxDistance := 0
		onlineOrderDeliveryConditionSetting.
			SetStoreId(data.StoreId).
			SetCheckStatus(true).
			SetStatus(true).
			SetCheckDeleteStatus(true).
			SetDeleteStatus(false).
			QueryAll(func(rs *model.OnlineOrderDeliveryConditionSetting) {

				if maxDistance < rs.Distance {
					maxDistance = rs.Distance
				}

				deliveryConditionList = append(deliveryConditionList, map[string]interface{}{
					"id":             rs.Id,
					"status":         rs.Status,
					"distance":       rs.Distance,
					"requiredAmount": rs.RequiredAmount,
					"deliveryFee":    rs.DeliveryFee,
				})
			})

		paymentMethod := map[string]interface{}{
			"creditCard":   false,
			"userPayToken": false,
		}

		posPaymentMethod.
			SetStoreId(data.StoreId).
			SetType(3).
			SetCheckStatus(true).
			SetStatus(true).
			QueryOne()

		userPayToken.
			SetUserId(handler.UserId).
			SetTokenType(1).
			QueryOneByTokeLife()

		if posPaymentMethod.Id > 0 {
			if storeInfo["storeShopCode"].(string) != "" {
				paymentMethod["creditCard"] = true
			}
			if userPayToken.Id > 0 {
				paymentMethod["userPayToken"] = true
			}
		}

		if handler.UserId > 0 {
			userAddress.
				SetUserId(handler.UserId).
				SetTargetLat(storeInfo["storeLat"].(float64)).
				SetTargetLng(storeInfo["storeLng"].(float64)).
				SetCheckIsDefault(true).
				SetIsDefault(true).
				SetCheckDeleteStatus(true).
				SetDeleteStatus(false).
				QueryOne()
		}

		onlineOrderDeliveryTimeSetting.
			SetStoreId(data.StoreId).
			SetDistance(userAddress.Distance).
			SetCheckDeleteStatus(true).
			SetDeleteStatus(false).
			QueryOne()

		userAddressData := map[string]interface{}{
			"id":         userAddress.Id,
			"name":       userAddress.Name,
			"title":      userAddress.Title,
			"city":       userAddress.City,
			"township":   userAddress.Township,
			"address":    userAddress.Address,
			"postalCode": userAddress.PostalCode,
			"content":    userAddress.Content,
			"isDefault":  userAddress.IsDefault,
			"lat":        userAddress.Lat,
			"lng":        userAddress.Lng,
			"distance":   userAddress.Distance,
			"time":       onlineOrderDeliveryTimeSetting.Minute,
		}

		return util.RS{
			Message: "",
			Data: map[string]interface{}{
				"storeInfo":             storeInfo,
				"deliveryConditionList": deliveryConditionList,
				"maxDistance":           maxDistance,
				"paymentMethod":         paymentMethod,
				"userAddressData":       userAddressData,
			},
			Status: true,
		}
	}
}
