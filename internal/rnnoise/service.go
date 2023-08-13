package rnnoise

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/voicedock/audio"
	"go.uber.org/zap"
	"io"
)

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -l:librnnoise.a -lm -Wl,-rpath=/usr/local/lib
#include "rnnoise.h"
*/
import "C"

type Service struct {
	logger    *zap.Logger
	frameSize int
}

func NewService(
	logger *zap.Logger,
) *Service {
	return &Service{
		logger:    logger,
		frameSize: 480,
	}
}

func (s *Service) ProcessSound(chIn chan []byte, chOut chan []byte) error {
	st := C.rnnoise_create(nil)
	defer func() {
		C.rnnoise_destroy(st)
		close(chOut)
	}()

	r := new(bytes.Buffer)
	for v := range chIn {
		r.Write(v)
		frameSize := r.Len() / 2 / s.frameSize
		if frameSize < 1 {
			continue
		}

		var buf bytes.Buffer
		s.processFrame(st, frameSize, r, &buf)
		chOut <- buf.Bytes()
	}

	if r.Len() == 0 {
		return nil
	}

	countInt := r.Len() / 2
	normalizeLen := countInt/s.frameSize*s.frameSize + s.frameSize
	if normalizeLen > s.frameSize {
		zero := int16(0)
		for i := 0; i < normalizeLen-countInt; i++ {
			binary.Write(r, binary.LittleEndian, &zero)
		}
	}

	var buf bytes.Buffer
	s.processFrame(st, normalizeLen, r, &buf)
	chOut <- buf.Bytes()

	return nil
}

func (s *Service) processFrame(st *C.DenoiseState, frameSize int, r io.Reader, w io.Writer) error {
	// read data
	int16Pcm := make([]int16, frameSize)
	err := binary.Read(r, binary.LittleEndian, &int16Pcm)
	if err != nil {
		return fmt.Errorf("unable to read pcm int16 LE format data from bytes: %w", err)
	}

	data := audio.ConvertNumbers[float32](int16Pcm)

	// processing data
	C.rnnoise_process_frame(st, (*C.float)(&data[0]), (*C.float)(&data[0]))

	// send data
	err = binary.Write(w, binary.LittleEndian, audio.ConvertNumbers[int16](data))
	if err != nil {
		return fmt.Errorf("unable to write pcm int16 LE format data: %w", err)
	}

	return nil
}
