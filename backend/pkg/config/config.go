package config

import "github.com/spf13/viper"

type Config struct {
	// Catalog
	CatalogGRPCPort string `mapstructure:"CATALOG_GRPC_PORT"`
	CatalogHTTPPort string `mapstructure:"CATALOG_HTTP_PORT"`
	CatalogGRPCAddr string `mapstructure:"CATALOG_GRPC_ADDR"`
	MongoURI        string `mapstructure:"MONGO_URI"`
	// Identity
	IdentityGRPCPort string `mapstructure:"IDENTITY_GRPC_PORT"`
	IdentityHTTPPort string `mapstructure:"IDENTITY_HTTP_PORT"`
	IdentityGRPCAddr string `mapstructure:"IDENTITY_GRPC_ADDR"`
	PostgresURI      string `mapstructure:"POSTGRES_URI"`
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
