package handler

import (
	"app_api/collections"
	"app_api/config"
	"app_api/content"
	"app_api/database/mysql"
	"app_api/model"
	"app_api/repository"
	"app_api/util"
	"app_api/util/crypto"
	digitalpay "app_api/util/digitalPay"
	"app_api/util/fcm"
	"app_api/util/hit"
	"app_api/util/image"
	"app_api/util/kafka"
	"app_api/util/log"
	"app_api/util/rpc"
	"app_api/util/validate"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type OrderHandler content.Handler

type ReceivedOrder struct {
	OrderId             int64              `json:"orderId"`
	TakeoutId           int64              `json:"takeoutId"`
	KeyWord             string             `json:"keyWord"`
	OrderSn             string             `json:"orderSn"`
	StartAt             string             `json:"startAt"`
	EndAt               string             `json:"endAt"`
	OffSet              int                `json:"OffSet"`
	LimitOrder          bool               `json:"limitOrder"`
	TotalPrice          float64            `json:"totalPrice"` //ๆไป้้ก
	StoreId             int64              `json:"storeId"`
	MealChoiceList      []MealChoiceList   `json:"mealChoiceList"`
	MealData            []MealData         `json:"mealData"`
	Remark              string             `json:"remark"`
	MealPrice           int64              `json:"mealPrice"`
	SumAmount           int64              `json:"sumAmount"`
	CouponChoiceList    []CouponChoiceList `json:"couponChoiceList"`
	OrderGift           []interface{}      `json:"OrderGift"`
	TaxIdNumber         string             `json:"taxIdNumber"`
	UserAddressId       int64              `json:"userAddressId"`
	DeliveryFee         int                `json:"deliveryFee"`
	RadioTake           int                `json:"radioTake"`
	IsAccepted          bool               `json:"isAccepted"`
	CancelReason        string             `json:"cancelReason"`
	OrderFoodList       []OrderFoodList    `json:"orderFoodList"`
	IsRequiredTableware bool               `json:"isRequiredTableware"`
	PaymentType         int                `json:"paymentType"`
	MealApiFormat       []MealApiFormat    `json:"mealApiFormat"`
	OrderTmpId          int64              `json:"orderTmpId"`
	ReservationTime     string             `json:"reservationTime"`
	PageSize            int                `json:"pageSize"`
	IsCordova           bool               `json:"isCordova"`
	Phone               string             `json:"phone"`
	TableId             int64              `json:"tableId"`
	Adult               int                `json:"adult"`
	Child               int                `json:"child"`
}

type MealApiFormat struct {
	Id                     int64                    `json:"id"`
	Name                   string                   `json:"name"`
	OrderType              int                      `json:"orderType"`
	Amount                 int64                    `json:"amount"`
	ChoicePrice            int64                    `json:"choicePrice"`
	TextNote               string                   `json:"textNote"`
	PosFoodSubmealList     []PosFoodSubmealList     `json:"posFoodSubmealList"`
	PosGoodsNoteDetailList []PosGoodsNoteDetailList `json:"posGoodsNoteDetailList"`
}

type PosFoodSubmealList struct {
	PosFoodId              int64                    `json:"posFoodId"`
	PosFoodName            string                   `json:"posFoodName"`
	IncreasePrice          int                      `json:"increasePrice"`
	PosGoodsNoteDetailList []PosGoodsNoteDetailList `json:"posGoodsNoteDetailList"`
}

type PosGoodsNoteDetailList struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	IncreasePrice int    `json:"increasePrice"`
}
type MealChoiceList struct {
	MealId          int64  `json:"meal_id"`
	Amount          int64  `json:"amount"`
	ChoicePrice     int64  `json:"choicePrice"`
	CombinationJson string `json:"combination_json"`
	TextNote        string `json:"text_note"`
}

type MealData struct {
	MealId   int64  `json:"meal_id"`
	DataJson string `json:"data_json"`
}

type CouponChoiceList struct {
	CouponId     int64  `json:"couponId"`
	CouponListId int64  `json:"couponListId"`
	Name         string `json:"name"`
	Code         string `json:"code"`
	CouponTypeId int    `json:"couponTypeId"`
	ChoiceAmount int    `json:"choiceAmount"`
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
}

type OrderFoodList struct {
	FoodId     int64  `json:"foodId"`
	Name       string `json:"name"`
	PosAppName string `json:"posAppName"`
}

type TakeoutOrderInfo struct {
	Id                   int64     `json:"takeoutId"`
	ClientName           string    `json:"clientName"`
	ClientPhone          string    `json:"clientPhone"`
	Gender               int       `json:"gender"`
	CustomerNote         string    `json:"customerNote"`
	BuyerBan             string    `json:"buyerBan"`
	CarrierNumber        string    `json:"carrierNumber"`
	ClientAddress        string    `json:"clientAddress"`
	SelectArea           string    `json:"selectArea"`
	TakeoutTypeId        int       `json:"takeoutTypeId"`
	TakeoutType          string    `json:"takeoutType"`
	EstimatedTime        time.Time `json:"estimatedTime"`
	Status               int64     `json:"status"`
	IsCheckedOut         bool      `json:"isCheckedOut"`
	DeliveryInstructions string    `json:"deliveryInstructions"`
	Type                 int       `json:"type"`
}

type OcpResultFormat struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Url  string `json:"url"`
}

const (
	TakeoutManagementId = 49
)

func (handler *OrderHandler) GetOrders() interface{} {
	var (
		data             ReceivedOrder
		orders           model.Orders
		orderPaymentTerm model.OrderPaymentTerm
		brand            model.Brand
		store            model.Store
		brandId          int64
		storeIdAry       []interface{}
	)

	if err := json.Unmarshal([]byte(handler.Parameter), &data); err != nil {
		return util.RS{Message: "", Status: false}
	}

	////ๅฎๅบ
	//brandId := brand.GetBrandIdByAppId(handler.AppId)
	//storeIdAry := store.GetStoreIDListBy(brandId)

	if handler.AppId <= 1 { //ๅคๆทๆฏๅฆ็บ็ธฝๅนณๅฐ
		brandId = 0 //ๆฏ็ธฝๅนณๅฐ
	} else {
		brandId = brand.GetBrandIdByAppId(handler.AppId)    //้็ธฝๅนณๅฐ๏ผๅฎๅบ๏ผ
		storeIdAry = store.GetStoreIDListByBrandId(brandId) //ๆๅ็ธฝๅบ็ๆๆ่ณๆ ๆ้คๅฎๅบ
	}

	if len(data.StartAt) > 0 {
		orders.SetStartAt(data.StartAt)
		orders.SetEndAt(data.EndAt)
	}

	dataList := make([]interface{}, 0)
	orders.SetUserId(handler.UserId).SetKeyWord(data.KeyWord).SetLimitOrderCount(data.LimitOrder).SetOffset(data.OffSet).SetStoreIdAry(storeIdAry).QueryAll(func(order *model.Orders) {
		orderPaymentTerm.SetOrderId(order.Id).SumPrice()
		data := map[string]interface{}{
			"orderSn":             order.OrderSn,
			"storeName":           order.StoreName,
			"invoicePrice":        order.InvoicePrice,
			"invoiceNumber":       order.InvoiceNumber,
			"createdAt":           order.CreatedAt.Format("2006-01-02 15:04:05"),
			"brandName":           order.BrandName,
			"brandImage":          image.ReturnPhotoPath(orders.BrandImage),
			"isOnlineOrder":       order.TakeoutOrderId != 0 && order.TakeoutOrderFrom == "app",
			"isContractAvailable": order.IsContractAvailable, //ๅ็ๅ็ดๆฏๅฆๆๆ๏ผ
			"isPosReserve":  	   order.PosReserveId > 0,
		}
		dataList = append(dataList, data)
	})

	return util.RS{Message: "", Status: true, Data: dataList}
}

