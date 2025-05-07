FROM golang:1.24.3-alpine AS builder

WORKDIR /builder
RUN apk update && apk --no-cache add ca-certificates tzdata

ADD source /builder
RUN CGO_ENABLED=0 go build -o pollendata

FROM ghcr.io/sollie/docker-upx:v5.0.0 AS upx
WORKDIR /upx
COPY --from=builder /builder/pollendata /upx/pollendata
RUN upx --best pollendata

FROM scratch AS final
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app
COPY --from=upx /upx/pollendata /app/

EXPOSE 8080

ENTRYPOINT ["/app/pollendata"]
