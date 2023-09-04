FROM golang:1.20

WORKDIR /app

COPY . .

CMD ["go", "test", "-v", "--tags=integration", "./tests/integration/..."]
