FROM public.ecr.aws/docker/library/golang:alpine AS builder
WORKDIR /builder
ADD source /builder
RUN CGO_ENABLED=0 go build -o pollendata

FROM public.ecr.aws/docker/library/alpine:latest
RUN apk --no-cache add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=builder /builder/pollendata /app/

ENTRYPOINT ["/app/pollendata"]

