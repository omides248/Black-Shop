package config

var (
	DefaultAppEnv   = "development"
	DefaultGRPCPort = ":50051"
	DefaultHTTPPort = ":8080"
	DefaultGRPCAddr = "127.0.0.1:50051"
	DefaultHost     = "192.168.8.140:8080"

	DefaultMinIOEndpoint  = "192.168.8.140:9000"
	DefaultMinIOAccessKey = "minioadmin"
	DefaultMinIOSecretKey = "minioadmin123123"
	DefaultMinIOPublicURL = "http://192.168.8.140:9000"

	DefaultMongoURI = "mongodb://omides248:123123@192.168.8.140:27017/?authSource=admin"

	DefaultPublicStoragePath  = "/var/lib/blackshop/storage/public"
	DefaultPrivateStoragePath = "/var/lib/blackshop/storage/private"
	DefaultStaticFilesPrefix  = "public"
)