func (handler *OrderHandler) GetOrder() interface{} {
	var (
		data                         ReceivedOrder
		orderRepository              repository.OrderRepository
		orderDetailRepository        repository.OrderDetailRepository
		orderCardRepository          repository.OrderCardRepository
		orderDiscountRepository      repository.OrderDiscountRepository
		orderPaymentTermRepository   repository.OrderPaymentTermRepository
		orderCouponRepository        repository.OrderCouponRepository
		orderOtherDiscountRepository repository.OrderOtherDiscountRepository
		brandId                      int64
		totalDiscount                int64
	)

	if err := json.Unmarshal([]byte(handler.Parameter), &data); err != nil {
		return util.RS{Message: "", Status: false}
	} else if order := orderRepository.QueryOrder(handler.UserId, data.OrderSn); order.Id == 0 {
		return util.RS{Message: "่จๅฎ็ทจ่ไธๅญๅจ", Status: false}
	} else {

		mysql.Model(new(model.Store)).Where("id", "=", order.StoreId).Select([]string{"brand_id"}).Find().Scan(&brandId)
		brandName := new(model.Brand).GetBrandName(brandId)

		goods, free := orderDetailRepository.QueryAll(order.Id)
		discount := orderDiscountRepository.QueryAll(data.OrderSn)
		payment := orderPaymentTermRepository.QueryAll(order.Id)
		gift := orderCouponRepository.QueryAll(order.Id)
		otherDiscount := orderOtherDiscountRepository.QueryAll(order.Id)
		invoiceStatus := new(model.Orders).GetOrderHasCarrierNumber(order.InvoiceNumber)

		loyaltyCardAmountList := new(model.LoyaltyCardDepositLog).UserGetLoyaltyAmountBYOrderId(order.Id)

		var cardList = make([]map[string]interface{}, 0)
		orderCardRepository.SetOrderId(order.Id).QueryAll(func(rs *repository.OrderCardRepository) {
			cardList = append(cardList, map[string]interface{}{
				"cardName":            rs.CardName,
				"cardDiscount":        rs.CardDiscount,
				"cardDiscountContent": rs.CardDiscountContent,
			})

			totalDiscount = rs.CardDiscount + totalDiscount
		})

		for _, val := range discount {
			item := val.(map[string]interface{})
			totalDiscount = hit.If(item["couponType"] == "addPurchase", item["totalPrice"].(int64)*-1, item["totalPrice"].(int64)).(int64) + totalDiscount
		}
		for _, val := range otherDiscount {
			item := val.(map[string]interface{})
			totalDiscount = item["price"].(int64) + totalDiscount
		}

		return map[string]interface{}{
			"message": "",
			"status":  true,
			"data": map[string]interface{}{
				"totalDiscount":         totalDiscount,
				"invoiceStatus":         invoiceStatus,
				"brandName":             brandName,
				"storeName":             order.StoreName,
				"subPrice":              order.SubPrice,
				"totalPrice":            order.TotalPrice,
				"freePrice":             order.FreePrice,                  //ๆๅพ้้ก
				"freeCount":             order.SubPrice - order.FreePrice, //ๅฐ่จๆฃ้คๆๅพ้้ก
				"invoicePrice":          order.InvoicePrice,
				"tip":                   order.Tip,
				"deliveryFee":           order.DeliveryFee,
				"invoiceNumber":         order.InvoiceNumber,
				"date":                  order.CreatedAt.Format("2006-01-02 15:04:05")[:10],
				"time":                  order.CreatedAt.Format("2006-01-02 15:04:05")[11:16],
				"goods":                 goods,
				"free":                  free,
				"discount":              discount, //ๆๆฃ
				"payment":               payment,  //ๆฏไปๆนๅผ
				"gift":                  gift,
				"otherDiscount":         otherDiscount, //ๅถไปๆๆฃ interface
				"loyaltyCardAmountList": loyaltyCardAmountList,
				"cardList":              cardList,
				"isOnlineOrder":         order.TakeoutOrderId != 0 && order.TakeoutOrderFrom == "app",
				"isPosReserve":          order.PosReserveId > 0,
			},
		}
	}
}

//่ด้ๅชๆ?
func (handler *OrderHandler) GetOrderGiftCoupon() interface{} {
	var (
		data ReceivedOrder
		// cardListRepository        repository.CardListRepository
		// rewardPointList           repository.RewardPointListRepository
		// loyaltyCardListRepository repository.LoyaltyCardListRepository
		couponListRepository repository.CouponListRepository
		//storeProjectContractRepository repository.StoreProjectContractRepository
	)

	if err := json.Unmarshal([]byte(handler.Parameter), &data); err != nil {
		// fmt.Println(err)
		return util.RS{Message: "", Status: false}
	} else if data.TotalPrice == 0 {
		return util.RS{Message: "่ซ่ผธๅฅๆไป้้ก", Status: false}
	} else {
		//guide := storeProjectContractRepository.QueryNew(0, handler.StoreId, 1)["guide"].([]string)
		guide := []string{"2", "3", "4", "5", "6"}
		var (
			givenData = make(map[string]interface{})
		)
		//ๅๆguide็1็งป้คไธฆไธ้?่จญ็ฉบ็ๆๅกๅกๆธๅฎ
		givenData["cardList"] = make([]map[string]interface{}, 0)
		// loyaltySettingList, err := loyaltyCardListRepository.GetLoyaltyPublishSetting(handler.StoreId)
		// if err != nil {
		// 	return util.RS{Message: "ๅๅพ้้ปๅก่ณๆๅคฑๆ", Status: false}
		// }
		// givenData["loyaltyList"] = loyaltyCardListRepository.GetGivenLoyaltyList(int64(data.TotalPrice), loyaltySettingList, data.UserId, handler.StoreId)
		for _, i := range guide {
			switch i {
			case "1":
				{
					// cardSettingList, err := cardListRepository.GetCardPublishSetting(handler.StoreId)
					// if err != nil {
					// 	return util.RS{Message: "ๅๅพๅก็่ณๆๅคฑๆ", Status: false}
					// }
					// givenData["cardList"] = cardListRepository.GetGivenCardList(int64(data.TotalPrice), cardSettingList, data.UserId, handler.StoreId)
				}
			case "2":
				{
					//offset
					offsetSettingList, err := couponListRepository.GetCouponPublishSetting(data.StoreId, 2, 1)
					if err != nil {
						return util.RS{Message: "ๅๅพๆๅนๅท่ณๆๅคฑๆ", Status: false}
					}
					givenData["offsetList"] = couponListRepository.GetGivenCouponList(int64(data.TotalPrice), offsetSettingList, handler.UserId, data.StoreId)

				}
			case "3":
				{
					//discount
					discountSettingList, err := couponListRepository.GetCouponPublishSetting(data.StoreId, 4, 2)
					if err != nil {
						return util.RS{Message: "ๅๅพๆๆฃๅท่ณๆๅคฑๆ", Status: false}
					}
					givenData["discountList"] = couponListRepository.GetGivenCouponList(int64(data.TotalPrice), discountSettingList, handler.UserId, data.StoreId)

				}
			case "4":
				{
					//freebie
					freebieSettingList, err := couponListRepository.GetCouponPublishSetting(data.StoreId, 1, 0)
					if err != nil {
						return util.RS{Message: "ๅๅพๅๆๅท่ณๆๅคฑๆ", Status: false}
					}
					givenData["freebieList"] = couponListRepository.GetGivenCouponList(int64(data.TotalPrice), freebieSettingList, handler.UserId, data.StoreId)
				}
			case "5":
				{
					//addPurchase
					addPurchaseSettingList, err := couponListRepository.GetCouponPublishSetting(data.StoreId, 3, 0)
					if err != nil {
						return util.RS{Message: "ๅๅพๅ?ๅนๅท่ณๆๅคฑๆ", Status: false}
					}
					givenData["addPurchaseList"] = couponListRepository.GetGivenCouponList(int64(data.TotalPrice), addPurchaseSettingList, handler.UserId, data.StoreId)

				}
			case "6":
				{
					// rewardPointSettingList, err := rewardPointList.GetRewardPointPublishSetting(handler.StoreId)
					// if err != nil {
					// 	return util.RS{Message: "ๅๅพ็ดๅฉ้ปๆธ่ณๆๅคฑๆ", Status: false}
					// }
					// givenData["rewardList"] = rewardPointList.GetGivenRewardPointList(int64(data.TotalPrice), rewardPointSettingList)

				}

			}
		}

		// fmt.Println("่ด้ๅชๆ?๏ผ", givenData)
		return util.RS{Message: "", Data: givenData, Status: true}
	}

}

func (handler *OrderHandler) ChkOnlineOrder() interface{} {
	type reciveMealPosFoodInfo struct {
		Id                      int64  `json:"id"`
		IsLimit                 bool   `json:"isLimit"`
		LimitAmount             int64  `json:"limitAmount"`
		OnlineOrderDescription  string `json:"onlineOrderDescription"`
		OnlineOrderFoodImage    string `json:"onlineOrderFoodImage"`
		OnlineOrderName         string `json:"onlineOrderName"`
		OnlineOrderPrice        int64  `json:"onlineOrderPrice"`
		OnlineOrderRemarkPrompt string `json:"onlineOrderRemarkPrompt"`
		OrderType               int64  `json:"orderType"`
	}
	type reciveMealData struct {
		DataList    []interface{}         `json:"dataList"`
		PosFoodInfo reciveMealPosFoodInfo `json:"posFoodInfo"`
	}
	var (
		data ReceivedOrder
	)
	paramStr, _ := crypto.KeyDecrypt(handler.Parameter) //่งฃๅฏ
	if err := json.Unmarshal([]byte(paramStr), &data); err != nil {
		return util.RS{Message: "", Status: false, Data: ""}
	}

	return util.RS{Message: "", Status: false, Data: ""}
}

