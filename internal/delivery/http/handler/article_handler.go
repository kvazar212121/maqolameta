package handler

import (
	"net/http"

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
}

func (a *ArticleHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Server is healthy and running",
		"status":  "OK",
	})
}

func (a *ArticleHandler) FetchArticles(c *gin.Context) {
	ctx := c.Request.Context()
	
	filter := domain.ArticleFilter{
		Title:      c.Query("title"),
		Journal:    c.Query("journal"),
		AccessType: c.Query("accessType"),
		Publisher:  c.Query("publisher"),
		AuthorName: c.Query("authorName"),
		StartDate:  c.Query("startDate"),
		EndDate:    c.Query("endDate"),
		KeyWord:    c.Query("keyWord"),
	}
	
	articles, err := a.ArticleUseCase.Fetch(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  articles,
		"total": len(articles),
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
