package notification

import (
	"dispatcher/helper"
	"dispatcher/lib/database"
	model "dispatcher/model"
	notify "dispatcher/model/notification"
	"dispatcher/service"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/isayme/go-amqp-reconnect/rabbitmq"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func NotifyUserService(db, dbUser *sqlx.DB, channel *rabbitmq.Channel, rdsClient *redis.Client, order model.Order, rmqUri, messagesExpired,
	urlPolyline, urlRefund, fcmKey, ApiToken string, versionCode int, rmqErr chan *amqp.Error) {
	bodyUser := notify.NewNotifyUser(order)

	if *order.StateOrder == 3 {
		bodyUser.StatusDriver = 0
		bodyUser.StateOrder = "01"
		bodyUser.Type = 3

		if *order.JenisLayananID == 1 {
			bodyUser.OrderStatus = "Driver menuju lokasi penjemputan"

		} else {
			bodyUser.OrderStatus = "Driver menuju " + *order.AlamatPenjemput
		}

		bodyDirectionCoreRequest := map[string]interface{}{
			"order_id":    order.OrderID,
			"id_merchant": order.MerchantID,
			"step":        1,
		}

		jsonBody, err := json.Marshal(bodyDirectionCoreRequest)
		if err != nil {
			logrus.Error(err)
		}

		result := helper.GeneratePolyline(jsonBody, ApiToken, urlPolyline)

		// Get Profile Picture
		bodyUser.DriverInfo.FotoProfil = database.GetDriverProfilePicture(db, order.DriverID)

		segmentedPolyline := helper.GetSegmentedPolyline(result.RespCore.ResponseData.OverviewPolyline)
		fmt.Println("poly line || " + result.RespCore.ResponseData.OverviewPolyline)
		fmt.Println("polyline segmented || " + segmentedPolyline)
		fmt.Println("driver || " + result.RespCore.ResponseData.DriverLocation)
		bodyUser.InfoRoute.Polyline = segmentedPolyline
		bodyUser.InfoRoute.DriverLocation = result.RespCore.ResponseData.DriverLocation
		database.InsertPolylineCore(bodyUser.OrderID, *order.StateOrder, segmentedPolyline, db)

	} else if *order.StateOrder == 4 {
		bodyUser.DriverInfo.FotoProfil = database.GetDriverProfilePicture(db, order.DriverID)
		bodyUser.StatusDriver = 0

		fmt.Printf("type: %v \n", bodyUser.Type)
		fmt.Printf("sub_state: %v \n", order.SubState)
		fmt.Printf("step: %v \n", order.Step)
		if *order.JenisLayananID == 1 {
			bodyUser.StateOrder = "03"

		} else if *order.JenisLayananID == 4 {
			bodyUser.StateOrder = "01"
			if order.Step == 2 {
				bodyUser.StateOrder = "03"
			}
		}

		bodyUser.Type = 4
		if *order.JenisLayananID == 4 && (order.SubState == 2 || order.SubState == 3) {
			bodyDirectionCoreRequest := map[string]interface{}{
				"order_id":    order.OrderID,
				"id_merchant": order.MerchantID,
				"step":        order.Step,
			}

			jsonBody, err := json.Marshal(bodyDirectionCoreRequest)
			if err != nil {
				logrus.Error(err)
			}

			result := helper.GeneratePolyline(jsonBody, ApiToken, urlPolyline)
			segmentedPolyline := helper.GetSegmentedPolyline(result.RespCore.ResponseData.OverviewPolyline)

			fmt.Println("polyline || " + result.RespCore.ResponseData.OverviewPolyline)
			fmt.Println("polyline segmented || " + segmentedPolyline)
			fmt.Println("driver || " + result.RespCore.ResponseData.DriverLocation)
			bodyUser.InfoRoute.Polyline = segmentedPolyline
			bodyUser.InfoRoute.DriverLocation = result.RespCore.ResponseData.DriverLocation
			database.InsertPolylineCore(bodyUser.OrderID, *order.StateOrder, segmentedPolyline, db)

		} else if *order.JenisLayananID == 1 {
			bodyDirectionCoreRequset := map[string]interface{}{
				"order_id": order.OrderID,
				"step":     order.Step,
			}

			jsonBody, err := json.Marshal(bodyDirectionCoreRequset)
			if err != nil {
				logrus.Error("error when marshaling ", err)
			}
			result := helper.GeneratePolyline(jsonBody, ApiToken, urlPolyline)
			segmentedPolyline := helper.GetSegmentedPolyline(result.RespCore.ResponseData.OverviewPolyline)

			fmt.Println("polyline || " + result.RespCore.ResponseData.OverviewPolyline)
			fmt.Println("polyline segmented || " + segmentedPolyline)
			fmt.Println("driver || " + result.RespCore.ResponseData.DriverLocation)
			bodyUser.InfoRoute.Polyline = segmentedPolyline
			bodyUser.InfoRoute.DriverLocation = result.RespCore.ResponseData.DriverLocation
			database.InsertPolylineCore(bodyUser.OrderID, *order.StateOrder, segmentedPolyline, db)

		} else {
			bodyUser.InfoRoute.Polyline = ""
			location, err := database.GetLocationDriver(order.DriverID, db)
			if err != nil {
				logrus.Error("any error when get latlon driver ", err)
			}

			bodyUser.InfoRoute.DriverLocation = location
		}

		err := service.EncodePolylineOrder(order.OrderID, 3)
		if err != nil {
			logrus.Error(fmt.Sprintf("got an error when encoding polyline :: %v || Order ID :: %v", err, order.OrderID))
		}

	} else if *order.StateOrder == 6 {
		bodyUser.StateOrder = "06"
		bodyUser.OrderStatus = "Menunggu Pembayaran"
		bodyUser.Type = 1
		bodyUser.PembayaranStatus = 0

	} else if *order.StateOrder == 2 {
		bodyUser.StateOrder = "05"
		bodyUser.OrderStatus = "Sedang mencari driver untukmu"

	} else if *order.StateOrder == 5 {
		bodyUser.StatusDriver = 1
		bodyUser.Type = 1
		bodyUser.StateOrder = "X1"
		bodyUser.OrderStatus = "Maaf, saat ini tidak ditemukan driver di sekitarmu"
	}
}
