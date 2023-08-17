# syntax=docker/dockerfile:1

FROM golang:1.21.0 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/app ./src/*.go


# Use the smallest base image possible for the final image
FROM scratch

COPY --from=builder /app/build/app /app

COPY .prod.env ./.env

EXPOSE 8081

ENTRYPOINT ["/app"]
