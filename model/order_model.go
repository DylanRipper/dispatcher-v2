package model

import "time"

type Order struct {
	OrderID                 uint64 `db:"order_id"`
	StatusBlacklistUser     uint64 `db:"status_blacklist_user"`
	SubState                uint64
	Step                    uint64
	DriverID                uint64   `db:"id_driver"`
	NftID                   uint64   `db:"id_nft"`
	PartnerID               *uint64  `db:"partner_id"`
	ReffUserPartnerID       *string  `db:"reff_id_user_partner"`
	ReffOrderPartnerID      *string  `db:"reff_id_order_partner"`
	BlacklistUserStatus     *uint64  `db:"status_blacklist_user"`
	JenisLayananID          *uint64  `db:"id_jenis_layanan"`
	JarakTempuh             *uint64  `db:"jarak_tempuh"`
	RadiusDispatch          *float64 `db:"radius_dispatch"`
	StateOrder              *int64   `db:"state_order"`
	JenisPembayaranID       *uint64  `db:"id_jenis_pembayaran"`
	WaktuTempuh             *uint64  `db:"waktu_tempuh"`
	KotaOriginID            *uint64  `db:"id_kota_origin"`
	MerchantID              *uint64  `db:"id_merchant"`
	DiskonVoucher           *int64   `db:"diskon_voucher"`
	IsKurir                 *uint64  `db:"is_kurir"`
	ModeKurir               *uint64  `db:"mode_kurir"`
	TotalJarakTempuh        *uint64  `db:"total_jarak_tempuh"`
	TotalDestinasiPaket     *uint64  `db:"total_destinasi_paket"`
	TotalBeratPaket         *uint64  `db:"total_berat_paket"`
	KomisiDriver            float64  `db:"komisi_driver"`
	PendapatanDriver        float64  `db:"pendapatan_driver"`
	BiayaBelanjaFinal       *float64 `db:"biaya_belanja_final"`
	WaitingTime             *float64 `db:"waiting_time"`
	LatPenjemputan          *float64 `db:"latitude_penjemputan"`
	LonPenjemputan          *float64 `db:"longitude_penjemputan"`
	NominalBiaya            *float64 `db:"nominal_biaya"`
	OngkosKirim             *float64 `db:"ongkos_kirim"`
	DIskonPromoUser         *float64 `db:"diskon_promo_user"`
	BiayaFoodMerchant       *float64 `db:"biaya_food_merchant"`
	TotalBayarTunaiMerchant *float64 `db:"total_bayar_tunai_merchant"`
	TitikPengantaran        *string  `db:"titik_pengantaran"`
	DetailLokasiPenjemputan *string  `db:"detail_lokasi_penjemputan"`
	DetailLokasiPengantaran *string  `db:"detail_lokasi_pengantaran"`
	TitikPenjemputan        *string  `db:"titik_penjemputan"`
	StatusOrder             *string  `db:"status_order"`
	DeskripsiStateOrder     *string  `db:"deskripsi_state_order"`
	NamaTempatPenjemput     *string  `db:"nama_tempat_penjemputan"`
	AlamatPenjemput         *string  `db:"alamat_penjemputan"`
	JamRilisPending         *string  `db:"jam_rilis_pending"`
	NamaTempatPengantaran   *string  `db:"nama_tempat_pengantaran"`
	AlamatPengantaran       *string  `db:"alamat_pengantaran"`
	NamaPengorder           *string  `db:"nama_pengorder"`
	NoHPPengorder           *string  `db:"no_hp_pengorder"`
	CatatanPengorder        *string  `db:"catatan_pengorder"`
	TipePembayaranMerchant  *string  `db:"tipe_pembayaran_merchant"`
	AutoBid                 *string  `db:"auto_bid"`
	NamaPengirim            *string  `db:"nama_pengirim"`
	NoHPPengirim            *string  `db:"no_hp_pengirim"`
	CatatanPengirim         *string  `db:"catatan_pengirim"`
	NamaPenerima            *string  `db:"nama_penerima"`
	NoHPPenerima            *string  `db:"no_hp_penerima"`
	DetailBarangKirim       *string  `db:"detail_barang_kiriman"`
	JamOrder                *string  `db:"jam_order"`
	StartTimeDriverSearch   *string  `db:"start_time_driver_search"`
	NamaDriver              *string  `db:"nama_driver"`
	NopolKendaranDriver     *string  `db:"nopol_kendaraan_driver"`
	JenisKendaraanDriver    *string  `db:"jenis_kendaraan_driver"`
	KategoriKEndaraanDriver *string  `db:"kategori_kendaraan_driver"`
	WarnaKendaraanDriver    *string  `db:"warna_kendaraan_driver"`
	AlasanPembatalanOrder   *string  `db:"alasan_pembatalan_order"`
	AlasanPembatalanDriver  *string  `db:"alasan_pembatalan_driver"`
	NoHpDriver              *string  `db:"no_hp_driver"`
	ETA                     string   `db:"eta"`
	Via                     *string  `db:"via"`
	VersionCode             *int     `db:"version_code"`
	PhotoProfileDriver      string
	TanggalOrder            *time.Time `db:"tanggal_order"`
	TanggalRilisPending     *time.Time `db:"tanggal_rilis_pending"`
}

type OrderDetail struct {
	MerchantID  uint64  `db:"id_merchant"`
	WaktuTempuh *uint64 `db:"waktu_tempuh"`
}

type StateOrder struct {
	IDOrder uint64 `json:"id_order"`
	State   int    `json:"state"`
}

func GetOrder(id uint64, state int) *StateOrder {
	return &StateOrder{
		IDOrder: id,
		State:   state,
	}
}
