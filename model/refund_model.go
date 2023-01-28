package model

type Refund struct {
	OrderID       uint64  `json:"order_id"`
	Amount        float64 `json:"amount"`
	UserID        uint64  `json:"id_user"`
	TransactionID string  `json:"transaction_id"`
	Via           string  `json:"via"`
	VersionName   string  `json:"version_name"`
	VersionCode   string  `json:"version_code"`
}

type RefundResponse struct {
	Rc           uint64 `json:"rc"`
	Rd           string `json:"rd"`
	MiD          string `json:"mid"`
	ResponseTime string `json:"response_time"`
}
