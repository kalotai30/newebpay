package handler

import (
	"encoding/json"
	. "trade/database/mysql"
	"trade/model"
	"trade/repository"
	"trade/util/kafka"
	"trade/util/log"
)

//產生訂單資料
func CreateOrder2(dataByte []byte) (isBreak bool) {
	var (
		data             model.ReceivedOrder
		orderDiscount    model.OrderDiscount
		orderCreditCard  model.OrderCreditCard
		orderPaymentTerm model.OrderPaymentTerm
		couponUser       model.CouponUser
		store            model.Store
		message          = ""
		status           = true
		saveOrderId      int64
		orderSn          string
		takeoutId        int64
		takeoutSnString  string
		createMenuOrder  repository.CreateMenuOrderRepository
	)

	if err := json.Unmarshal(dataByte, &data); err != nil {
		log.Error(err)
		message = "Error"
	} else {

		for _, couponList := range data.CouponList {
			var (
				userCoupon    model.CouponUser
				userCouponAry []interface{}
			)
			userCoupon.
				SetCouponId(couponList.CouponListId).
				SetStatusAry([]interface{}{0}).
				SetCouponTypeId(couponList.CouponTypeId).
				SetUserId(data.UserId).
				SetStartTime(couponList.StartTime).
				SetEndTime(couponList.EndTime).
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
			residue := 0
			for _, item2 := range userCouponAry {
				if item2.(map[string]interface{})["isUse"] != 0 && item2.(map[string]interface{})["isStart"] != 0 {
					residue += 1
				}
			}
			if data.IsEditOrder {
				var orderDiscount model.OrderDiscount
				orderDiscount.
					SetOrderSn(data.OrderSn).
					SetCouponListId(userCoupon.CouponListId).
					SetStartTime(userCoupon.StartTime).
					SetEndTime(userCoupon.EndTime).
					QueryAllJoinList(func(rs *model.OrderDiscount) {
						residue += 1
					})
			}
			if couponList.ChoiceAmount > residue {
				message = couponList.CouponListName + "數量不足"
				status = false
				break
			}
		}

		newOrder := createMenuOrder.CreateMenuOrder2(data)

		saveOrderId = newOrder["orderId"].(int64)
		orderSn = newOrder["orderSn"].(string)
		takeoutId = newOrder["takeoutId"].(int64)
		takeoutSnString = newOrder["takeoutSnString"].(string)
		status = newOrder["status"].(bool)
		message = newOrder["message"].(string)

		if status {
			NewDB().Transaction(func(database Database) {

				store.SetId(data.StoreId).QueryOne()

				//紀錄信用卡使用資訊
				//這邊在實際開始用時要再確認各種情況
				if data.CardType != 0 {
					_, err = orderCreditCard.SetOrderId(saveOrderId).SetAmount(data.InvoicePrice).SetCardNumber(data.CardNumber).SetPaymentNote(data.PaymentNote).
						SetCardType(data.CardType).CreateByTransaction(database)
					if err != nil {
						message = "紀錄訂單失敗"
						status = false
						log.Error(database.Rollback())
					}
				}

				//新增訂單折扣項目order_discount
				var couponUserIdAry []interface{}
				var resumCouponUserIdAry []interface{}
				if status {

					if data.IsEditOrder {
						var orderDiscount model.OrderDiscount
						orderDiscountIdAry := make([]interface{}, 0)
						orderDiscount.
							SetOrderSn(orderSn).
							QueryAll(func(rs *model.OrderDiscount) {
								orderDiscountIdAry = append(orderDiscountIdAry, rs.Id)
								resumCouponUserIdAry = append(resumCouponUserIdAry, rs.CouponUserId)
							})
						if len(orderDiscountIdAry) > 0 {
							orderDiscount.DeleteByIdAry(orderDiscountIdAry, database)
						}

						if len(resumCouponUserIdAry) > 0 {
							err2 := couponUser.
								SetStatus(0).
								SetUseStoreId(0).
								SetUseObject(0).
								SetUseObjectId(0).
								UpdateInAryByTransaction([]string{"status", "updated_at", "use_store_id", "use_object", "use_object_id"}, resumCouponUserIdAry, database)

							if err2 != nil {
								log.Error(err2)
								message = "更新使用優惠券失敗"
								status = false
								log.Error(database.Rollback())
								return
							}
						}

						err := database.Model(new(model.OrderPaymentTerm)).Where("order_id", "=", saveOrderId).Delete()
						if err != nil {
							log.Error(err)
							message = "刪除訂單支付方式失敗"
							status = false
							log.Error(database.Rollback())
							return
						}
					}
					for _, couponList := range data.CouponList {
						for i := 0; i < couponList.ChoiceAmount; i++ {
							orderDiscount.SetUserId(data.UserId).SetId(0).
								SetStoreId(data.StoreId).
								SetOrderSn(orderSn).
								SetCouponUserId(couponList.CouponUserIdAry[i]).
								SetCouponListId(couponList.CouponListId).
								SetCouponName(couponList.CouponListName)

							if len(couponList.DiscountComputedAry) == couponList.ChoiceAmount {
								couponList.DiscountComputed = couponList.DiscountComputedAry[i]
							}
							if couponList.DiscountComputed != 0 {
								orderDiscount.SetCouponDiscount(couponList.DiscountComputed)
							} else {
								orderDiscount.SetCouponDiscount(couponList.Amount)
							}

							_, err := orderDiscount.SetCouponType(couponList.CouponType).InsertByTransaction(database)

							if err != nil {
								log.Error(err)
								message = "新增訂單折扣失敗"
								status = false
								log.Error(database.Rollback())
								return
							}
							//將coupon_user_id組成陣列(更新消費者優惠券已使用)
							couponUserIdAry = append(couponUserIdAry, couponList.CouponUserIdAry[i])
						}
					}
					if len(couponUserIdAry) > 0 {
						//更新消費者優惠券為已使用
						err := couponUser.
							SetStatus(1).
							SetUseStoreId(data.StoreId).
							SetUseObject(ObjectTypeStore).
							SetUseObjectId(data.StoreId).
							UpdateInAryByTransaction([]string{"status", "updated_at", "use_store_id", "use_object", "use_object_id"}, couponUserIdAry, database)
						if err != nil {
							log.Error(err)
							message = "更新使用優惠券失敗"
							status = false
							log.Error(database.Rollback())
							return
						}
					}
				}

				//計算現金
				var cashAmount = data.InvoicePrice
				//新增訂單支付方式order_payment_term
				if saveOrderId != 0 && status { //先確認訂單已成功新增
					//並且如果選擇的是會立即付款的方式時要把 takeout_order is_checked_out 寫 1
					if data.PaymentType == 1 {
						_, err = orderPaymentTerm.SetOrderId(saveOrderId).
							SetName("現金").
							SetPaymentMethodId(0).
							SetNote(data.CashNote).
							SetAmount(0).
							SetPrice(int64(cashAmount)).
							SetIsCoupon(0).CreateByTransaction(database)
					} else {
						var posPaymentMethod model.PosPaymentMethod
						var paymentName string

						if data.PaymentType == 2 || data.PaymentType == 8 {
							paymentName = "信用卡"
						}

						posPaymentMethod.
							SetStoreId(data.StoreId).
							SetType(3).
							QueryOne()

						_, err = orderPaymentTerm.SetOrderId(saveOrderId).
							SetName(paymentName).
							SetPaymentMethodId(posPaymentMethod.Id).
							SetNote(data.CashNote).
							SetAmount(0).
							SetPrice(int64(cashAmount)).
							SetIsCoupon(0).CreateByTransaction(database)
					}

					if err != nil {
						log.Error(err)
						message = "新增訂單支付方式失敗"
						status = false
						log.Error(database.Rollback())
						return
					}
				}

				if status {
					log.Error(database.Commit())
				} else {
					log.Error(database.Rollback())
				}
			})
		}
	}

	sendData := map[string]interface{}{
		"orderId":         saveOrderId,
		"orderSn":         orderSn,
		"takeoutId":       takeoutId,
		"totalPrice":      data.TotalPrice,
		"takeoutSnString": takeoutSnString,
	}

	if _, err := kafka.Push("OnlineAppApiCreateOrder", map[string]interface{}{
		"uuid":    data.Uuid,
		"status":  status,
		"message": message,
		"data":    sendData,
	}); err != nil {
		log.Error(err)
	}

	return false
}
