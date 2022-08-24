package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port                 int    `mapstructure:"SERVER_PORT" required:"true"`
	DBConnStr            string `mapstructure:"DB_CONN_STR" required:"true"`
	DriverName           string `mapstructure:"DB_DRIVER_NAME" required:"true"`
	ClusterName          string `mapstructure:"NATS_CLUSTER_NAME" required:"true"`
	NatsURL              string `mapstructure:"NATS_URL" required:"true"`
	NatsClient           string `mapstructure:"NATS_CLIENT" required:"true"`
	NatsSubject          string `mapstructure:"NATS_SUBJECT" required:"true"`
	CacheExpTime         string `mapstructure:"CACHE_KEY_EXPIRATION_TIME" required:"true"`
	CacheCleanupInterval string `mapstructure:"CACHE_CLEANUP_INTERVAL" required:"true"`
	DurableName          string `mapstructure:"DURABLE_NAME" required:"true"`
	CacheSize            int    `mapstructure:"CACHE_SIZE" required:"true"`
	ShutdownTimeout      int    `mapstructure:"SHUTDOWN_TIMEOUT" required:"true"`
	GatewayURL           string `mapstructure:"GATEWAY_URL" required:"true"`
	DebugAuth            bool   `mapstructure:"DEBUG_AUTH" default:"false"` // DebugAuth флаг для отключения проверки авторизации
	TokenSignedKey       string `mapstructure:"TOKEN_SIGNED_KEY" default:"AAA"`

	// MinIO
	MinIOURL       string `mapstructure:"MINIO_URL" required:"true"`
	MinIOUser      string `mapstructure:"MINIO_USER" required:"true"`
	MinIOPass      string `mapstructure:"MINIO_PASS" required:"true"`
	MinIODocBucket string `mapstructure:"MINIO_DOC_BUCKET" required:"true"`
}

//func getTime(input string) (time.Duration, error) {
//	return time.ParseDuration(input)
//}

func NewConfig() (*Config, error) {
	viper.SetConfigName("config.env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var conf Config
	err = viper.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, err
}
