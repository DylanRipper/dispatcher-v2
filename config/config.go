package config

import (
	"dispatcher/model"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/isayme/go-amqp-reconnect/rabbitmq"
	"github.com/jmoiron/sqlx"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initLogger(config model.Config) {
	writer, err := rotatelogs.New(
		config.Logger.Path+".%Y%m%d",
		rotatelogs.WithLinkName(config.Logger.Path),
		rotatelogs.WithMaxAge(time.Duration(config.Logger.MaxAge*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(config.Logger.RotationTime*24)*time.Hour),
	)
	if err != nil {
		panic(err)
	}

	customFormatter := &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     true,
	}
	logrus.SetFormatter(customFormatter)
	logrus.SetReportCaller(true)
	logrus.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.WarnLevel:  writer,
			logrus.ErrorLevel: writer,
			logrus.FatalLevel: writer,
		},
		customFormatter,
	))
}

func initRabbitMQ(config model.Config) (*rabbitmq.Channel, chan *amqp.Error) {
	switch config.MessageBroker.Driver {
	case "rabbitmq":
		conn, err := rabbitmq.Dial(config.MessageBroker.Uri)
		if err != nil {
			logrus.Error(fmt.Errorf("dial rabbitMQ error: %s", err))
			panic(err)
		}
		notify := conn.NotifyClose(make(chan *amqp.Error))
		go func() {
			<-conn.NotifyClose(make(chan *amqp.Error)) //Listen to NotifyClose
			config.MessageBroker.Err <- errors.New("connection closed")
		}()
		channel, err := conn.Channel()
		if err != nil {
			logrus.Error(fmt.Errorf("channel: %s", err))
			panic(err)
		}
		return channel, notify

	default:
		panic(fmt.Sprintf("driver msgbroker %s not implemented", config.MessageBroker.Driver))
	}

}

func initRedis(config model.Config) *redis.Client {
	switch config.Cache.Driver {
	case "redis":
		rdsClient := redis.NewClient(&redis.Options{
			Addr:     config.Cache.Addr,
			Password: config.Cache.Password,
			DB:       config.Cache.Database,
			PoolSize: config.Cache.PoolSize,
		})
		return rdsClient

	default:
		panic(fmt.Sprintf("driver idempotence %s not implemented", config.Cache.Driver))

	}

}

func InitDB() (*rabbitmq.Channel, chan *amqp.Error, *gorm.DB, *redis.Client) {
	viper.SetConfigFile("config-dev.yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var config model.Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	// ===== Logger ======
	initLogger(config)
	os.Setenv("ID", "Asia/Jakarta")

	// ===== RabbitMQ ======
	rmqChannel, rmqError := initRabbitMQ(config)

	// ===== Redis Cache ======
	rdsClient := initRedis(config)

	// ===== Database ======
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=public",
		config.DB.Host,
		config.DB.Port,
		config.DB.Username,
		config.DB.Password,
		config.DB.DbName)

	sqlx, err := sqlx.Open("pgx", dsn)
	if err != nil {
		logrus.Error(err)
		panic(err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlx,
	}), &gorm.Config{})

	if err != nil {
		logrus.Error(err)
		panic(err)
	} else {
		logrus.Println("Connected to database")
	}

	return rmqChannel, rmqError, db, rdsClient
}
