package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func GetLocationDriver(idDriver uint64, db *sqlx.DB) (string, error) {
	var location string
	q := `
	select 
		concat_ws(',', latitude, longitude) as loc 
	from 
		bj_m_driver bmd 
	where 
		id_driver = $1`
	err := db.Get(&location, q, idDriver)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return location, err
}

func GetDriverProfilePicture(db *sqlx.DB, driverID uint64) (profilePicture string) {
	q := `
	select 
		coalesce(foto_profil, '') as foto_profil 
	from 
		bj_m_driver 
	where 
		id_driver = $1`
	err := db.Get(&profilePicture, q, driverID)
	if err != nil {
		logrus.Error(err)
	}

	return
}
