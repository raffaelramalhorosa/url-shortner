package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

type URL struct {
	ID        int       `json:"id"`
	ShortCode string    `json:"short_code"`
	LongURL   string    `json:"long_url"`
	Clicks    int       `json:"clicks"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateURLRequest struct {
	LongURL string `json:"long_url" binding:"required"`
}

var db *sql.DB

func main() {
	// Inicializar banco de dados
	initDB()
	defer db.Close()

	// Configurar Gin
	r := gin.Default()

	// Configurar CORS para permitir requisições do frontend
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rotas da API
	api := r.Group("/api")
	{
		api.POST("/shorten", createShortURL)
		api.GET("/urls", listURLs)
		api.GET("/urls/:code", getURLStats)
		api.DELETE("/urls/:code", deleteURL)
	}

	// Rota de redirecionamento
	r.GET("/:code", redirectURL)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Pegar porta do ambiente (Railway) ou usar 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Servidor rodando em http://localhost:%s\n", port)
	r.Run(":" + port)
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "./urls.db")
	if err != nil {
		log.Fatal(err)
	}

	// Criar tabela se não existir
	createTableSQL := `CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		short_code TEXT UNIQUE NOT NULL,
		long_url TEXT NOT NULL,
		clicks INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("✅ Banco de dados inicializado")
}

// Gerar código curto aleatório
func generateShortCode() string {
	b := make([]byte, 6)
	rand.Read(b)
	code := base64.URLEncoding.EncodeToString(b)[:8]
	return code
}

// Criar URL curta
func createShortURL(c *gin.Context) {
	var req CreateURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL inválida"})
		return
	}

	// Gerar código curto único
	shortCode := generateShortCode()

	// Verificar se o código já existe (improvável, mas possível)
	for {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = ?)", shortCode).Scan(&exists)
		if err == nil && !exists {
			break
		}
		shortCode = generateShortCode()
	}

	// Inserir no banco
	_, err := db.Exec(
		"INSERT INTO urls (short_code, long_url) VALUES (?, ?)",
		shortCode, req.LongURL,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar URL"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"short_code": shortCode,
		"short_url":  "http://localhost:8080/" + shortCode,
		"long_url":   req.LongURL,
	})
}

// Redirecionar para URL longa
func redirectURL(c *gin.Context) {
	code := c.Param("code")

	var longURL string
	err := db.QueryRow("SELECT long_url FROM urls WHERE short_code = ?", code).Scan(&longURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL não encontrada"})
		return
	}

	// Incrementar contador de cliques
	db.Exec("UPDATE urls SET clicks = clicks + 1 WHERE short_code = ?", code)

	c.Redirect(http.StatusFound, longURL)
}

// Listar todas as URLs
func listURLs(c *gin.Context) {
	rows, err := db.Query(`
		SELECT id, short_code, long_url, clicks, created_at 
		FROM urls 
		ORDER BY created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar URLs"})
		return
	}
	defer rows.Close()

	var urls []URL
	for rows.Next() {
		var u URL
		err := rows.Scan(&u.ID, &u.ShortCode, &u.LongURL, &u.Clicks, &u.CreatedAt)
		if err != nil {
			continue
		}
		urls = append(urls, u)
	}

	c.JSON(http.StatusOK, urls)
}

// Obter estatísticas de uma URL
func getURLStats(c *gin.Context) {
	code := c.Param("code")

	var u URL
	err := db.QueryRow(`
		SELECT id, short_code, long_url, clicks, created_at 
		FROM urls 
		WHERE short_code = ?
	`, code).Scan(&u.ID, &u.ShortCode, &u.LongURL, &u.Clicks, &u.CreatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL não encontrada"})
		return
	}

	c.JSON(http.StatusOK, u)
}

// Deletar URL
func deleteURL(c *gin.Context) {
	code := c.Param("code")

	result, err := db.Exec("DELETE FROM urls WHERE short_code = ?", code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar URL"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL não encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL deletada com sucesso"})
}
