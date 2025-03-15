FROM golang:1.22.5

RUN apt update && apt upgrade -y && \
    apt install -y git \
    make openssh-client

WORKDIR /app

RUN go install github.com/cosmtrek/air@v1.48.0

CMD air