package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	
	_ "github.com/lib/pq"

	"maqola-backent/internal/delivery/http/handler"
	"maqola-backent/internal/infrastructure/repository/postgres"
	"maqola-backent/internal/usecase"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "root"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "maqolameta"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
		dbHost, dbPort, dbUser, dbPassword, dbName)

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Ma'lumotlar bazasiga ulanishda xatolik:", err)
	}
	defer dbConn.Close()

	for i := 0; i < 5; i++ {
		if err = dbConn.Ping(); err == nil {
			break
		}
		log.Println("Ma'lumotlar bazasi ulanishi kutilmoqda...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Ma'lumotlar bazasiga ulanib bo'lmadi:", err)
	}

	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	timeoutContext := time.Duration(2) * time.Second
	
	articleRepo := postgres.NewPostgresArticleRepository(dbConn)
	articleUseCase := usecase.NewArticleUseCase(articleRepo, timeoutContext)
	
	handler.NewArticleHandler(router, articleUseCase)

	log.Println("Server 8080 portda ishga tushirildi...")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server ishdan chiqdi:", err)
	}
}
