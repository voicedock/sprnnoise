package main

import (
	"github.com/alexflint/go-arg"
	grpcapi "github.com/voicedock/sprnnoise/internal/api/grpc"
	spv1 "github.com/voicedock/sprnnoise/internal/api/grpc/gen/voicedock/core/sp/v1"
	"github.com/voicedock/sprnnoise/internal/rnnoise"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

var cfg AppConfig
var logger *zap.Logger

func init() {
	arg.MustParse(&cfg)
	logger = initLogger(cfg.LogLevel, cfg.LogJson)
}

func main() {
	defer logger.Sync()

	logger.Info("Starting SP RNNoise")

	lis, err := net.Listen("tcp", cfg.GrpcAddr)
	if err != nil {
		logger.Fatal("Failed to listen gRPC server", zap.Error(err))
	}

	rnn := rnnoise.NewService(logger)

	srv := grpcapi.NewServerSp(rnn, logger)

	s := grpc.NewServer()
	spv1.RegisterSpAPIServer(s, srv)
	reflection.Register(s)

	logger.Info("gRPC server listen", zap.String("addr", cfg.GrpcAddr))
	err = s.Serve(lis)
	if err != nil {
		logger.Fatal("gRPC server error", zap.Error(err))
	}
}
