# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copiar go.mod e go.sum
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Build do aplicativo
RUN go build -o main .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copiar binário do stage anterior
COPY --from=builder /app/main .

# Expor porta
EXPOSE 8080

# Comando para rodar
CMD ["./main"]