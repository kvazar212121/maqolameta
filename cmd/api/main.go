package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	
	_ "github.com/lib/pq" // PostgreSQL driver (importni unutmang: go get github.com/lib/pq)

	"maqola-backent/internal/delivery/http/handler"
	"maqola-backent/internal/infrastructure/repository/postgres"
	"maqola-backent/internal/usecase"
)

func main() {
	// 1. PostgreSQL bazasiga ulanish (Bu yerda ulanish simini o'zingizga moslang)
	dbConn, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=root dbname=maqolameta sslmode=disable")
	if err != nil {
		log.Fatal("Ma'lumotlar bazasiga ulanishda xatolik:", err)
	}
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		log.Fatal("Ma'lumotlar bazasiga ulanib bo'lmadi:", err)
	}

	// 2. Gin serverini yaratish
	router := gin.Default()

	// 3. Clean Architecture qatlamlarini bog'lash (Dependency Injection)
	timeoutContext := time.Duration(2) * time.Second
	
	articleRepo := postgres.NewPostgresArticleRepository(dbConn)
	articleUseCase := usecase.NewArticleUseCase(articleRepo, timeoutContext)
	
	// Handlerlarni ulab qo'yish
	handler.NewArticleHandler(router, articleUseCase)

	// 4. Serverni ishga tushirish
	log.Println("Server 8080 portda ishga tushirildi...")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server ishdan chiqdi:", err)
	}
}
