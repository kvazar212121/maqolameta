package postgres

import (
	"context"
	"database/sql"
	"maqola-backent/internal/domain"
)

type postgresArticleRepository struct {
	Conn *sql.DB
}

// NewPostgresArticleRepository - Repository initsializatsiyasi
func NewPostgresArticleRepository(conn *sql.DB) domain.ArticleRepository {
	return &postgresArticleRepository{Conn: conn}
}

// Fetch - PostgreSQL'dan maqolalarni olish
func (m *postgresArticleRepository) Fetch(ctx context.Context) ([]domain.Article, error) {
	// Bu yerda haqiqiy SQL so'rovi yoziladi (hozircha mock/dummy data)
	// PostgreSQL dan ma'lumotlarni o'qish logikasi bo'lishi kerak.
	// Hozirgi loyiha poydevori uchun bo'sh massiv qaytaramiz:
	
	articles := []domain.Article{}
	
	// TODO: Quyidagicha SQL ishlashi kerak:
	// query := `SELECT id, title, access_type, abstract, journal, publisher, publisher_date, doi, url, pdf_url, source_url FROM articles`
	// keyin authors va keywords ni alohida o'qib kelish mumkin (yoki JOIN orqali).
	
	return articles, nil
}
