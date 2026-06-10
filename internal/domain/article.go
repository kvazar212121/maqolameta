package domain

import "context"

type Author struct {
	Name        string `json:"name"`        // ismi sharifi
	Affiliation string `json:"affiliation"` // ish joyi
	ORCID       string `json:"orcid"`       // xalqaro mualliflar id si
}

type Article struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`         // 1 mavzu uzb rus eng
	AccessType    string   `json:"accessType"`    // 2 foydalanuvchi huquqi (ochiq/yopiq)
	Authors       []Author `json:"authors"`       // 3 mualliflar ro'yxati
	Abstract      string   `json:"abstract"`      // 4 annotatsiyasi
	KeyWords      []string `json:"keyWords"`      // 5 kalit so'zlar
	Journal       string   `json:"journal"`       // 6 jurnal nomi
	Publisher     string   `json:"publisher"`     // 7 jurnal chiqaruvchi tashkilot
	PublisherDate string   `json:"publisherDate"` // 8 chop etilgan sana
	DOI           string   `json:"doi"`           // 9 maqola DOI
	URL           string   `json:"url"`           // 10 maqola url manzil
	PDFUrl        string   `json:"pdfUrl"`        // 11 pdf fayl manzili
	SourceURL     string   `json:"sourceUrl"`     // 12 manba URL
}

type ArticleFilter struct {
	Title      string `json:"title"`
	Journal    string `json:"journal"`
	AccessType string `json:"accessType"`
	Publisher  string `json:"publisher"`
	AuthorName string `json:"authorName"`
	StartDate  string `json:"startDate"`
	EndDate    string `json:"endDate"`
	KeyWord    string `json:"keyWord"`
}

type ArticleRepository interface {
	Fetch(ctx context.Context, filter ArticleFilter) ([]Article, error)
	GetUniqueKeyWords(ctx context.Context) ([]string, error)
}

type ArticleUseCase interface {
	Fetch(ctx context.Context, filter ArticleFilter) ([]Article, error)
	GetUniqueKeyWords(ctx context.Context) ([]string, error)
}
