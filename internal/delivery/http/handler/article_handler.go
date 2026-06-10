package handler

import (
	"net/http"

	"maqola-backent/internal/domain"

	"github.com/gin-gonic/gin"
)

// ArticleHandler HTTP so'rovlarni boshqaruvchi tuzilma
type ArticleHandler struct {
	ArticleUseCase domain.ArticleUseCase
}

// NewArticleHandler - yangi handler yaratadi va yo'nalishlarni (marshrutlarni) bog'laydi
func NewArticleHandler(r *gin.Engine, us domain.ArticleUseCase) {
	handler := &ArticleHandler{
		ArticleUseCase: us,
	}

	// Yo'nalishlar (Routes)
	r.GET("/ping", handler.Ping)
	r.GET("/api/v1/articles", handler.FetchArticles)
}

// Ping funksiyasi - server holatini tekshirish uchun
func (a *ArticleHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "backend to'g'ri ishlayapti",
	})
}

// FetchArticles funksiyasi - barcha maqolalarni qaytarish
func (a *ArticleHandler) FetchArticles(c *gin.Context) {
	ctx := c.Request.Context()
	
	articles, err := a.ArticleUseCase.Fetch(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  articles,
		"total": len(articles),
	})
}