func (handler *OrderHandler) GetOnlineOrderList() interface{} {
	var (
		data                           ReceivedOrder
		dataList                       = make([]map[string]interface{}, 0)
		orderDeliveryTimeList          = make([]map[string]interface{}, 0)
		onlineOrderDeliveryTimeSetting model.OnlineOrderDeliveryTimeSetting
	)
	if err := json.Unmarshal([]byte(handler.Parameter), &data); err != nil {
		log.Error(err)
		return util.RS{Message: "", Status: false}
	}

	orderList := mysql.Model(new(model.Orders))
	if data.StoreId != 0 {
		orderList.Where("orders.store_id", "=", data.StoreId)
	}

	onlineOrderDeliveryTimeSetting.
		SetCheckDeleteStatus(true).
		SetDeleteStatus(false).
		QueryAll(func(rs *model.OnlineOrderDeliveryTimeSetting) {
			orderDeliveryTimeList = append(orderDeliveryTimeList, map[string]interface{}{
				"id":      rs.Id,
				"storeId": rs.StoreId,
				"kmStart": rs.KmStart,
				"kmEnd":   rs.KmEnd,
				"minute":  rs.Minute,
			})
		})

	if data.PageSize != 0 {
		orderList.Limit(data.OffSet*data.PageSize, data.PageSize)
	}

	orderList.
		Select([]string{
			"orders.id", "orders.store_id", "IfNULL(orders.invoice_price,0)", "IfNULL(orders.sub_price,0)", "IfNULL(orders.created_at,'0000-00-00 00:00:00')", "order_sn", "(select sum(amount) from order_detail where order_detail.order_id = orders.id)",
			"takeout_order.id", "takeout_order.order_from", "takeout_order.order_type", "takeout_order.estimated_time", "takeout_order.cancel_reason", "takeout_order.status", "takeout_order.delivery_distance", "takeout_order.is_user_cancel",
			"store.name", "store.city", "store.township", "store.address", "brand.name", "brand.image", "takeout_order.type",
			collections.SQL_IsContractAvailable, //ๅคๆทๅ็ๅ็ดๆๆ
		}).
		Where("orders.user_id", "=", handler.UserId).
		WhereIn("and", "order_from", []interface{}{"app"}).
		InnerJoin("takeout_order", "takeout_order.order_id", "=", "orders.id").
		InnerJoin("store", "store.id", "=", "orders.store_id").
		InnerJoin("brand", "brand.id", "=", "store.brand_id").
		LeftJoin("brand_contract", "brand_contract.brand_id", "=", "store.brand_id"). //joinๅ็ๅ็ด่ณๆ่กจ
		OrderBy([]string{"takeout_order.created_at"}, []string{"desc"}).
		Get(func(rows *sql.Rows) (isBreak bool) {
			var (
				// subPrice       float64
				orderData                           model.Orders
				orderCreatedAt, estimatedTimeString string
				mealAmount                          int16
				takeoutOrder                        model.TakeoutOrder
				store                               model.Store
				brand                               model.Brand
				isContractAvailable                 int
			)
			log.Error(rows.Scan(
				&orderData.Id, &orderData.StoreId, &orderData.InvoicePrice, &orderData.SubPrice, &orderCreatedAt, &orderData.OrderSn, &mealAmount,
				&takeoutOrder.Id, &takeoutOrder.OrderFrom, &takeoutOrder.OrderType, &takeoutOrder.EstimatedTime, &takeoutOrder.CancelReason, &takeoutOrder.Status, &takeoutOrder.DeliveryDistance, &takeoutOrder.IsUserCancel,
				&store.Name, &store.City, &store.Township, &store.Address, &brand.Name, &brand.Image, &takeoutOrder.Type,
				&isContractAvailable,
			))

			layout := "2006-01-02 15:04:05 -0700 MST"
			// checkTime, err := time.Parse(layout, "2014-11-17 23:02:03 +0800 UTC") // ่ฅๅฐๆชไปๆฌพ createdAt ็ๆ้็บ 0000-00-00 00:00:00 ๆไปฅ็จไปปๆๆ้ๆชขๆฅๆฏๆชขๆฅๆฏๅฆๅคงๆผๅณๅฏ			log.Error(err)
			// log.Error(err)

			if orderCreatedAt != "0000-00-00 00:00:00" && orderCreatedAt != "0001-01-01 00:00:00" {
				orderCreatedAt = orderCreatedAt + " +0800 UTC"
				setCreatedAt, err := time.Parse(layout, orderCreatedAt)
				orderData.SetCreatedAt(setCreatedAt)
				log.Error(err)
			}

			if takeoutOrder.OrderType == 1 {
				deliveryTime := 0
				for _, item := range orderDeliveryTimeList {
					if item["storeId"].(int64) == orderData.StoreId {
						if item["kmStart"].(float64) <= takeoutOrder.DeliveryDistance && item["kmEnd"].(float64) >= takeoutOrder.DeliveryDistance {
							deliveryTime = item["minute"].(int)
							break
						}
					}
				}
				mm, _ := time.ParseDuration(strconv.Itoa(deliveryTime) + "m")
				takeoutOrder.EstimatedTime = takeoutOrder.EstimatedTime.Add(mm)
			}

			estimatedTimeString = hit.If(takeoutOrder.OrderType == 0, "้?่จๅ้คๆ้", "้?่จ้้คๆ้").(string)

			WeekDayMap := map[string]string{
				"Monday":    "ไธ",
				"Tuesday":   "ไบ",
				"Wednesday": "ไธ",
				"Thursday":  "ๅ",
				"Friday":    "ไบ",
				"Saturday":  "ๅญ",
				"Sunday":    "ๆฅ",
			}

			hour, _ := strconv.Atoi(takeoutOrder.EstimatedTime.Format("15"))
			hourString := strconv.Itoa(hour)
			minute, _ := strconv.Atoi(takeoutOrder.EstimatedTime.Format("04"))
			timeStart := ""
			timeEnd := ""
			if minute >= 0 && minute < 15 {
				timeStart = hourString + ":00"
				timeEnd = hourString + ":15"
			} else if minute >= 15 && minute < 30 {
				timeStart = hourString + ":15"
				timeEnd = hourString + ":30"
			} else if minute >= 30 && minute < 45 {
				timeStart = hourString + ":30"
				timeEnd = hourString + ":45"
			} else if minute >= 45 && minute <= 59 {
				nextHour := hit.If(hour == 23, "00", strconv.Itoa(hour+1)).(string)
				timeStart = hourString + ":45"
				timeEnd = nextHour + ":00"
			}

			day := takeoutOrder.EstimatedTime.Format("01ๆ02ๆฅ")
			week := WeekDayMap[takeoutOrder.EstimatedTime.Weekday().String()]
			timeString := timeStart + "~" + timeEnd

			estimatedTimeString = estimatedTimeString + " " + day + "(" + week + ")"

			data := map[string]interface{}{
				"takeoutId":           takeoutOrder.Id,
				"orderId":             orderData.Id,
				"storeId":             orderData.StoreId,
				"invoicePrice":        orderData.InvoicePrice,
				"totalPrice":          orderData.SubPrice,
				"mealAmount":          mealAmount,
				"orderCreatedAt":      orderData.CreatedAt.Format("2006-01-02 15:04:05"), //ไปๆฌพๆ้
				"orderSn":             orderData.OrderSn,
				"orderFromEN":         takeoutOrder.OrderFrom,
				"takeoutType":         takeoutOrder.OrderType,
				"storeInfo":           map[string]interface{}{"storeName": store.Name, "city": store.City, "township": store.Township, "address": store.Address, "brandName": brand.Name, "brandImage": image.ReturnPhotoPath(brand.Image)},
				"estimatedTime":       takeoutOrder.EstimatedTime.Format("2006-01-02 15:04"),
				"cancelReason":        takeoutOrder.CancelReason,
				"takeoutStatus":       takeoutOrder.Status,
				"isUserCancel":        takeoutOrder.IsUserCancel,
				"estimatedTimeString": estimatedTimeString,
				"timeString":          timeString,
				"isContractAvailable": isContractAvailable, //ๅ็ๅ็ดๆฏๅฆๆๆ๏ผ
			}
			dataList = append(dataList, data)
			return
		})
	return util.RS{Message: "", Data: dataList, Status: true}
}

func (handler *OrderHandler) GetOnlineOrderDetail() interface{} {
	var (
		data            ReceivedOrder
		orderRepository repository.OrderRepository
		takeoutList     repository.TakeoutDetailRepository
		sendData        interface{}
		// dataList = make([]map[string]interface{}, 0)
	)
	if err := json.Unmarshal([]byte(handler.Parameter), &data); err != nil {
		return util.RS{Message: "", Status: false}
	} else if order := orderRepository.QueryOrder(handler.UserId, data.OrderSn); order.Id == 0 {
		return util.RS{Message: "่จๅฎ็ทจ่ไธๅญๅจ", Status: false}
	} else {
		sendData = takeoutList.GetTakeoutOrderDetail(order)

	}
	return util.RS{Message: "", Data: sendData, Status: true}
}

