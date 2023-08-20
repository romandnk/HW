FROM golang:1.20 as build

WORKDIR /app

COPY . .
RUN go mod download

ARG BUILD_VERSION

RUN CGO_ENABLED=0 GOOS=linux go build -v -o "./bin/calendar" \
    -ldflags "-X main.release=$BUILD_VERSION -X main.buildDate=$(date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(git log --format="%h" -n 1)" \
    ./cmd/calendar

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/bin/calendar ./bin/

COPY ./configs/calendar_config.toml ./configs/
COPY ./configs/calendar.env ./configs/
RUN mkdir -p ./logging && touch ./logging/logging.txt

CMD ["./bin/calendar", "-config", "./configs/calendar_config.toml"]