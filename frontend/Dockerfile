# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copiar package.json e package-lock.json
COPY package*.json ./

# Instalar dependências
RUN npm ci

# Copiar código fonte
COPY . .

# Build do projeto
RUN npm run build

# Production stage
FROM node:20-alpine

WORKDIR /app

# Instalar serve para servir os arquivos estáticos
RUN npm install -g serve

# Copiar build do stage anterior
COPY --from=builder /app/dist ./dist

# Expor porta
EXPOSE 3000

# Comando para servir
CMD ["serve", "-s", "dist", "-l", "3000"]