func (handler *OrderHandler) CreateOrder2() interface{} {

	var (
		data                                ReceivedOrder
		fcm                                 fcm.FCM
		store                               model.Store
		users                               model.Users
		userAddress                         model.UserAddress
		orderRepository                     repository.OrderRepository
		storeRepository                     repository.StoreRepository
		onlineOrderDeliveryConditionSetting model.OnlineOrderDeliveryConditionSetting
		onlineOrderDeliveryTimeSetting      model.OnlineOrderDeliveryTimeSetting
		posMenuRepository                   repository.PosMenuRepository
		posMenuTypeDetail                   model.PosMenuTypeDetail
		orderDiscountLimitSetting           model.OrderDiscountLimitSetting
	)

	paramStr, _ := crypto.KeyDecrypt(handler.Parameter) //่งฃๅฏ
	if err := json.Unmarshal([]byte(paramStr), &data); err != nil {
		log.Error(err)
		return util.RS{Message: "", Status: false}
	} else if handler.UserId == 0 {
		return util.RS{Message: "่ซ็ปๅฅๆๅก", Status: false}
	} else if handler.UserId == -1 && data.Phone == "" {
		return util.RS{Message: "่ซ่ผธๅฅๆๆฉ่็ขผ", Status: false}
	} else if msg, isValidate := validate.CheckPhone(data.Phone); !isValidate && handler.UserId == -1 {
		return util.RS{Message: msg, Status: false}
	} else if handler.UserId == -1 && data.TableId == 0 {
		log.Error(errors.New("handler.UserId == -1 && data.TableId == "))
		return util.RS{Message: "", Status: false}
	} else if data.StoreId == 0 {
		log.Error(errors.New("data.StoreId == 0"))
		return util.RS{Message: "", Status: false}
	} else if len(data.TaxIdNumber) > 0 && len(data.TaxIdNumber) != 8 {
		return util.RS{Message: "่ซ็ขบ่ช็ตฑไธ็ทจ่้ทๅบฆ", Status: false}
	} else if data.PaymentType == 0 {
		return util.RS{Message: "่ซ้ธๆไปๆฌพๆนๅผ", Status: false}
	} else if data.RadioTake == 2 && data.UserAddressId == 0 {
		return util.RS{Message: "่ซ้ธๆๅค้ๅฐๅ", Status: false}
	}

	var orderId int64
	var orderSn string
	var isEditOrder, isTableOrder bool
	var oldTakeoutOrderEstimatedTime time.Time
	var oldTakeoutOrderTakeoutSerialNumber string

	if handler.UserId == -1 && data.Phone != "" {
		var (
			checkUsers model.Users
		)
		if checkUsers.SetCountryCode("+886").SetPhone(data.Phone).GetUsersIdByPhone(); checkUsers.Id != 0 {
			handler.UserId = checkUsers.Id
		} else {
			userId, err := checkUsers.SetName("็ดฏ็ฉๆๅก").Insert()

			if err == nil {
				handler.UserId = userId
			} else {
				log.Error(err)
				return util.RS{Message: "ๆฐๅข็ดฏ็ฉๆๅก็ฐๅธธ", Status: false}
			}
		}

		isTableOrder = true
	}

	if data.TakeoutId != 0 {
		var takeoutOrder model.TakeoutOrder
		if takeoutOrder.SetId(data.TakeoutId).SetUserId(handler.UserId).QueryOneJoinOrder(); takeoutOrder.Name == "" {
			return util.RS{Message: "ๆฅ็กๆญค่จๅฎ", Status: false}
		} else if takeoutOrder.Status == 4 {
			return util.RS{Message: "่จๅฎๅทฒๅๆถ๏ผ็กๆณไฟฎๆน่จๅฎ", Status: false}
		} else if takeoutOrder.Status == 5 {
			return util.RS{Message: "่จๅฎๅทฒๅฎๆ๏ผ็กๆณไฟฎๆน่จๅฎ", Status: false}
		} else if takeoutOrder.Status != 0 && data.IsAccepted == false {
			return util.RS{Message: "ไฟฎๆนๅคฑๆ๏ผๅบๅฎถๅทฒๆฅๅฎ๏ผ่ฅๆฌฒไฟฎๆน่ซ้ๆฐๆไฝ", Status: false}
		}

		orderId = takeoutOrder.OrderId
		orderSn = takeoutOrder.OrderSn
		isEditOrder = true
		oldTakeoutOrderEstimatedTime = takeoutOrder.EstimatedTime
		oldTakeoutOrderTakeoutSerialNumber = takeoutOrder.TakeoutSerialNumber
	}

	if data.ReservationTime != "" {
		local, _ := time.LoadLocation("Asia/Taipei") //ไฟฎๆนๆๅฐๅๆ้
		timeNow, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), local)
		timeReservation, _ := time.ParseInLocation("2006-01-02 15:04:05", data.ReservationTime, local)
		if timeReservation.Before(timeNow) {
			return util.RS{Message: "ไธๅฏ้?้ป็พๅจไนๅ", Status: false}
		}
	}

	type PosGoodsNoteDetailList struct {
		Id            int64  `json:"id"`
		Name          string `json:"name"`
		IncreasePrice int    `json:"increasePrice"`
	}

	type PosFoodGoodsNoteList struct {
		PosGoodsNoteDetailList []PosGoodsNoteDetailList `json:"posGoodsNoteDetailList"`
	}

	type PosFoodSubmealList struct {
		PosFoodId            int64                  `json:"posFoodId"`
		PosFoodName          string                 `json:"posFoodName"`
		IncreasePrice        int                    `json:"increasePrice"`
		PosFoodGoodsNoteList []PosFoodGoodsNoteList `json:"posFoodGoodsNoteList"`
	}

	type DataList struct {
		PosGoodsNoteDetailList []PosGoodsNoteDetailList `json:"posGoodsNoteDetailList"`
		PosFoodSubmealList     []PosFoodSubmealList     `json:"posFoodSubmealList"`
	}

	type PosFoodInfo struct {
		Id              int64  `json:"id"`
		OnlineOrderName string `json:"onlineOrderName"`
		OrderType       int    `json:"orderType"`
	}

	type DataJson struct {
		DataList    []DataList  `json:"dataList"`
		PosFoodInfo PosFoodInfo `json:"posFoodInfo"`
	}

	now := time.Now()
	var takeoutInfo TakeoutOrderInfo
	var totalPrice, deliveryFee, deliveryDistance float64
	var subPrice int64

	if len(data.MealChoiceList) == 0 {
		return util.RS{Message: "่ซ้ธๆๅๅ", Status: false}
	}
	storeInfo := storeRepository.SetId(data.StoreId).GetStoreOrderSetting()
	menuData := posMenuRepository.SetStoreId(data.StoreId).SetReservationTime(data.ReservationTime).GetStoreMenuData()

	if !storeInfo["onlineOrderStatus"].(bool) || menuData["id"].(int64) == 0 {
		return util.RS{Message: "็ฎๅ้่ฉฒๅบๅฎถ็ทไธ้ป้คๆๆฎต๏ผ็กๆณ่จ้คใ", Status: false}
	}

	orderFoodList := make([]map[string]interface{}, 0)
	orderFoodIdList := make([]interface{}, 0)
	mealChoiceList := data.MealApiFormat
	for _, item := range mealChoiceList {
		subPrice += item.ChoicePrice
		isRepeat := false
		for _, id := range orderFoodIdList {
			if id == item.Id {
				isRepeat = true
				break
			}
		}

		if !isRepeat {
			orderFoodIdList = append(orderFoodIdList, item.Id)
			orderFoodList = append(orderFoodList, map[string]interface{}{
				"id":   item.Id,
				"name": item.Name,
			})
		}
	}

	foodList := posMenuTypeDetail.
		SetMenuTypeIdList(menuData["menuTypeIdList"].([]interface{})).
		SetFoodIdList(orderFoodIdList).
		GetOnlineOrderFoodList()

	emptyFoodNameList := make([]string, 0)
	orderFoodDetailList := make([]map[string]interface{}, 0)

	for _, orderFoodItem := range orderFoodList {
		hasData := false
		for _, foodItem := range foodList {
			if orderFoodItem["id"].(int64) == foodItem["posFoodId"] {
				hasData = true
				break
			}
		}

		if hasData {
			var posFoodRepository repository.PosFoodRepository
			posFoodData := posFoodRepository.SetId(orderFoodItem["id"].(int64)).GetPosFoodDetail()
			orderFoodDetailList = append(orderFoodDetailList, posFoodData)
		} else {
			emptyFoodNameList = append(emptyFoodNameList, orderFoodItem["name"].(string))
		}
	}

	message := ""

	if len(emptyFoodNameList) > 0 {
		message += strings.Join(emptyFoodNameList, ",") + "๏ผ็ฎๅ้ป้คๆๆฎตไธๆไพ่ฉฒ้ค้ป"
		return util.RS{Message: message, Status: false}
	}

	updateFoodIdList := make([]int64, 0)
	updateFoodNameList := make([]string, 0)

	for _, choiceItem := range mealChoiceList {
		isRepeatUpdate := false
		for _, id := range updateFoodIdList {
			if choiceItem.Id == id {
				isRepeatUpdate = true
				break
			}
		}

		if !isRepeatUpdate {
			for _, foodItem := range orderFoodDetailList {
				posFoodInfo := foodItem["posFoodInfo"].(map[string]interface{})
				dataList := foodItem["dataList"].([]map[string]interface{})
				var price int
				price = posFoodInfo["onlineOrderPrice"].(int)

				if choiceItem.Id == posFoodInfo["id"] {
					var isUpdate bool

					if posFoodInfo["orderType"] == 1 {
						for _, choiceNoteItem := range choiceItem.PosGoodsNoteDetailList {
							for _, noteTypeItem := range dataList {
								for _, noteItem := range noteTypeItem["posGoodsNoteDetailList"].([]map[string]interface{}) {
									if choiceNoteItem.Id == noteItem["id"] {
										price += noteItem["increasePrice"].(int)
										break
									}
								}
							}
						}
					} else {
						if len(choiceItem.PosFoodSubmealList) == len(dataList) {

							for foodSubMealIndex, foodSubMealItem := range choiceItem.PosFoodSubmealList {
								for _, foodSubItem := range dataList[foodSubMealIndex]["posFoodSubmealList"].([]map[string]interface{}) {
									if foodSubMealItem.PosFoodId == int64(foodSubItem["posFoodId"].(float64)) {
										price += int(foodSubItem["increasePrice"].(float64))
										for _, choiceNoteItem := range foodSubMealItem.PosGoodsNoteDetailList {
											for _, noteTypeItem := range foodSubItem["posFoodGoodsNoteList"].([]interface{}) {
												ins := reflect.ValueOf(noteTypeItem.(map[string]interface{})["posGoodsNoteDetailList"])
												if ins.Len() > 0 {
													for _, noteItem := range noteTypeItem.(map[string]interface{})["posGoodsNoteDetailList"].([]interface{}) {
														if choiceNoteItem.Id == int64(noteItem.(map[string]interface{})["id"].(float64)) {
															price += int(noteItem.(map[string]interface{})["increasePrice"].(float64))
															break
														}
													}
												}
											}
										}
										break
									}
								}
							}
						} else {
							isUpdate = true
						}
					}

					price = price * int(choiceItem.Amount)

					if int(choiceItem.ChoicePrice) != price || isUpdate {
						updateFoodIdList = append(updateFoodIdList, choiceItem.Id)
						updateFoodNameList = append(updateFoodNameList, choiceItem.Name)
					}
				}
			}
		}
	}

	if len(updateFoodNameList) > 0 {
		message += strings.Join(updateFoodNameList, ",") + "๏ผ้ค้ป่ณๆๆ็ฐๅ๏ผ่ซๅช้ค่ฉฒ้ค้ปๅพ้ๆฐ้ธๅ"
		return util.RS{Message: message, Status: false}
	}

	totalPrice = float64(subPrice)

	offsetList := make([]map[string]interface{}, 0)
	discountList := make([]map[string]interface{}, 0)

	for _, item := range data.CouponChoiceList {
		residue := 0
		idAry := make([]int64, 0)

		var userCoupon model.CouponUser
		var userCouponAry []interface{}
		var (
			useCouponListId   int64
			useCouponListName string
			useDiscountAmount int
		)

		userCoupon.
			SetCouponId(item.CouponListId).
			SetStatusAry([]interface{}{0}).
			SetCouponTypeId(item.CouponTypeId).
			SetUserId(handler.UserId).
			SetStartTime(item.StartTime).
			SetEndTime(item.EndTime).
			QueryAllJoinList(func(rs *model.CouponUser) {
				data := map[string]interface{}{
					"id":             rs.Id,
					"status":         rs.Status,
					"code":           rs.Code,
					"couponListId":   rs.CouponListId,
					"couponTypeId":   rs.CouponTypeId,
					"isStart":        rs.IsStart,
					"isUse":          rs.IsUse,
					"discountAmount": rs.DiscountAmount,
					"startTime":      rs.StartTime,
					"endTime":        rs.EndTime,
				}
				userCouponAry = append(userCouponAry, data)
			})

		for _, item2 := range userCouponAry {
			if item2.(map[string]interface{})["isUse"] != 0 && item2.(map[string]interface{})["isStart"] != 0 {
				idAry = append(idAry, item2.(map[string]interface{})["id"].(int64))
				residue += 1
			}
		}

		//ๆฌๆฌกๆๅขๅ?ไฝฟ็จ็ๅชๆ?ๅธ
		if len(userCouponAry) > 0 {
			useCouponListId = userCoupon.CouponListId
			useCouponListName = userCoupon.CouponListName
			useDiscountAmount = userCoupon.DiscountAmount
		}

		if isEditOrder {
			var orderDiscount model.OrderDiscount
			orderDiscount.
				SetOrderSn(orderSn).
				SetCouponListId(item.CouponListId).
				SetStartTime(item.StartTime).
				SetEndTime(item.EndTime).
				QueryAllJoinList(func(rs *model.OrderDiscount) {
					idAry = append(idAry, rs.CouponUserId)
					residue += 1
				})
			//ๅทฒไฝฟ็จ็ๅชๆ?ๅธ
			//ๅ?็บไฟฎๆน่จๅฎๆๅฏ่ฝๆuserCouponๆฒๆ็ๆๆณ๏ผไปฃ่กจๅชๆ?ๅธๆฏๅทฒ่ขซไฝฟ็จ็๏ผๆไปฅๅฐฑ่ฆ่จๅฎ็่ณๆ
			if len(userCouponAry) == 0 && residue > 0 {
				useCouponListId = orderDiscount.CouponListId
				useCouponListName = orderDiscount.CouponName
				useDiscountAmount = orderDiscount.DiscountAmount
			}
		}

		if item.ChoiceAmount > residue {
			return util.RS{Message: item.Name + "ๆธ้ไธ่ถณ", Status: false}
		} else if item.CouponTypeId == 2 {
			offsetList = append(offsetList, map[string]interface{}{
				"idAry":          idAry,
				"choiceAmount":   item.ChoiceAmount,
				"couponListId":   useCouponListId,
				"couponTypeId":   item.CouponTypeId,
				"couponListName": useCouponListName,
				"discountAmount": useDiscountAmount,
				"startTime":      item.StartTime,
				"endTime":        item.EndTime,
			})
		} else if item.CouponTypeId == 4 {
			discountList = append(discountList, map[string]interface{}{
				"idAry":          idAry,
				"choiceAmount":   item.ChoiceAmount,
				"couponListId":   useCouponListId,
				"couponTypeId":   item.CouponTypeId,
				"couponListName": useCouponListName,
				"discountAmount": useDiscountAmount,
				"startTime":      item.StartTime,
				"endTime":        item.EndTime,
			})
		}
	}

	couponList := make([]map[string]interface{}, 0)

	for _, item := range offsetList {
		totalPrice -= float64(item["discountAmount"].(int)) * float64(item["choiceAmount"].(int))
		if totalPrice < 0 {
			totalPrice = 0
		}

		couponList = append(couponList, map[string]interface{}{
			"couponUserIdAry": item["idAry"].([]int64),
			"couponListId":    item["couponListId"].(int64),
			"couponListName":  item["couponListName"].(string),
			"couponDiscount":  float64(item["discountAmount"].(int)),
			"couponTypeId":    item["couponTypeId"].(int),
			"choiceAmount":    item["choiceAmount"].(int),
			"couponType":      "offset",
			"startTime":       item["startTime"].(string),
			"endTime":         item["endTime"].(string),
		})
	}

	for _, item := range discountList {
		var rate, discountComputed float64
		discountComputedAry := make([]float64, 0)
		rate = float64(100-item["discountAmount"].(int)) / 100

		for i := 0; i < item["choiceAmount"].(int); i++ {
			if storeInfo["discountRole"].(int) == 0 {
				discountComputed = totalPrice - math.Floor((totalPrice*rate)+0.5)
				totalPrice = math.Floor((totalPrice * rate) + 0.5)
			} else if storeInfo["discountRole"].(int) == 1 {
				discountComputed = totalPrice - math.Ceil(totalPrice*rate)
				totalPrice = math.Ceil(totalPrice * rate)
			} else if storeInfo["discountRole"].(int) == 2 {
				discountComputed = totalPrice - math.Floor(totalPrice*rate)
				totalPrice = math.Floor(totalPrice * rate)
			}

			discountComputedAry = append(discountComputedAry, discountComputed)
		}

		couponList = append(couponList, map[string]interface{}{
			"couponUserIdAry":     item["idAry"].([]int64),
			"couponListId":        item["couponListId"].(int64),
			"couponListName":      item["couponListName"].(string),
			"choiceAmount":        item["choiceAmount"].(int),
			"couponDiscount":      discountComputed,
			"couponTypeId":        item["couponTypeId"].(int),
			"couponType":          "discount",
			"startTime":           item["startTime"].(string),
			"endTime":             item["endTime"].(string),
			"discountComputedAry": discountComputedAry,
		})

	}
	orderDiscountLimitSetting.SetStoreId(data.StoreId).QueryOne()

	switch data.RadioTake {
	case 1:
		switch orderDiscountLimitSetting.OnlineTakeoutType {
		case 1:
			if tempLimitPrice := orderDiscountLimitSetting.OnlineTakeoutPrice; float64(subPrice)-totalPrice > float64(tempLimitPrice) {
				totalPrice = float64(subPrice) - float64(tempLimitPrice)
			}
		case 2:
			if tempLimitPrice := math.Floor((float64(subPrice) * float64(orderDiscountLimitSetting.OnlineTakeoutPercentage) / 100) + 0.5); float64(subPrice)-totalPrice > tempLimitPrice {
				totalPrice = float64(subPrice) - tempLimitPrice
			}
		}
	case 2:
		switch orderDiscountLimitSetting.OnlineDeliveryType {
		case 1:
			if tempLimitPrice := orderDiscountLimitSetting.OnlineDeliveryPrice; float64(subPrice)-totalPrice > float64(tempLimitPrice) {
				totalPrice = float64(subPrice) - float64(tempLimitPrice)
			}
		case 2:
			if tempLimitPrice := math.Floor((float64(subPrice) * float64(orderDiscountLimitSetting.OnlineDeliveryPercentage) / 100) + 0.5); float64(subPrice)-totalPrice > tempLimitPrice {
				totalPrice = float64(subPrice) - tempLimitPrice
			}
		}
	}
	estimatedTime := now
	reduceMm, _ := time.ParseDuration("-0m")

	if data.RadioTake == 2 {
		userAddress.
			SetId(data.UserAddressId).
			SetUserId(handler.UserId).
			SetTargetLat(storeInfo["storeLat"].(float64)).
			SetTargetLng(storeInfo["storeLng"].(float64)).
			SetCheckDeleteStatus(true).
			SetDeleteStatus(false).
			QueryOne()

		if userAddress.Address == "" {
			return util.RS{Message: "ๆฅ็กๅค้ๅฐๅ่ณๆ", Status: false}
		} else {
			takeoutInfo.ClientAddress = userAddress.City + userAddress.Township + userAddress.Address
			takeoutInfo.SelectArea = userAddress.PostalCode
			takeoutInfo.DeliveryInstructions = userAddress.Content

			onlineOrderDeliveryConditionSetting.
				SetStoreId(data.StoreId).
				SetCheckStatus(true).
				SetStatus(true).
				SetCheckDeleteStatus(true).
				SetDeleteStatus(false).
				SetFloatDistance(userAddress.Distance).
				SetRequiredAmount(int(subPrice)).
				GetMinDeliveryFee()

			if onlineOrderDeliveryConditionSetting.Id == 0 {
				return util.RS{Message: "ๆฅ็กๆปฟ่ถณๆขไปถ็้่ฒป่ณๆ๏ผๅฏ่ฝ็กๆณๆไพๅค้", Status: false}
			} else {
				deliveryFee = float64(onlineOrderDeliveryConditionSetting.DeliveryFee)
				deliveryDistance = userAddress.Distance

				deliveryMinute := onlineOrderDeliveryTimeSetting.
					SetStoreId(data.StoreId).
					SetDistance(userAddress.Distance).
					SetCheckDeleteStatus(true).
					SetDeleteStatus(false).
					QueryOne().
					Minute

				addDeliveryMm, _ := time.ParseDuration(strconv.Itoa(deliveryMinute) + "m")
				reduceDeliveryMm, _ := time.ParseDuration("-" + strconv.Itoa(deliveryMinute) + "m")

				estimatedTime = estimatedTime.Add(addDeliveryMm)
				reduceMm = reduceDeliveryMm
			}
		}
	}

	totalPrice += deliveryFee
	if data.TotalPrice != 0 && data.TotalPrice != totalPrice {
		log.Error(errors.New("็ทไธ้ป้คๅๅพ็ซฏๆไป้้กๆ็ฐ๏ผ่ซ็ขบ่ชๅๅพ็ซฏ้็ฎๆฏๅฆไธ่ดๆ่จญๅฎๆๆ็ฐๅ" + strconv.Itoa(int(data.TotalPrice)) + " != " + strconv.Itoa(int(totalPrice))))
		return util.RS{Message: "ๆไป้้กๆ็ฐ๏ผ่ซ้ๆฐ้ฒๅฅ่ณผ็ฉ่ป", Status: false}
	}

	userInfo := users.GetUsersInfoWithId(handler.UserId)

	gender := 2

	if int(userInfo["gender"].(int64)) == 1 {
		gender = 0
	} else if int(userInfo["gender"].(int64)) == 2 {
		gender = 1
	}

	takeoutInfo.ClientName = userInfo["name"].(string)
	takeoutInfo.ClientPhone = userInfo["phone"].(string)
	takeoutInfo.BuyerBan = data.TaxIdNumber
	takeoutInfo.Gender = gender

	takeoutInfo.CustomerNote = data.Remark
	takeoutInfo.TakeoutTypeId = data.RadioTake
	takeoutInfo.IsCheckedOut = data.PaymentType != 1

	addMm, _ := time.ParseDuration("0m")
	mealPreparationTimeType := storeInfo["mealPreparationTimeType"].(int)
	switch mealPreparationTimeType {
	case 1:
		addMm, _ = time.ParseDuration("15m")
	case 2:
		addMm, _ = time.ParseDuration("30m")
	case 3:
		addMm, _ = time.ParseDuration("45m")
	case 4:
		addMm, _ = time.ParseDuration("60m")
	}

	if data.ReservationTime == "" {
		takeoutInfo.EstimatedTime = now.Add(addMm)
		estimatedTime = estimatedTime.Add(addMm)
	} else {
		local, _ := time.LoadLocation("Asia/Taipei") //ไฟฎๆนๆๅฐๅๆ้
		reservationTime, _ := time.ParseInLocation("2006-01-02 15:04:05", data.ReservationTime, local)
		takeoutInfo.EstimatedTime = reservationTime.Add(reduceMm)
		estimatedTime = reservationTime
		takeoutInfo.Type = 1
	}

	uuid, err := crypto.KeyEncrypt(strconv.FormatInt(data.StoreId, 10) + "OnlineCreateOrder" + "")
	log.Error(err)

	appListData := store.SetId(data.StoreId).GetAppListData()
	appListId := appListData["appListId"].(int64)

	jsonMap := map[string]interface{}{
		"takeoutId":           data.TakeoutId,
		"orderId":             orderId,
		"orderSn":             orderSn,
		"isEditOrder":         isEditOrder,
		"mealChoiceList":      mealChoiceList,
		"isTakeoutOrder":      true,
		"takeoutInfo":         takeoutInfo,
		"uuid":                uuid,
		"userId":              handler.UserId,
		"appId":               appListId,
		"storeId":             data.StoreId,
		"subPrice":            subPrice,
		"totalPrice":          totalPrice,
		"invoicePrice":        totalPrice,
		"actualPrice":         totalPrice,
		"paymentType":         data.PaymentType,
		"deliveryFee":         deliveryFee,
		"deliveryDistance":    deliveryDistance,
		"couponList":          couponList,
		"buyerBan":            data.TaxIdNumber,
		"isRequiredTableware": data.IsRequiredTableware,
		"isCordova":           data.IsCordova,
		"isTableOrder":        isTableOrder,
		"tableId":             data.TableId,
		"adult":               data.Adult,
		"child":               data.Child,
	}

	if data.PaymentType == 1 {
		_, err = kafka.Push("OnlineCreateOrder2", jsonMap)
		log.Error(err)
		dataList := make(map[string]interface{}, 0)
		//kafka.Listen(strconv.FormatInt(data.StoreId, 10)+"OnlineAppApiCreateOrder", uuid, func(kafKaRS util.KafKaRS) {
		//	dataList["status"] = kafKaRS.Status
		//	dataList["message"] = kafKaRS.Message
		//	dataList["data"] = kafKaRS.Data
		//})

		for {
			if kafKaRS := <-kafka.KafkaChan["OnlineAppApiCreateOrder"]; kafKaRS.Uuid == uuid {
				dataList["status"] = kafKaRS.Status
				dataList["message"] = kafKaRS.Message
				dataList["data"] = kafKaRS.Data
				break
			}
		}

		if dataList["status"] == true {

			takeoutSnString := dataList["data"].(map[string]interface{})["takeoutSnString"].(string)
			orderSn := dataList["data"].(map[string]interface{})["orderSn"].(string)
			isFcm := false
			fmcBody := ""
			fmcTitle := ""
			toDay := now.Format("2006-01-02")

			//้ไฟฎๆน็็ซๅณๆ้?้ป็ถๅคฉๆ
			if !isEditOrder &&
				((data.ReservationTime == "") ||
					(toDay == takeoutInfo.EstimatedTime.Format("2006-01-02"))) {

				isFcm = true
				fmcBody = "#" + takeoutSnString + " ่จๅฎ็ทจ่:" + orderSn + " ่ซ็ขบ่ช่จๅฎ"
				fmcTitle = "ๆๆฐ็ๅคๅธถๅค้่จๅฎ๏ผ่ซๅๅพๅพ็ขบ่ช่็"
				//ไฟฎๆน่จๅฎๆ
			} else if isEditOrder {
				//ๆนๅ่ทๆนๅพ้?้ปๆ้้ฝๆฏไปๅคฉๆ
				if (toDay == oldTakeoutOrderEstimatedTime.Format("2006-01-02")) &&
					(toDay == takeoutInfo.EstimatedTime.Format("2006-01-02")) {

					isFcm = true
					fmcBody = "#" + takeoutSnString + " ่จๅฎ็ทจ่:" + orderSn + " ่ซ็ขบ่ช่จๅฎ"
					fmcTitle = "ๆไฟฎๆนๅคๅธถๅค้่จๅฎ๏ผ่ซๅๅพๅพ็ขบ่ช่็"
					//ไปๅคฉ็่จๅฎๆนๅฐๅถไปๅคฉๆ
				} else if (toDay == oldTakeoutOrderEstimatedTime.Format("2006-01-02")) &&
					(toDay != takeoutInfo.EstimatedTime.Format("2006-01-02")) {

					oldEstimatedTimeYear, _ := strconv.Atoi(oldTakeoutOrderEstimatedTime.Format("2006"))
					newEstimatedTimeYear, _ := strconv.Atoi(takeoutInfo.EstimatedTime.Format("2006"))

					oldEstimatedTimeString := strconv.Itoa(oldEstimatedTimeYear-1911) + "/" + oldTakeoutOrderEstimatedTime.Format("01/02 15:04")
					newEstimatedTimeString := strconv.Itoa(newEstimatedTimeYear-1911) + "/" + takeoutInfo.EstimatedTime.Format("01/02 15:04")

					isFcm = true
					fmcBody = "่จๅฎๆ้็ฐๅ๏ผๅคๅธถๅค้่จๅฎ #" + oldTakeoutOrderTakeoutSerialNumber + " ๅไปๆฅ้?่จๅ้คๆ้ " + oldEstimatedTimeString + "๏ผๆถ่ฒป่ๅทฒๆดๆน่ณ" + newEstimatedTimeString + " ่จๅฎๆจ่็บ#" + takeoutSnString
					fmcTitle = "ๆไฟฎๆนๅคๅธถๅค้่จๅฎ๏ผ่ซๅๅพๅพ็ขบ่ช่็"
					//ๅถไปๅคฉ็่จๅฎๆนๅฐไปๅคฉๆ
				} else if (toDay != oldTakeoutOrderEstimatedTime.Format("2006-01-02")) &&
					(toDay == takeoutInfo.EstimatedTime.Format("2006-01-02")) {

					isFcm = true
					fmcBody = "#" + takeoutSnString + " ่จๅฎ็ทจ่:" + orderSn + " ่ซ็ขบ่ช่จๅฎ"
					fmcTitle = "ๆๆฐ็ๅคๅธถๅค้่จๅฎ๏ผ่ซๅๅพๅพ็ขบ่ช่็"
				}
			}

			if isFcm {
				fcm.SetBody(fmcBody).
					SetStoreId(data.StoreId).
					SetTitle(fmcTitle).
					SetFcmServerKey(config.ServerInfo.FcmServerKey).
					SetData(map[string]interface{}{
						"action":    "takeoutOrder",
						"type":      hit.If(isEditOrder, "editOrder", "newOrder"),
						"takeoutId": dataList["data"].(map[string]interface{})["takeoutId"],
					}).
					SetManagementId(TakeoutManagementId).
					SendAdminStoreId()
			}

			WeekDayMap := map[string]string{
				"Monday":    "ๆๆไธ",
				"Tuesday":   "ๆๆไบ",
				"Wednesday": "ๆๆไธ",
				"Thursday":  "ๆๆๅ",
				"Friday":    "ๆๆไบ",
				"Saturday":  "ๆๆๅญ",
				"Sunday":    "ๆๆๆฅ",
			}

			hour, _ := strconv.Atoi(estimatedTime.Format("15"))
			hourString := strconv.Itoa(hour)
			minute, _ := strconv.Atoi(estimatedTime.Format("04"))
			timeStart := ""
			timeEnd := ""
			if minute >= 0 && minute < 15 {
				timeStart = hourString + ":00"
				timeEnd = hourString + ":15"
			} else if minute >= 15 && minute < 30 {
				timeStart = hourString + ":15"
				timeEnd = hourString + ":30"
			} else if minute >= 30 && minute < 45 {
				timeStart = hourString + ":30"
				timeEnd = hourString + ":45"
			} else if minute >= 45 && minute <= 59 {
				nextHour := hit.If(hour == 23, "00", strconv.Itoa(hour+1)).(string)
				timeStart = hourString + ":45"
				timeEnd = nextHour + ":00"
			}

			data := map[string]interface{}{
				"orderId":    dataList["data"].(map[string]interface{})["orderId"],
				"orderSn":    dataList["data"].(map[string]interface{})["orderSn"],
				"takeoutId":  dataList["data"].(map[string]interface{})["takeoutId"],
				"userName":   userInfo["name"].(string),
				"radioTake":  data.RadioTake,
				"totalPrice": totalPrice,
				"day":        estimatedTime.Format("2006-01-02"),
				"time":       timeStart + "~" + timeEnd,
				"week":       WeekDayMap[estimatedTime.Weekday().String()],
			}

			return util.RS{Message: dataList["message"].(string), Status: true, Data: data}
		} else {
			return util.RS{Message: dataList["message"].(string), Status: false}
		}
	} else if data.PaymentType == 9 {
		var (
			orderTmp     model.OrderTmp
			responseData OcpResultFormat
			storeUid     = "451339810001" //ไธฒ้ๅญๅญๅบไปฃ็ขผ(ๅบๅฎ)
		)

		jsonData, _ := json.Marshal(jsonMap)
		orderTmpId, _ := orderTmp.SetJsonContent(string(jsonData)).SetCreatedAt(time.Now()).Create()

		topUpItems := make([]map[string]interface{}, 0)

		//ๅๅพ่ณผ็ฉ่ปๆธๅฎ
		for _, item := range mealChoiceList {
			topUpItems = append(topUpItems, map[string]interface{}{
				"id":     strconv.FormatInt(item.Id, 10),
				"name":   item.Name,
				"cost":   item.ChoicePrice / item.Amount,
				"amount": item.Amount,
				"total":  item.ChoicePrice,
			})
		}

		//้ซ้้ปๅญ้ขๅ
		if paymentJsonText, err := json.Marshal(map[string]interface{}{ //ๆๅฎไบคๆ่ณ่จๅงๅฎนไธฆ่ฝๆjsonๅญไธฒ
			"store_uid": storeUid,
			"cost":      strconv.FormatFloat(jsonMap["totalPrice"].(float64), 'f', -1, 64),
			"items":     topUpItems,
			"order_id":  strconv.FormatInt(orderTmpId, 10),
			"echo_0":    config.MyPayInfo.MyPayFcmKey,
			"echo_1":    "SendOinApp",
			"user_id":   strconv.FormatInt(jsonMap["userId"].(int64), 10),
			"ip":        "127.0.0.1",
		}); err != nil {
			log.Error(err) //"ๅปบ็ซๅณ่ผธ่ณๆ้ฏ่ชค"
		} else if encryData, err := crypto.AesCBC256Encrpty(config.MyPayInfo.MypayKey, paymentJsonText); err != nil { //ๅฐไบคๆ่ณ่จๅ?ๅฏ
			log.Error(err) //"ๅฐไบคๆ่ณ่จๅ?ๅฏ้ฏ่ชค"
		} else if serviceJsonText, err := json.Marshal(map[string]interface{}{ //ๆๅ่ณ่จ่ฝไนjsonๅญไธฒ
			"service_name": "ocpap",
			"cmd":          "api/directdeal",
		}); err != nil {
			log.Error(err) //"ๅปบ็ซๆๅ่ณ่จ้ฏ่ชค"
		} else if service, err := crypto.AesCBC256Encrpty(config.MyPayInfo.MypayKey, serviceJsonText); err != nil { //ๆๅ่ณ่จๅ?ๅฏ
			log.Error(err) //"ๆๅ่ณ่จๅ?ๅฏ้ฏ่ชค"
		} else if response, err := rpc.ApiCall(config.MyPayInfo.MypayInnerPayUrl, "POST", map[string]interface{}{ //ๅผๅซ้ๆตapi
			"agent_uid":  config.MyPayInfo.AgentUid,
			"service":    service,
			"encry_data": encryData,
		}); err != nil {
			log.Error(err) //"ๅผๅซ้ๆตapi้ฏ่ชค"
		} else if err := json.Unmarshal([]byte(response), &responseData); err != nil { //ๅฐๅๅณjsonๅญไธฒ่ฝๆ็ฉไปถ
			log.Error(err) //"ๅฐๅๅณjsonๅญไธฒ่ฝๆ็ฉไปถ้ฏ่ชค"
		}

		data := map[string]interface{}{
			"orderTmpId": orderTmpId,
			"code":       responseData.Code,
			"url":        responseData.Url,
		}

		return util.RS{Message: "่ฝๅพ้ปๅญ้ขๅไปๆฌพไธญ", Status: true, Data: data}
	} else if data.PaymentType == 2 {
		var orderTmp model.OrderTmp
		jsonData, _ := json.Marshal(jsonMap)
		orderTmpId, _ := orderTmp.SetJsonContent(string(jsonData)).SetCreatedAt(time.Now()).Create()

		//็ข็้?ไผฐๅ้คๆ้
		mapResult := orderRepository.GetEstimatedTime(estimatedTime)

		data := map[string]interface{}{
			"orderTmpId": orderTmpId,
			"userName":   userInfo["name"].(string),
			"radioTake":  data.RadioTake,
			"totalPrice": totalPrice,
		}

		tradeInfoData := digitalpay.MergeMaps(mapResult, data)

		return util.RS{Message: "ไฟก็จๅกไปๆฌพ่็ไธญ...", Status: true, Data: tradeInfoData}
	} else {
		var orderTmp model.OrderTmp
		jsonData, _ := json.Marshal(jsonMap)
		orderTmpId, _ := orderTmp.SetJsonContent(string(jsonData)).SetCreatedAt(time.Now()).Create()

		//็ข็้?ไผฐๅ้คๆ้
		mapResult := orderRepository.GetEstimatedTime(estimatedTime)

		data := map[string]interface{}{
			"orderTmpId": orderTmpId,
			"userName":   userInfo["name"].(string),
			"radioTake":  data.RadioTake,
			"totalPrice": totalPrice,
		}

		tradeInfoData := digitalpay.MergeMaps(mapResult, data)

		return util.RS{Message: "็ถๅฎไฟก็จๅกไปๆฌพ่็ไธญ...", Status: true, Data: tradeInfoData}
	}
}

