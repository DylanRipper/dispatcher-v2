package database

import (
	"dispatcher/model"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func GetOrderSendNotify(db *sqlx.DB) ([]model.Order, error) {
	var listOrderNotify []model.Order
	qGetNotify := `
		select
			order_id,
			state_order,
			id_merchant,
			deskripsi_state_order,
			start_time_driver_search,
			reff_id_user_partner,
			coalesce(id_driver, 0) as id_driver,
			coalesce(nama_driver, '') as nama_driver,
			coalesce(nopol_kendaraan_driver, '') as nopol_kendaraan_driver,
			coalesce(jenis_kendaraan_driver, '') as jenis_kendaraan_driver,
			coalesce(kategori_kendaraan_driver, '') as kategori_kendaraan_driver,
			coalesce(warna_kendaraan_driver, '') as warna_kendaraan_driver,
			coalesce(alasan_pembatalan_driver, '') as alasan_pembatalan_driver,
			coalesce(no_hp_driver, '')  as no_hp_driver,
			coalesce(eta, '')  as eta,
			coalesce(nama_tempat_penjemputan, '')  as nama_tempat_penjemputan,
			coalesce(waktu_tempuh, 0) as waktu_tempuh,
			coalesce(via, '') as via,
			coalesce(version_code, '0') as version_code,
			jarak_tempuh,
			coalesce(id_jenis_layanan, 0) as id_jenis_layanan,
			coalesce(id_jenis_pembayaran, 0) as id_jenis_pembayaran,
			no_hp_pengorder,
			tanggal_rilis_pending,
			jam_rilis_pending
		from 
			bj_t_order
		where
			status_send_notif = 0
		limit 20 `

	err := db.Select(&listOrderNotify, qGetNotify)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return listOrderNotify, nil
}

func UpdateOrderSendNotify(state int64, orderID uint64, eta string, db *sqlx.DB) error {
	q := `update bj_t_order
		set status_send_notif = 1`
	if state == 2 {
		q += `, eta = :eta`
	}

	q += ` where order_id = :order_id`
	_, err := db.NamedExec(q, map[string]interface{}{
		"order_id": orderID,
		"eta":      eta,
	})

	if err != nil {
		return err
	}
	return nil
}

func InsertLog(orderID, merchantID uint64, corrID, dataFcm string, db *sqlx.DB) error {
	q := `
 insert into 
 	log_fcm(tanggal,
		jam,
		status_sending,
		multicast_id,
		message_id,
		order_id,
		id_driver,
		id_merchant,
		is_received,
		data_fcm)
 values (:tanggal,
	:jam,
	:status_sending,
	:multicast_id,
	:message_id,
	:order_id,
	0,
	:id_merchant,
	false,
	:data_fcm)`

	_, err := db.NamedExec(q, map[string]interface{}{
		"tanggal":        time.Now().Format("2006-01-02"),
		"jam":            time.Now().Format("15:04:05"),
		"status_sending": 1,
		"multicast_id":   "merchant notif using rabbitmq",
		"message_id":     corrID,
		"order_id":       orderID,
		"id_merchant":    merchantID,
		"data_fcm":       dataFcm,
	})

	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}
