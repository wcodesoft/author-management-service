FROM golang:1.18.3-alpine3.16 AS builder

# Change the workdir to inside the service folder
WORKDIR /app/service

COPY . .

RUN go get .

# Build the image
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /bin

COPY --from=builder /app/service .

ENTRYPOINT [ "./app" ]

EXPOSE 9000
EXPOSE 9001