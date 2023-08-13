package grpc

import (
	commonv1 "github.com/voicedock/sprnnoise/internal/api/grpc/gen/voicedock/core/common/v1"
	spv1 "github.com/voicedock/sprnnoise/internal/api/grpc/gen/voicedock/core/sp/v1"
	"github.com/voicedock/sprnnoise/internal/rnnoise"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

func NewServerSp(rnnoise *rnnoise.Service, logger *zap.Logger) *ServerSp {
	return &ServerSp{
		rnnoise: rnnoise,
		logger:  logger,
	}
}

type ServerSp struct {
	rnnoise *rnnoise.Service
	logger  *zap.Logger
	spv1.UnimplementedSpAPIServer
}

func (s *ServerSp) ProcessSound(srv spv1.SpAPI_ProcessSoundServer) error {
	chIn := make(chan []byte)
	chOut := make(chan []byte)

	go func() {
		err := s.readData(srv, chIn)
		if err != nil {
			s.logger.Error("Failed read data", zap.Error(err))
		}
	}()

	go func() {
		err := s.writeData(srv, chOut)
		if err != nil {
			s.logger.Error("Failed write data", zap.Error(err))
		}
	}()

	err := s.rnnoise.ProcessSound(chIn, chOut)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (s *ServerSp) readData(srv spv1.SpAPI_ProcessSoundServer, chIn chan []byte) error {
	defer close(chIn)

	for {
		req, err := srv.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "failed to read data: %s", err)
		}

		if req.GetAudio().SampleRate != 48000 {
			return status.Error(codes.InvalidArgument, "unsupported sample rate: required 48000")
		}
		if req.GetAudio().Channels != 0 {
			return status.Error(codes.InvalidArgument, "unsupported audio channels: required 1")
		}

		chIn <- req.Audio.Data
	}

	return nil
}

func (s *ServerSp) writeData(srv spv1.SpAPI_ProcessSoundServer, chOut chan []byte) error {
	for v := range chOut {
		err := srv.Send(&spv1.ProcessSoundResponse{Audio: &commonv1.AudioContainer{
			Data:       v,
			SampleRate: 48000,
			Channels:   1,
		}})

		if err != nil {
			return status.Errorf(codes.Internal, "failed send data: %", err)
		}
	}

	return nil
}
