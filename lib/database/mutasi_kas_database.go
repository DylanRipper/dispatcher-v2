package database

import (
	"dispatcher/model"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func GetMutasiKas(orderID uint64, db *sqlx.DB) ([]model.LogMutasiKas, error) {
	var listMutasiKas []model.LogMutasiKas
	q := `
	select 
		order_id, 
		nominal, 
		coalesce(transaction_reference, '') as transaction_reference, 
		coalesce(reference_id, 0) as reference_id, 
		coalesce(cb_success, 0) as cb_success, 
		jenis_kas, cb_success, 
		id_kas
	from 
		bj_t_mutasi_kas
	where 
		order_id = $1`

	err := db.Select(&listMutasiKas, q, orderID)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return listMutasiKas, nil

}

func InsertMutasiKas(nominal float64, orderID uint64, db *sqlx.DB) error {
	var ID int
	q := `
	insert into 
		bj_t_mutasi_kas
		(uraian, nominal, id_kas, kas, order_id, tanggal, jam, jenis_kas) 
	values
		($1,$2,$3,$4,$5,$6,$7,$8)
	returning id`

	err := db.QueryRowx(q,
		"Refund Biaya Pesanan",
		nominal,
		7,
		"refund_biaya_pesanan_user",
		orderID,
		time.Now().Format("2006-01-02"),
		time.Now().Format("15:04:05"),
		"KREDIT",
	).Scan(&ID)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func UpdateMutasiKasRefund(ID int, respRefund model.RefundResponse, tx *sqlx.Tx) error {
	q := `
	update 
		bj_t_mutasi_kas 
	set 
		cb_success=:cb_success, 
		cb_message=:cb_message, 
		transaction_reference=:transaction_reference, 
		cb_response_code=:cb_response_code
	where 
		id=:id`

	_, err := tx.NamedExec(q, map[string]interface{}{
		"cb_success":            respRefund.Rc,
		"cb_message":            respRefund.Rd,
		"transaction_reference": respRefund.MiD,
		"cb_response_code":      strconv.Itoa(int(respRefund.Rc)),
		"id":                    ID,
	})

	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
