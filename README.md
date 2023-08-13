# SP RNNoise
RNNoise based [VoiceDock SP](https://github.com/voicedock/voicedock-specs/blob/main/proto/voicedock/core/sp/v1/) implementation for audio stream noise reduction.

> Provides gRPC API for noise reduction (from raw PCM stream) based on C-binding with 
> [rnnoise library](https://github.com/xiph/rnnoise) ([more about RNNoise](https://hacks.mozilla.org/2017/09/rnnoise-deep-learning-noise-suppression/) neural network design).

# Usage
Run docker container on CPU:
```bash
docker run --rm \
  -p 9999:9999 \
  ghcr.io/voicedock/sprnnoise:latest sprnnoise
```

Show more options:
```bash
docker run --rm ghcr.io/voicedock/sprnnoise sprnnoise -h
```
```
Usage: sprnnoise [--grpcaddr GRPCADDR] [--loglevel LOGLEVEL] [--logjson]

Options:
  --grpcaddr GRPCADDR    gRPC API host:port [default: 0.0.0.0:9999, env: GRPC_ADDR]
  --loglevel LOGLEVEL    log level: debug, info, warn, error, dpanic, panic, fatal [default: info, env: LOG_LEVEL]
  --logjson              set to true to use JSON format [env: LOG_JSON]
  --help, -h             display this help and exit
```

## API
See implementation in [proto file](https://github.com/voicedock/voicedock-specs/blob/main/proto/voicedock/core/sp/v1/sp_api.proto).

## CONTRIBUTING
Lint proto files:
```bash
docker run --rm -w "/work" -v "$(pwd):/work" bufbuild/buf:latest lint internal/api/grpc/proto
```
Generate grpc interface:
```bash
docker run --rm -w "/work" -v "$(pwd):/work" ghcr.io/voicedock/protobuilder:1.0.0 generate internal/api/grpc/proto --template internal/api/grpc/proto/buf.gen.yaml
```
Manual build docker image:
```bash
docker build -t ghcr.io/voicedock/sprnnoise:latest .
```

## Thanks
 * [The Mozilla Research](https://jmvalin.ca/demo/rnnoise/)