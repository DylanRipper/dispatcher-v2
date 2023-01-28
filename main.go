package main

import "dispatcher/config"

func main() {
	rmqChannel, _, db, rdsClient := config.InitDB()
	defer func() {
		if rmqChannel != nil {
			_ = rmqChannel.Close()
		}
	}()

	defer func() {
		if db != nil {
			sql, _ := db.DB()
			sql.Close()
		}
	}()

	defer func() {
		if rdsClient != nil {
			_ = rdsClient.Close()
		}
	}()
}
