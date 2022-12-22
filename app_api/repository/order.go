package repository

import (
	"app_api/model"
	"app_api/util/hit"
	"strconv"
	"time"
)

type OrderRepository struct{}
type jsonNonStationary struct {
	Money  string `json:"money"`
	Number string `json:"number"`
}

func (repository *OrderRepository) QueryOrder(userId int64, orderSn string) *model.Orders {
	var orders model.Orders
	return orders.SetUserId(userId).SetOrderSn(orderSn).QueryOne()
}

func (repository *OrderRepository) GetEstimatedTime(estimatedTime time.Time) map[string]interface{} {
	WeekDayMap := map[string]string{
		"Monday":    "星期一",
		"Tuesday":   "星期二",
		"Wednesday": "星期三",
		"Thursday":  "星期四",
		"Friday":    "星期五",
		"Saturday":  "星期六",
		"Sunday":    "星期日",
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

	mapResult := map[string]interface{}{
		"day":  estimatedTime.Format("2006-01-02"),
		"time": timeStart + "~" + timeEnd,
		"week": WeekDayMap[estimatedTime.Weekday().String()],
	}

	return mapResult
}