func (handler *OrderHandler) CancelOrder() interface{} {
	var (
		data          ReceivedOrder
		order         model.Orders
		takeoutOrder  model.TakeoutOrder
		orderDiscount model.OrderDiscount
		couponUser    model.CouponUser
		fcm           fcm.FCM
	)

	if err := json.Unmarshal([]byte(handler.Parameter), &data); err != nil {
		log.Error(err)
		return util.RS{Message: "", Status: false}
	} else if data.OrderId == 0 && data.TakeoutId == 0 {
		log.Error(errors.New("data.OrderId == 0 && data.TakeoutId == 0"))
		return util.RS{Message: "", Status: false}
	} else if data.CancelReason == "" {
		return util.RS{Message: "่ซ้ธๆๅๆถ่จๅฎๅๅ?", Status: false}
	} else if takeoutOrder.SetId(data.TakeoutId).SetOrderId(data.OrderId).SetUserId(handler.UserId).QueryOneJoinOrder(); takeoutOrder.Name == "" {
		return util.RS{Message: "ๆฅ็กๆญค่จๅฎ", Status: false}
	} else if takeoutOrder.Status == 4 {
		return util.RS{Message: "่จๅฎๅทฒๅๆถ๏ผ็กๆณ้่คๅๆถ่จๅฎ", Status: false}
	} else if takeoutOrder.Status == 5 {
		return util.RS{Message: "่จๅฎๅทฒๅฎๆ๏ผ็กๆณๅๆถ่จๅฎ", Status: false}
	} else if takeoutOrder.Status != 0 && data.IsAccepted == false {
		return util.RS{Message: "ๅๆถๅคฑๆ๏ผๅบๅฎถๅทฒๆฅๅฎ๏ผ่ฅๆฌฒๅๆถ่ซ้ๆฐๆไฝ", Status: false}
	}

	takeoutOrderStatus := takeoutOrder.Status

	err := takeoutOrder.
		SetId(takeoutOrder.Id).
		SetStatus(4).
		SetCancelReason(data.CancelReason).
		SetCanceledAt(time.Now()).
		SetIsUserCancel(true).
		Update([]string{"status", "cancel_reason", "canceled_at", "is_user_cancel"})

	if err == nil {
		if _, err := new(model.OrderNotes).SetOrderId(takeoutOrder.OrderId).SetType(2).
			SetContents(data.CancelReason).Insert(); err != nil {
			log.Error(errors.New(fmt.Sprintf("order id: %v, reason: %v", takeoutOrder.OrderId, data.CancelReason)))
			log.Verbose("ๆฐๅข็ทไธ้ป้ค่จๅฎๅๆถๅ่จปๅคฑๆ")
			log.Error(err)
		}
	}

	if err == nil {

		if takeoutOrder.OrderId > 0 {
			order.
				SetId(takeoutOrder.OrderId).
				SetCanceledAt(time.Now()).
				SetStatus(2).
				Update([]string{"canceled_at", "status"})
		}

		if takeoutOrderStatus == 0 {
			var couponUserIdAry []interface{}
			orderDiscount.SetOrderSn(takeoutOrder.OrderSn).QueryAll(func(rs *model.OrderDiscount) {
				couponUserIdAry = append(couponUserIdAry, rs.CouponUserId)
			})
			if len(couponUserIdAry) > 0 {
				couponUser.
					SetIdAry(couponUserIdAry).
					SetUserId(handler.UserId).
					SetStatus(0).
					SetUseStoreId(0).
					SetUpdatedAt(time.Now()).
					UpdateInAryById([]string{"status", "updated_at", "use_store_id"})
			}
		}

		fcm.SetBody("#" + takeoutOrder.TakeoutSerialNumber + " ่จๅฎ็ทจ่:" + takeoutOrder.OrderSn + " ๆถ่ฒป่ๅๆถ่จๅฎ").
			SetStoreId(takeoutOrder.StoreId).
			SetTitle("ๆ่จๅฎ่ขซๅๆถ๏ผ่ซๅๅพๅคๅธถๅค้็ขบ่ช").
			SetFcmServerKey(config.ServerInfo.FcmServerKey).
			SetData(map[string]interface{}{
				"action":    "takeoutOrder",
				"type":      "cancelOrder",
				"takeoutId": takeoutOrder.Id,
			}).
			SetManagementId(TakeoutManagementId).
			SendAdminStoreId()

		return util.RS{Message: "ๅๆถ่จๅฎๆๅ", Status: true}
	} else {
		log.Error(err)
		return util.RS{Message: "ๅๆถ่จๅฎๅคฑๆ", Status: false}
	}
}

