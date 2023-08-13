FROM golang:1.20 as builder

RUN apt update && apt install -y \
        build-essential \
        wget  \
        libtool \
        autoconf && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /build

ADD . /usr/src/app

RUN wget https://github.com/xiph/rnnoise/archive/refs/heads/master.tar.gz && \
    tar -xvf "master.tar.gz" --strip-components 1 -C "./" && \
    ./autogen.sh && \
    ./configure && \
    make && \
    make install && \
    cd /usr/src/app && \
    go build -o ./sprnnoise ./cmd/sprnnoise

FROM debian:12

RUN apt update && \
    apt install -y ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/src/app/sprnnoise /usr/local/bin/sprnnoise

CMD ["sprnnoise"]