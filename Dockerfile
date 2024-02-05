FROM public.ecr.aws/docker/library/golang:1.20-alpine AS builder
WORKDIR /builder
RUN apk --no-cache add ca-certificates zoneinfo

ADD source /builder
RUN CGO_ENABLED=0 go build -o pollendata
RUN apk --no-cache add upx
RUN upx --best pollendata

FROM scratch
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app
COPY --from=builder /builder/pollendata /app/

ENTRYPOINT ["/app/pollendata"]

