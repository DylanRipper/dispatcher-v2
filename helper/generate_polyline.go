package helper

import (
	"bytes"
	"dispatcher/model"
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/twpayne/go-polyline"
)

func GeneratePolyline(reqBody []byte, token string, urlPolyline string) model.Data {
	req, err := http.NewRequest("POST", urlPolyline, bytes.NewBuffer(reqBody))
	if err != nil {
		logrus.Error(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-token", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
	}

	var resultResponse model.Data
	if err := json.Unmarshal(body, &resultResponse); err != nil {
		logrus.Error(err)
	}

	return resultResponse

}

func GetSegmentedPolyline(generatedPolyline string) string {
	buf := []byte(generatedPolyline)
	coords, _, _ := polyline.DecodeCoords(buf)
	//fmt.Println("slice II :: ", coords)

	var newPolyline [][]float64
	iter := 0
	for {
		newPolyline = [][]float64{}
		for i := 0; i < len(coords)-1; i++ {
			dist := CalcDistance(coords[i][0], coords[i][1], coords[i+1][0], coords[i+1][1], "K")
			newPolyline = append(newPolyline, []float64{coords[i][0], coords[i][1]})

			if dist > 0.02 {
				lat3, lon3 := getMiddlePoint(coords[i][0], coords[i][1], coords[i+1][0], coords[i+1][1])
				newPolyline = append(newPolyline, []float64{lat3, lon3})
			}
		}
		newPolyline = append(newPolyline, []float64{coords[len(coords)-1][0], coords[len(coords)-1][1]})
		coords = newPolyline

		if iter > 4 {
			break
		}
		iter++
	}

	//fmt.Println("slice :: ", newPolyline)
	encodedPoly := polyline.EncodeCoords(newPolyline)
	//fmt.Println("polyline :: " + string(encodedPoly))

	return string(encodedPoly)
}

func getMiddlePoint(currlat, currlon, targetlat, targetlon float64) (float64, float64) {
	lat1 := currlat * math.Pi / 180.0
	lat2 := targetlat * math.Pi / 180.0

	lon1 := currlon * math.Pi / 180.0
	dLon := (targetlon - currlon) * math.Pi / 180.0

	bx := math.Cos(lat2) * math.Cos(dLon)
	by := math.Cos(lat2) * math.Sin(dLon)

	lat3Rad := math.Atan2(
		math.Sin(lat1)+math.Sin(lat2),
		math.Sqrt(math.Pow(math.Cos(lat1)+bx, 2)+math.Pow(by, 2)),
	)
	lon3Rad := lon1 + math.Atan2(by, math.Cos(lat1)+bx)

	lat3 := lat3Rad * 180.0 / math.Pi
	lon3 := lon3Rad * 180.0 / math.Pi

	return lat3, lon3
}

func CalcDistance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * lat1 / 180)
	radlat2 := float64(PI * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}

	return dist
}
