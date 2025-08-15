package config

import "github.com/spf13/viper"

type Config struct {
	// Catalog Service
	CatalogGRPCPort string `mapstructure:"CATALOG_GRPC_PORT"`
	CatalogHTTPPort string `mapstructure:"CATALOG_HTTP_PORT"`
	CatalogGRPCAddr string `mapstructure:"CATALOG_GRPC_ADDR"`
	CatalogHost     string `mapstructure:"CATALOG_HOST"`

	// Identity Service
	IdentityGRPCPort string `mapstructure:"IDENTITY_GRPC_PORT"`
	IdentityHTTPPort string `mapstructure:"IDENTITY_HTTP_PORT"`
	IdentityGRPCAddr string `mapstructure:"IDENTITY_GRPC_ADDR"`

	// Order Service
	OrderGRPCPort string `mapstructure:"ORDER_GRPC_PORT"`
	OrderHTTPPort string `mapstructure:"ORDER_HTTP_PORT"`
	OrderGRPCAddr string `mapstructure:"ORDER_GRPC_ADDR"`

	// Database
	MongoURI            string `mapstructure:"MONGO_URI"`
	PostgresIdentityURI string `mapstructure:"POSTGRES_IDENTITY_URI"`
	PostgresOrderURI    string `mapstructure:"POSTGRES_ORDER_URI"`
	RedisAddr           string `mapstructure:"REDIS_ADDR"`

	// MinIO
	MinioEndpoint  string `mapstructure:"MINIO_ENDPOINT"`
	MinioAccessKey string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioSecretKey string `mapstructure:"MINIO_SECRET_KEY"`
	MinioPublicURL string `mapstructure:"MINIO_PUBLIC_URL"`

	// Local Storage
	LocalStoragePath string `mapstructure:"LOCAL_STORAGE_PATH"`

	// General
	AppEnv string `mapstructure:"APP_ENV"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
