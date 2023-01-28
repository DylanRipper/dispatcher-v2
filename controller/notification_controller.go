package controller

import (
	"dispatcher/helper"
	"dispatcher/lib/database"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func NotifyUserController(db *sqlx.DB) {
	listOrderNotify, err := database.GetOrderSendNotify(db)
	if err != nil {
		logrus.Error(err)
	}

	if len(listOrderNotify) > 0 {
		// isAnyData := true
		for _, orderNotify := range listOrderNotify {
			fmt.Printf("test notif :%v || status :%v \n", orderNotify.OrderID, *orderNotify.StateOrder)

			// if state order = 2, do calculate ETA
			var etaStr string
			if *orderNotify.StateOrder == 2 {
				orderDetail, err := database.GetOrderDetail(orderNotify.OrderID, db)
				if err != nil {
					logrus.Error(err)
				}

				etaStr = helper.CalculateETA(orderDetail, *orderNotify.WaktuTempuh)
				orderNotify.ETA = etaStr
			}

			err := database.UpdateOrderSendNotify(*orderNotify.StateOrder, orderNotify.OrderID, etaStr, db)
			if err != nil {
				logrus.Error(err)
			}

			item := orderNotify
			go func() {
				if *item.StateOrder == 0 {

				}
			}()
		}

	}
}
