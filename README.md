# URL Shortener

Sistema de encurtamento de URLs desenvolvido com React e Go.

## Tecnologias Utilizadas

### Backend
- Go 1.21+
- Gin (framework web)
- SQLite (banco de dados)
- CORS middleware

### Frontend
- React 18
- TypeScript
- Tailwind CSS
- Vite

## Funcionalidades

- Criação de URLs curtas a partir de URLs longas
- Redirecionamento automático
- Contador de cliques
- Listagem de todas as URLs criadas
- Exclusão de URLs
- Interface responsiva

## Pré-requisitos

- Go 1.21 ou superior
- Node.js 18 ou superior
- npm ou yarn

## Instalação

### Backend

```bash
cd backend
go mod download
go run main.go
```

O servidor estará disponível em `http://localhost:8080`

### Frontend

```bash
cd frontend
npm install
npm run dev
```

A aplicação estará disponível em `http://localhost:5173`

## Executar Ambos Simultaneamente

### Opção 1: Script Batch (Windows)
```bash
start.bat
```

### Opção 2: Script PowerShell (Windows)
```bash
.\start.ps1
```

### Opção 3: Manualmente
Abra dois terminais e execute os comandos de backend e frontend separadamente.

## Estrutura do Projeto

```
url-shortener/
├── backend/
│   ├── main.go
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── App.tsx
│   │   ├── main.tsx
│   │   └── index.css
│   ├── package.json
│   └── vite.config.ts
├── .gitignore
└── README.md
```

## API Endpoints

### POST /api/shorten
Cria uma URL curta
```json
{
  "long_url": "https://exemplo.com/url-muito-longa"
}
```

### GET /api/urls
Lista todas as URLs criadas

### GET /api/urls/:code
Retorna estatísticas de uma URL específica

### DELETE /api/urls/:code
Remove uma URL

### GET /:code
Redireciona para a URL original

## Banco de Dados

O sistema utiliza SQLite com a seguinte estrutura:

```sql
CREATE TABLE urls (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    short_code TEXT UNIQUE NOT NULL,
    long_url TEXT NOT NULL,
    clicks INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## Desenvolvimento

### Backend
O backend utiliza o framework Gin para roteamento e middlewares. O banco de dados SQLite é criado automaticamente ao iniciar a aplicação.

### Frontend
O frontend foi desenvolvido com React e TypeScript, utilizando Tailwind CSS para estilização. As requisições HTTP são feitas com a Fetch API nativa.

## Licença

MIT