func (handler *OrderHandler) AgainOrder() interface{} {
	var (
		data               ReceivedOrder
		onlineOrderSetting model.OnlineOrderSetting
		posMenuRepository  repository.PosMenuRepository
		posMenuTypeDetail  model.PosMenuTypeDetail
	)

	if err := json.Unmarshal([]byte(handler.Parameter), &data); err != nil {
		log.Error(err)
		return util.RS{Message: "", Status: false}
	} else if data.StoreId == 0 {
		log.Error(errors.New("data.StoreId == 0"))
		return util.RS{Message: "", Status: false}
	} else if len(data.OrderFoodList) == 0 {
		return util.RS{Message: "ๆฅ็ก้ค้ป่ณๆ", Status: false}
	}

	onlineOrderSetting.SetStoreId(data.StoreId).QueryOne()

	menuData := posMenuRepository.SetStoreId(data.StoreId).GetStoreMenuData()

	if !onlineOrderSetting.StrictStatus || menuData["id"].(int64) == 0 {
		return util.RS{Message: "็ฎๅ้่ฉฒๅบๅฎถ็ทไธ้ป้คๆๆฎต๏ผ็กๆณ่จ้คใ", Status: false}
	}

	orderFoodList := make([]map[string]interface{}, 0)
	orderFoodIdList := make([]interface{}, 0)
	orderFoodIdWhereList := make([]interface{}, 0)

	for _, orderFoodItem := range data.OrderFoodList {
		isRepeat := false
		for _, id := range orderFoodIdList {
			if id == orderFoodItem.FoodId {
				isRepeat = true
				break
			}
		}

		if !isRepeat {
			orderFoodIdList = append(orderFoodIdList, orderFoodItem.FoodId)
			orderFoodIdWhereList = append(orderFoodIdWhereList, "?")
			orderFoodList = append(orderFoodList, map[string]interface{}{
				"foodId":     orderFoodItem.FoodId,
				"name":       orderFoodItem.Name,
				"posAppName": orderFoodItem.PosAppName,
			})
		}
	}

	foodList := posMenuTypeDetail.
		SetMenuTypeIdList(menuData["menuTypeIdList"].([]interface{})).
		SetFoodIdList(orderFoodIdList).
		GetOnlineOrderFoodList()

	emptyFoodNameList := make([]string, 0)
	orderFoodDetailList := make([]map[string]interface{}, 0)

	for _, orderFoodItem := range orderFoodList {
		hasData := false
		for _, foodItem := range foodList {
			if orderFoodItem["foodId"].(int64) == foodItem["posFoodId"] {
				hasData = true
				break
			}
		}

		if hasData {
			var posFoodRepository repository.PosFoodRepository
			posFoodData := posFoodRepository.SetId(orderFoodItem["foodId"].(int64)).GetPosFoodDetail()
			orderFoodDetailList = append(orderFoodDetailList, posFoodData)
		} else {
			emptyFoodNameList = append(emptyFoodNameList, orderFoodItem["posAppName"].(string))
		}
	}

	message := ""

	if len(emptyFoodNameList) > 0 {
		message += strings.Join(emptyFoodNameList, ",") + "๏ผ็ฎๅ้ป้คๆๆฎตไธๆไพ่ฉฒ้ค้ป"
	}

	dataList := map[string]interface{}{
		"posMenuId":            menuData["id"],
		"orderFoodIdList":      orderFoodIdList,
		"orderFoodIdWhereList": orderFoodIdWhereList,
		"orderFoodDetailList":  orderFoodDetailList,
	}

	return util.RS{Message: message, Status: true, Data: dataList}

}

