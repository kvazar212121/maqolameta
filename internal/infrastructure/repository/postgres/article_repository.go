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

func (m *postgresArticleRepository) Fetch(ctx context.Context, filter domain.ArticleFilter) ([]domain.Article, int, error) {
	baseQuery := ` FROM articles a LEFT JOIN article_views v ON a.id = v.article_id WHERE 1=1`
	var args []interface{}
	argId := 1

	if filter.Title != "" {
		baseQuery += ` AND a.title ILIKE $` + fmt.Sprint(argId)
		args = append(args, "%"+filter.Title+"%")
		argId++
	}
	if filter.Journal != "" {
		baseQuery += ` AND a.journal ILIKE $` + fmt.Sprint(argId)
		args = append(args, "%"+filter.Journal+"%")
		argId++
	}
	if filter.AccessType != "" {
		baseQuery += ` AND a.access_type = $` + fmt.Sprint(argId)
		args = append(args, filter.AccessType)
		argId++
	}
	if filter.Publisher != "" {
		baseQuery += ` AND a.publisher ILIKE $` + fmt.Sprint(argId)
		args = append(args, "%"+filter.Publisher+"%")
		argId++
	}
	if filter.AuthorName != "" {
		baseQuery += ` AND EXISTS (SELECT 1 FROM jsonb_array_elements(a.authors) AS elem WHERE elem->>'name' ILIKE $` + fmt.Sprint(argId) + `)`
		args = append(args, "%"+filter.AuthorName+"%")
		argId++
	}
	if filter.StartDate != "" {
		baseQuery += ` AND a.publisher_date >= $` + fmt.Sprint(argId)
		args = append(args, filter.StartDate)
		argId++
	}
	if filter.EndDate != "" {
		baseQuery += ` AND a.publisher_date <= $` + fmt.Sprint(argId)
		args = append(args, filter.EndDate)
		argId++
	}
	if filter.KeyWord != "" {
		baseQuery += ` AND $` + fmt.Sprint(argId) + ` ILIKE ANY(a.key_words)`
		args = append(args, "%"+filter.KeyWord+"%")
		argId++
	}

	var total int
	if argId == 1 {
		// Hech qanday filter yo'q, barcha maqolalar sonini pg_class orqali tezkor olamiz
		m.Conn.QueryRowContext(ctx, `SELECT reltuples::bigint FROM pg_class WHERE relname = 'articles'`).Scan(&total)
	} else {
		// Filter bor, shuning uchun faqat mos tushganlarni sanaymiz, lekin 1000 tadan ortig'iga vaqt ketkazmaymiz
		countQuery := `SELECT COUNT(*) FROM (SELECT 1 ` + baseQuery + ` LIMIT 1000) as sub`
		m.Conn.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	}

	var query string
	if argId == 1 {
		// Bosh sahifada (hech qanday filter yo'q) barcha 3 million qatorni ORDER BY qilish serverni qotiradi.
		// Shuning uchun eng so'nggi maqolalarni to'g'ridan-to'g'ri tezkor index (id) orqali olamiz.
		query = `SELECT a.id, a.title, a.access_type, a.abstract, a.authors, a.journal, a.publisher, a.publisher_date, a.doi, a.url, a.pdf_url, a.source_url, a.key_words, COALESCE(v.views_count, 0) ` + baseQuery
		query += ` ORDER BY a.id DESC`
	} else {
		// Qidiruv va filterlar mavjud bo'lsa, natijalar kam bo'lgani uchun view bo'yicha saralash xavfsiz.
		query = `SELECT a.id, a.title, a.access_type, a.abstract, a.authors, a.journal, a.publisher, a.publisher_date, a.doi, a.url, a.pdf_url, a.source_url, a.key_words, COALESCE(v.views_count, 0) ` + baseQuery
		query += ` ORDER BY COALESCE(v.views_count, 0) DESC, a.id ASC`
	}

	if filter.Limit > 0 {
		query += ` LIMIT $` + fmt.Sprint(argId)
		args = append(args, filter.Limit)
		argId++
	}
	if filter.Offset > 0 {
		query += ` OFFSET $` + fmt.Sprint(argId)
		args = append(args, filter.Offset)
		argId++
	}
	
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var articles []domain.Article
	for rows.Next() {
		var a domain.Article
		var authorsJSON []byte
		var keyWords pq.StringArray 
		var abstract, journal, publisher, pubDate, doi, url, pdfUrl, sourceUrl sql.NullString

		err := rows.Scan(
			&a.ID, &a.Title, &a.AccessType, &abstract,
			&authorsJSON, &journal, &publisher, &pubDate,
			&doi, &url, &pdfUrl, &sourceUrl, &keyWords, &a.ViewsCount,
		)
		if err != nil {
			return nil, 0, err
		}
		
		a.Abstract = abstract.String
		a.Journal = journal.String
		a.Publisher = publisher.String
		a.PublisherDate = pubDate.String
		a.DOI = doi.String
		a.URL = url.String
		a.PDFUrl = pdfUrl.String
		a.SourceURL = sourceUrl.String
		a.KeyWords = keyWords
		if len(authorsJSON) > 0 {
			json.Unmarshal(authorsJSON, &a.Authors) 
		}

		articles = append(articles, a)
	}

	if articles == nil {
		articles = []domain.Article{} 
	}

	return articles, total, nil
}

func (m *postgresArticleRepository) GetUniqueKeyWords(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT unnest(key_words) FROM articles WHERE key_words IS NOT NULL`
	
	rows, err := m.Conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keywords []string
	for rows.Next() {
		var kw string
		if err := rows.Scan(&kw); err != nil {
			return nil, err
		}
		if kw != "" {
			keywords = append(keywords, kw)
		}
	}

	if keywords == nil {
		keywords = []string{}
	}

	return keywords, nil
}

func (m *postgresArticleRepository) AddViews(ctx context.Context, articleID string, viewsToAdd int) error {
	query := `
		INSERT INTO article_views (article_id, views_count)
		VALUES ($1, $2)
		ON CONFLICT (article_id) 
		DO UPDATE SET views_count = article_views.views_count + EXCLUDED.views_count
	`
	
	_, err := m.Conn.ExecContext(ctx, query, articleID, viewsToAdd)
	return err
}