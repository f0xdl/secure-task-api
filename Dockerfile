FROM golang:1.23 AS builder
ARG CGO_ENABLED=0
LABEL authors="f0xdl"
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o ./bin/app ./cmd/main.go

FROM scratch
COPY --from=builder /src/bin/app /app
#EXPOSE 80
ENTRYPOINT ["/app"]