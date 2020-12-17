# build an executable binary
FROM golang:alpine AS builder
RUN apk --no-cache --update upgrade && apk --no-cache add ca-certificates
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main cmd/backend/main.go

# build an image from scratch
FROM scratch
COPY --from=builder /build/main /app/
COPY --from=builder /build/config.yml /app/config.yml
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
WORKDIR /app
CMD ["./main"]
