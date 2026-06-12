package handler

import (
	"io"
	"net/http"
	"net/url"
	"strconv"

	"maqola-backent/internal/domain"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	ArticleUseCase domain.ArticleUseCase
}

func NewArticleHandler(r *gin.Engine, us domain.ArticleUseCase) {
	handler := &ArticleHandler{
		ArticleUseCase: us,
	}

	r.GET("/health", handler.Health)
	r.GET("/api/v1/articles", handler.FetchArticles)
	r.GET("/api/v1/keywords", handler.FetchUniqueKeyWords)
	r.POST("/api/v1/articles/:id/views", handler.AddViews)
	r.GET("/api/v1/proxy/pdf", handler.ProxyPDF)
}

func (a *ArticleHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Server is healthy and running",
		"status":  "OK",
	})
}

func (a *ArticleHandler) FetchArticles(c *gin.Context) {
	ctx := c.Request.Context()
	
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")
	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	filter := domain.ArticleFilter{
		Title:      c.Query("title"),
		Journal:    c.Query("journal"),
		AccessType: c.Query("accessType"),
		Publisher:  c.Query("publisher"),
		AuthorName: c.Query("authorName"),
		StartDate:  c.Query("startDate"),
		EndDate:    c.Query("endDate"),
		KeyWord:    c.Query("keyWord"),
		Limit:      limit,
		Offset:     offset,
	}
	
	articles, total, err := a.ArticleUseCase.Fetch(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  articles,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (a *ArticleHandler) FetchUniqueKeyWords(c *gin.Context) {
	ctx := c.Request.Context()

	keywords, err := a.ArticleUseCase.GetUniqueKeyWords(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  keywords,
		"total": len(keywords),
	})
}

type AddViewsRequest struct {
	Views int `json:"views" binding:"required"`
}

func (a *ArticleHandler) AddViews(c *gin.Context) {
	articleID := c.Param("id")
	if articleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "article id is required"})
		return
	}

	var req AddViewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	ctx := c.Request.Context()
	err := a.ArticleUseCase.AddViews(ctx, articleID, req.Views)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "views added successfully",
	})
}

// ProxyPDF - Tashqi URL'dan PDF'ni yuklab, inline ko'rinishda frontendga uzatadi va keshlaydi
func (a *ArticleHandler) ProxyPDF(c *gin.Context) {
	pdfURL := c.Query("url")
	if pdfURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
		return
	}

	// URL xavfsizligini tekshirish
	parsedURL, err := url.Parse(pdfURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid URL"})
		return
	}

	// Keshga saqlamasdan to'g'ridan-to'g'ri foydalanuvchiga uzatish
	resp, err := http.Get(pdfURL)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch PDF"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "external server returned error"})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "inline")
	c.Status(http.StatusOK)
	
	io.Copy(c.Writer, resp.Body)
}
