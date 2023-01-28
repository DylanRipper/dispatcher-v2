package model

import (
	"dispatcher/model"
	"strings"
)

type NotifOrderUser struct {
	OrderID          uint64       `json:"order_id"`
	Duration         uint64       `json:"duration"`
	IDJenisLayanan   uint64       `json:"id_jenis_layanan"`
	StateOrder       string       `json:"state_order"`
	OrderStatus      string       `json:"order_status"`
	NotifType        string       `json:"notif_type"`
	PembayaranStatus int64        `json:"status_bayar"`
	Type             uint64       `json:"type"`
	StatusDriver     uint64       `json:"status_driver"`
	Unread           uint64       `json:"unread"`
	DriverInfo       DriverDetail `json:"driver_info"`
	InfoRoute        RouteDetail  `json:"info_route"`
	TotalJarakTempuh uint64       `json:"total_jarak_tempuh"`
	CancelDetail     string       `json:"cancel_detail"`
	IsEtaChange      uint64       `json:"is_eta_change"`
	ETAMin           string       `json:"eta_min"`
	ETAMax           string       `json:"eta_max"`
	WaktuOrder       string       `json:"waktu_order"`
}

type DriverDetail struct {
	DriverID          uint64 `json:"id_driver"`
	NamaDriver        string `json:"nama_driver"`
	Nopol             string `json:"nopol"`
	KategoriKendaraan string `json:"kategori_kendaraan"`
	JenisKendaraan    string `json:"jenis_kendaraan"`
	WarnaKendaraan    string `json:"warna_kendaraan"`
	NoHp              string `json:"no_hp"`
	FotoProfil        string `json:"foto_profil"`
}

type RouteDetail struct {
	Polyline       string `json:"polyline"`
	DriverLocation string `json:"driver_location"`
}

func NewNotifyUser(o model.Order) NotifOrderUser {
	driverDetail := DriverDetail{
		DriverID:          o.DriverID,
		NamaDriver:        *o.NamaDriver,
		Nopol:             *o.NopolKendaranDriver,
		KategoriKendaraan: *o.KategoriKEndaraanDriver,
		JenisKendaraan:    *o.JenisKendaraanDriver,
		WarnaKendaraan:    *o.WarnaKendaraanDriver,
		NoHp:              *o.NoHpDriver,
		FotoProfil:        o.PhotoProfileDriver,
	}

	route := RouteDetail{
		Polyline:       "",
		DriverLocation: "",
	}

	if o.SubState == 1 {
		route.DriverLocation = *o.TitikPenjemputan
	}

	notifOrderUser := NotifOrderUser{
		OrderID:          o.OrderID,
		OrderStatus:      *o.DeskripsiStateOrder,
		IDJenisLayanan:   *o.JenisLayananID,
		NotifType:        "notif",
		Type:             0,
		StatusDriver:     0,
		PembayaranStatus: -1,
		Unread:           0,
		DriverInfo:       driverDetail,
		InfoRoute:        route,
		TotalJarakTempuh: *o.JarakTempuh,
		CancelDetail:     *o.AlasanPembatalanDriver,
		WaktuOrder:       o.TanggalRilisPending.Format("2006-01-02") + " " + *o.JamRilisPending,
	}

	listEta := strings.Split(o.ETA, "|")
	if len(listEta) == 2 {
		notifOrderUser.ETAMin = listEta[0]
		notifOrderUser.ETAMax = listEta[1]
	}

	return notifOrderUser
}
