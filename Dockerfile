FROM golang:1.23.2-alpine AS builder

WORKDIR /builder
RUN apk update && apk --no-cache add ca-certificates tzdata

ADD source /builder
RUN CGO_ENABLED=0 go build -o pollendata

FROM ghcr.io/sollie/docker-upx:v4.2.4 AS upx
WORKDIR /upx
COPY --from=builder /builder/pollendata /upx/pollendata
RUN upx --best pollendata

FROM scratch AS final
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app
COPY --from=upx /upx/pollendata /app/

ENTRYPOINT ["/app/pollendata"]

