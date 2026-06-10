package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	
	"maqola-backent/internal/domain"
	
	"github.com/lib/pq"
)

type postgresArticleRepository struct {
	Conn *sql.DB
}

func NewPostgresArticleRepository(conn *sql.DB) domain.ArticleRepository {
	return &postgresArticleRepository{Conn: conn}
}

func (m *postgresArticleRepository) Fetch(ctx context.Context, filter domain.ArticleFilter) ([]domain.Article, error) {
	query := `SELECT id, title, access_type, abstract, authors, journal, publisher, publisher_date, doi, url, pdf_url, source_url, key_words FROM articles WHERE 1=1`
	var args []interface{}
	argId := 1

	if filter.Title != "" {
		query += ` AND title ILIKE $` + fmt.Sprint(argId)
		args = append(args, "%"+filter.Title+"%")
		argId++
	}
	if filter.Journal != "" {
		query += ` AND journal ILIKE $` + fmt.Sprint(argId)
		args = append(args, "%"+filter.Journal+"%")
		argId++
	}
	if filter.AccessType != "" {
		query += ` AND access_type = $` + fmt.Sprint(argId)
		args = append(args, filter.AccessType)
		argId++
	}
	if filter.Publisher != "" {
		query += ` AND publisher ILIKE $` + fmt.Sprint(argId)
		args = append(args, "%"+filter.Publisher+"%")
		argId++
	}
	if filter.AuthorName != "" {
		query += ` AND EXISTS (SELECT 1 FROM jsonb_array_elements(authors) AS elem WHERE elem->>'name' ILIKE $` + fmt.Sprint(argId) + `)`
		args = append(args, "%"+filter.AuthorName+"%")
		argId++
	}
	if filter.StartDate != "" {
		query += ` AND publisher_date >= $` + fmt.Sprint(argId)
		args = append(args, filter.StartDate)
		argId++
	}
	if filter.EndDate != "" {
		query += ` AND publisher_date <= $` + fmt.Sprint(argId)
		args = append(args, filter.EndDate)
		argId++
	}
	
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []domain.Article
	for rows.Next() {
		var a domain.Article
		var authorsJSON []byte
		var keyWords pq.StringArray 

		err := rows.Scan(
			&a.ID, &a.Title, &a.AccessType, &a.Abstract,
			&authorsJSON, &a.Journal, &a.Publisher, &a.PublisherDate,
			&a.DOI, &a.URL, &a.PDFUrl, &a.SourceURL, &keyWords,
		)
		if err != nil {
			return nil, err
		}
		
		a.KeyWords = keyWords
		if len(authorsJSON) > 0 {
			json.Unmarshal(authorsJSON, &a.Authors) 
		}

		articles = append(articles, a)
	}

	if articles == nil {
		articles = []domain.Article{} 
	}

	return articles, nil
}