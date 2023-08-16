# syntax=docker/dockerfile:1

FROM golang:1.21.0-bullseye

WORKDIR /app

COPY . ./

RUN make build

CMD ["make", "start"]
