FROM golang:1.20

WORKDIR /app

COPY . .

CMD ["sh", "-c", "if go test -v --tags=integration /app/tests/integration/...; then echo 1; else echo 0; fi"]