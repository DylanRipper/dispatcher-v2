package model

type Config struct {
	Server struct {
		Rest string `validate:"required" mapstructure:"rest"`
	} `mapstructure:"server"`
	ServerCore struct {
		Url         string `validate:"required" mapstructure:"url"`
		UrlPolyline string `validate:"required" mapstructure:"url_polyline"`
		UrlRefund   string `validate:"required" mapstructure:"url_refund"`
	} `mapstructure:"server_core"`
	AppConfig struct {
		LimitGetOrder int    `validate:"required" mapstructure:"limit_get_order"`
		MaxPoolWorker int    `validate:"required" mapstructure:"max_pool_worker"`
		MaxRadius     int    `validate:"required" mapstructure:"max_radius"`
		MinRadius     int    `validate:"required" mapstructure:"min_radius"`
		VersionCode   int    `validate:"required" mapstructure:"version_code"`
		ApiToken      string `validate:"required" mapstructure:"api_token"`
	} `mapstructure:"app_config"`
	DB struct {
		Driver         string `validate:"required" mapstructure:"driver"`
		Host           string `validate:"required" mapstructure:"host"`
		Port           int    `validate:"required" mapstructure:"port"`
		Username       string `mapstructure:"username"`
		Password       string `mapstructure:"password"`
		DbName         string `validate:"required" mapstructure:"dbname"`
		DbNameMerchant string `validate:"required" mapstructure:"dbname_merchant"`
		DbNameUser     string `validate:"required" mapstructure:"dbname_user"`
		PoolSize       uint64 `mapstructure:"pool_size"`
	} `mapstructure:"database"`
	Cache struct {
		Driver   string `validate:"required" mapstructure:"driver"`
		Addr     string `validate:"required" mapstructure:"addr"`
		Password string `mapstructure:"password"`
		Database int    `mapstructure:"database"`
		PoolSize int    `mapstructure:"pool_size"`
	} `mapstructure:"cache"`
	MessageBroker struct {
		Driver      string `validate:"required" mapstructure:"driver"`
		Uri         string `validate:"required" mapstructure:"uri"`
		mssgExpired string `validate:"required" mapstructure:"mssg_expired"`
		Exchange    string `mapstructure:"exchange"`
		RouteKey    string `mapstructure:"routekey"`
		Err         chan error
	} `mapstructure:"message_broker"`
	Logger struct {
		Path         string `validate:"required" mapstructure:"path"`
		MaxAge       int    `validate:"required,gte=1,lte=365" mapstructure:"max_age"`
		RotationTime int    `validate:"required,gte=1,lte=365" mapstructure:"rotation_time"`
	} `mapstructure:"logger"`
	FCM struct {
		Token string `validate:"required" mapstructure:"token"`
	} `mapstructure:"fcm"`
}
