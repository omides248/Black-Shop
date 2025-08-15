package config

import "github.com/spf13/viper"

type Config struct {
	General      General      `mapstructure:"general"`
	MinIO        MinIO        `mapstructure:"minio"`
	Database     Database     `mapstructure:"database"`
	LocalStorage LocalStorage `mapstructure:"local_storage"`
}

type General struct {
	AppEnv   string `mapstructure:"app_env"`
	GRPCPort string `mapstructure:"grpc_port"`
	HTTPPort string `mapstructure:"http_port"`
	GRPCAddr string `mapstructure:"grpc_addr"`
	Host     string `mapstructure:"host"`
}

type MinIO struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	PublicURL string `mapstructure:"public_url"`
}

type Database struct {
	MongoURI string `mapstructure:"mongo_uri"`
}

type LocalStorage struct {
	PublicStoragePath  string `mapstructure:"public_storage_path"`
	PrivateStoragePath string `mapstructure:"private_storage_path"`
	StaticFilesPrefix  string `mapstructure:"static_files_prefix"`
}

func setDefault(v *viper.Viper) {
	v.SetDefault("general.app_env", DefaultAppEnv)
	v.SetDefault("general.grpc_port", DefaultGRPCPort)
	v.SetDefault("general.http_port", DefaultHTTPPort)
	v.SetDefault("general.grpc_addr", DefaultGRPCAddr)
	v.SetDefault("general.host", DefaultHost)

	v.SetDefault("minio.endpoint", DefaultMinIOEndpoint)
	v.SetDefault("minio.access_key", DefaultMinIOAccessKey)
	v.SetDefault("minio.secret_key", DefaultMinIOSecretKey)
	v.SetDefault("minio.public_url", DefaultMinIOPublicURL)

	v.SetDefault("database.mongo_uri", DefaultMongoURI)

	v.SetDefault("local_storage.public_storage_path", DefaultPublicStoragePath)
	v.SetDefault("local_storage.private_storage_path", DefaultPrivateStoragePath)
	v.SetDefault("local_storage.static_files_prefix", DefaultStaticFilesPrefix)
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	setDefault(viper.GetViper())

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
