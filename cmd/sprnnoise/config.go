package main

type AppConfig struct {
	GrpcAddr string `arg:"env:GRPC_ADDR" help:"gRPC API host:port" default:"0.0.0.0:9999"`
	LogLevel string `arg:"env:LOG_LEVEL" help:"log level: debug, info, warn, error, dpanic, panic, fatal" default:"info"`
	LogJson  bool   `arg:"env:LOG_JSON" help:"set to true to use JSON format"`
}
