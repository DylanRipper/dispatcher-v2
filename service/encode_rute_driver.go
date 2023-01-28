package service

import (
	"bytes"
	"dispatcher/model"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

func EncodePolylineOrder(id uint64, state int) error {
	// idStr := strconv.Itoa(id)
	data := model.GetOrder(id, state)
	request, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", "http://localhost:8080/encode-polyline-onduty", bytes.NewBuffer(request))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
	}

	if resp.StatusCode != 200 {
		return errors.New("Error encode polyline on duty:" + strconv.Itoa(resp.StatusCode) + "and status" + resp.Status)
	}

	return nil
}