func (handler *OrderHandler) UseGoldFlowCreateOrder(orderTmpId int64) interface{} {
	var (
		orderTmp    model.OrderTmp
		defaultTime time.Time
	)

	orderTmp.SetId(orderTmpId).QueryOne()

	if orderTmp.JsonContent == "" {
		return util.RS{Message: "ๆฅ็ก่จๅฎ่ณๆ", Status: false}
	} else if orderTmp.PaymentAt != defaultTime {
		return util.RS{Message: "ๅทฒไปๆฌพ๏ผ็กๆณ้่คไปๆฌพ", Status: false}
	}

	var mapResult map[string]interface{}
	err := json.Unmarshal([]byte(orderTmp.JsonContent), &mapResult)
	if err != nil {
		fmt.Println("JsonToMapDemo err: ", err)
	}

	mapResult["paymentNote"] = ""
	mapResult["cashNote"] = ""
	mapResult["cardNumber"] = "1234"
	mapResult["cardType"] = -1

	_, err = kafka.Push("OnlineCreateOrder", mapResult)

	dataList := make(map[string]interface{}, 0)
	//kafka.Listen(strconv.FormatInt(int64(mapResult["storeId"].(float64)), 10)+"OnlineAppApiCreateOrder", mapResult["uuid"].(string), func(kafKaRS util.KafKaRS) {
	//	dataList["status"] = kafKaRS.Status
	//	dataList["message"] = kafKaRS.Message
	//	dataList["data"] = kafKaRS.Data
	//})

	for {
		if kafKaRS := <-kafka.KafkaChan["OnlineAppApiCreateOrder"]; kafKaRS.Uuid == mapResult["uuid"].(string) {
			dataList["status"] = kafKaRS.Status
			dataList["message"] = kafKaRS.Message
			dataList["data"] = kafKaRS.Data
			break
		}
	}

	fmt.Println(dataList)

	if dataList["status"] == true {
		orderData := dataList["data"].(map[string]interface{})

		orderTmp.
			SetPaymentAt(time.Now()).
			SetOrderId(int64(orderData["orderId"].(float64))).
			Update([]string{"payment_at", "order_id"})
	} else {
		//ๅฏ่ฝ้ๅ้ๆฌพๅไฝ
	}

	return util.RS{Message: dataList["message"].(string), Status: dataList["status"].(bool), Data: dataList["data"].(map[string]interface{})}
}

func (handler *OrderHandler) GetOrderSnByTmpId() interface{} {
	var (
		data     ReceivedOrder
		orderTmp model.OrderTmp
	)

	if err := json.Unmarshal([]byte(handler.Parameter), &data); err != nil {
		return util.RS{Message: "", Status: false}
	} else {
		orderTmp.SetId(data.OrderTmpId).GetOrderData()

		if orderTmp.OrderSn != "" {
			return util.RS{Message: "", Status: true, Data: orderTmp.OrderSn}
		} else {
			return util.RS{Message: "", Status: false}
		}
	}
}
