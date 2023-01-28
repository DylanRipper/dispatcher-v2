package database

import (
	"dispatcher/model"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func GetOrder(timeNowStr string, config model.Config, db *sqlx.DB) ([]model.Order, error) {
	startTime := `cast(coalesce(start_time_driver_search,concat(tanggal_rilis_pending,' ',jam_rilis_pending)) as timestamp)`
	endTime := `cast('` + timeNowStr + `' as timestamp)`
	interval := `EXTRACT(EPOCH FROM(` + endTime + `-` + startTime + `))`
	batasWaktuTunggu := "21"
	queryWaitingTime := interval + " > " + batasWaktuTunggu
	queryWaitingTimeState0 := interval + " > 0 "
	var resultOrder []model.Order
	q := `
	select
		order_id,
		partner_id,
		get_status_blacklist_user(reff_id_user_partner::bigint) as status_blacklist_user,
		reff_id_user_partner,
		reff_id_order_partner,
		id_jenis_layanan,
		jarak_tempuh,
		status_order,
		state_order,
		deskripsi_state_order,
		id_jenis_pembayaran,
		waktu_tempuh,
		id_kota_origin,
		id_merchant,
		coalesce(id_nft, 0) as id_nft,
		coalesce(diskon_voucher, 0) as diskon_voucher,
		coalesce(radius_dispatch, 0) as radius_dispatch,
		is_kurir,
		mode_kurir,
		total_jarak_tempuh,
		total_destinasi_paket,
		total_berat_paket,
		latitude_penjemputan,
		longitude_penjemputan,
		titik_penjemputan,
		titik_pengantaran,
		detail_lokasi_penjemputan,
		detail_lokasi_pengantaran,
		nominal_biaya,
		ongkos_kirim,
		diskon_promo_user,
		biaya_food_merchant,
		nama_tempat_penjemputan,
		alamat_penjemputan,
		tanggal_rilis_pending,
		jam_rilis_pending,
		nama_tempat_pengantaran,
		alamat_pengantaran,
		nama_pengorder,
		no_hp_pengorder,
		catatan_pengorder,
		tipe_pembayaran_merchant,
		nama_pengirim,
		no_hp_pengirim,
		catatan_pengirim,
		nama_penerima,
		no_hp_penerima,
		detail_barang_kiriman,
		jam_order,
		biaya_belanja_final,
		coalesce(total_bayar_tunai_merchant, 0) as total_bayar_tunai_merchant,
		tanggal_rilis_pending, ` + interval + ` as waiting_time
	from 
		bj_t_order
	where ((state_order = 0 and is_pending=0) 
		or (state_order = 0 and is_pending=1 and ` + queryWaitingTimeState0 + `) 
		or (state_order= 2 and ` + queryWaitingTime + `))
	order by order_id desc
	limit $1`
	err := db.Select(&resultOrder, q, config.AppConfig.LimitGetOrder)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return resultOrder, nil
}

func GetOrderDetail(orderID uint64, db *sqlx.DB) ([]model.OrderDetail, error) {
	var listOrderDetail []model.OrderDetail
	q := `
	select
		distinct id_merchant,
		coalesce(waktu_tempuh_pengiriman, 0) as waktu_tempuh
	from
		bj_t_detail_order
	where 
		id_order = $1`
	err := db.Select(&listOrderDetail, q, orderID)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return listOrderDetail, nil

}

func UpdateStateOrderSearchDriver(timeNowStr string, orderID uint64, db *sqlx.DB) error {
	q := `
	update 
		bj_t_order
	set 
		state_order = :state_order,
		deskripsi_state_order = 'Sedang mencari driver',
		start_time_driver_search = :start_time_driver_search 
	where 
		order_id = :order_id 
		and (state_order = 0 or state_order = 2)`

	_, err := db.NamedExec(q, map[string]interface{}{
		"order_id":                 orderID,
		"state_order":              2,
		"start_time_driver_search": timeNowStr,
	})

	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func UpdateETAOrder(eta string, orderID uint64, db *sqlx.DB) error {
	q := `
    update 
        bj_t_order 
    set
		eta = :eta
	where
	    order_id = :order_id`

	_, err := db.NamedExec(q, map[string]interface{}{
		"eta":      eta,
		"order_id": orderID,
	})

	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func InsertPolylineCore(idOrder uint64, state int64, polyline string, db *sqlx.DB) {
	times := time.Now().Format("15:04:05")
	date := time.Now().Format("2006-01-02")
	q := `update 
		public.bj_t_ongoing_order 
	set 
		polyline = :polyline,
		time_update = :time_update,
		date_update = :date_update
	where 
		id_order = :idOrder 
		and state_order = :state`
	_, err := db.NamedExec(q, map[string]interface{}{
		"polyline":    polyline,
		"time_update": times,
		"date_update": date,
		"idOrder":     idOrder,
		"state":       state,
	})
	if err != nil {
		logrus.Error(err)
	}
}
