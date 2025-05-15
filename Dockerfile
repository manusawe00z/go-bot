FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go build -o app ./cmd/main.go

FROM python:3.11-slim
WORKDIR /app
COPY --from=builder /app/app .
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY . .

CMD ["./app"]