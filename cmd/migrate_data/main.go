package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Connect to Postgres
	pgDSN := "host=localhost port=5432 user=postgres password=root dbname=maqolameta sslmode=disable"
	pgDB, err := sql.Open("postgres", pgDSN)
	if err != nil {
		log.Fatal("Error connecting to Postgres:", err)
	}
	defer pgDB.Close()

	if err := pgDB.Ping(); err != nil {
		log.Fatal("Postgres ping failed:", err)
	}
	log.Println("Connected to Postgres")

	// Apply Schema
	schema, err := os.ReadFile("migrations/000001_init_schema.up.sql")
	if err != nil {
		log.Fatal("Failed to read schema:", err)
	}
	_, err = pgDB.Exec(string(schema))
	if err != nil {
		log.Fatal("Failed to apply schema:", err)
	}
	log.Println("Schema applied")

	// Connect to SQLite
	slDB, err := sql.Open("sqlite3", "../maqola_arxiv.db")
	if err != nil {
		log.Fatal("Error connecting to SQLite:", err)
	}
	defer slDB.Close()

	if err := slDB.Ping(); err != nil {
		log.Fatal("SQLite ping failed:", err)
	}
	log.Println("Connected to SQLite")

	// Fetch data
	query := `
	SELECT 
		m.id, 
		IFNULL(m.mavzu, 'Nomsiz maqola'), 
		IFNULL(m.foydalanish_huquqi, 'open'), 
		IFNULL(m.annotatsiya, ''), 
		IFNULL(m.jurnal_nomi, ''), 
		IFNULL(m.nashriyot, ''), 
		IFNULL(m.chop_qilingan_sana, ''), 
		IFNULL(m.doi, ''), 
		IFNULL(m.url, ''), 
		IFNULL(m.kalit_sozlar, ''),
		(SELECT json_group_array(json_object('name', fio, 'affiliation', '', 'orcid', IFNULL(science_id, ''))) 
		 FROM maqola_mualliflar a WHERE a.maqola_id = m.id ORDER BY tartib) as authors_json,
		(SELECT zenodo_url FROM maqola_fayllar f WHERE f.maqola_id = m.id AND fayl_nomi LIKE '%.pdf' LIMIT 1) as pdf_url
	FROM ilmiy_maqolalar m
	`

	rows, err := slDB.Query(query)
	if err != nil {
		log.Fatal("Failed to query SQLite:", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var title, accessType, abstract, journal, publisher, pubDate, doi, url, keywordsRaw, authorsJSON sql.NullString
		var pdfUrl sql.NullString

		err := rows.Scan(&id, &title, &accessType, &abstract, &journal, &publisher, &pubDate, &doi, &url, &keywordsRaw, &authorsJSON, &pdfUrl)
		if err != nil {
			log.Fatal("Failed to scan row:", err)
		}

		// Process keywords
		var keywords []string
		if keywordsRaw.Valid && keywordsRaw.String != "" {
			// They might be comma separated or JSON. Assume comma separated for now.
			parts := strings.Split(keywordsRaw.String, ",")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p != "" {
					keywords = append(keywords, p)
				}
			}
		}

		// Ensure authorsJSON is valid
		authors := "[]"
		if authorsJSON.Valid && authorsJSON.String != "" {
			authors = authorsJSON.String
		}

		insertQuery := `
			INSERT INTO articles (title, access_type, abstract, authors, journal, publisher, publisher_date, doi, url, pdf_url, key_words)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`
		
		doiStr := sql.NullString{String: doi.String, Valid: doi.String != ""}
		pdfStr := sql.NullString{String: pdfUrl.String, Valid: pdfUrl.String != ""}

		_, err = pgDB.Exec(insertQuery, 
			title.String, 
			accessType.String, 
			abstract.String, 
			authors, 
			journal.String, 
			publisher.String, 
			pubDate.String, 
			doiStr, 
			url.String, 
			pdfStr, 
			pq.Array(keywords),
		)
		if err != nil {
			// Skip duplicates or unique violations
			if !strings.Contains(err.Error(), "unique constraint") {
				log.Printf("Failed to insert article %d: %v", id, err)
			}
		} else {
			count++
		}
	}

	log.Printf("Successfully migrated %d articles to Postgres", count)
}
