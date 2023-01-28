package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func GetTokenFCM(userID uint64, db *sqlx.DB) (string, error) {
	q := `
	select 
		coalesce(fcm_token, '') as fcm_token
	from 
        aci_users.users
	where
        id_user = $1
	`
	var token string
	err := db.Get(&token, q, userID)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	return token, nil
}
