FROM golang:1.20 as build

WORKDIR /app

COPY . .
RUN go mod download

ARG BUILD_VERSION

RUN CGO_ENABLED=0 GOOS=linux go build -v -o "./bin/sender" \
    -ldflags "-X main.release=$BUILD_VERSION -X main.buildDate=$(date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(git log --format="%h" -n 1)" \
    ./cmd/calendar_sender

FROM scratch

WORKDIR /app

COPY --from=build /app/bin/sender ./bin/

COPY ./configs/sender_config.toml ./configs/
COPY ./configs/sender.env ./configs/

CMD ["./bin/sender", "-config", "./configs/sender_config.toml"]