package model

type LogMutasiKas struct {
	OrderID              uint64  `db:"order_id"`
	KasID                uint64  `db:"id_kas"`
	CbSuccess            uint64  `db:"cb_success"`
	Nominal              float64 `db:"nominal"`
	TransactionReference string  `db:"transaction_reference"`
	ReferenceID          string  `db:"reference_id"`
	CbResponseCode       string  `db:"cb_response_code"`
	JenisKas             string  `db:"jenis_kas"`
}
