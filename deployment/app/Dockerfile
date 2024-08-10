FROM golang:1.21 as build

WORKDIR /app

COPY . .
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -v -o ./bin/app ./cmd/app

FROM scratch

WORKDIR /app

COPY --from=build /app/bin/app ./bin/
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY ./config/ ./config/

CMD ["./bin/app"]