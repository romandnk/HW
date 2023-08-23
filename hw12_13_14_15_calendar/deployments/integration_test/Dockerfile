FROM golang:1.20

WORKDIR /app

COPY . .

CMD ["sh", "-c", "if go test -v --tags=integration /app/tests/integration/...; then echo 0; else echo 1; fi"]