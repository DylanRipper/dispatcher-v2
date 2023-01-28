package helper

import (
	"dispatcher/model"
	"time"
)

func CalculateETA(listOrderDetail []model.OrderDetail, totalSpendTime uint64) string {
	for _, OrderDetail := range listOrderDetail {
		totalSpendTime += *OrderDetail.WaktuTempuh
		totalSpendTime += 600
	}

	etaStr := time.Now().Add(time.Duration(totalSpendTime+600)*time.Second).Format("15:04") + "|" +
		time.Now().Add(time.Duration(totalSpendTime+1500)*time.Second).Format("15:04")
	return etaStr
}
