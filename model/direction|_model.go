package model

type RespDataDetail struct {
	OverviewPolyline string `json:"overview_polyline"`
	DriverLocation   string `json:"driver_location"`
}

type RespData struct {
	ResponseCode      string         `json:"response_code"`
	ResponseDeskripsi string         `json:"response_deskripsi"`
	ResponseData      RespDataDetail `json:"response_data"`
}

type Data struct {
	RespCore RespData `json:"data"`
}
