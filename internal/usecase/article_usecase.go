package usecase

import (
	"context"
	"maqola-backent/internal/domain"
	"time"
)

type articleUseCase struct {
	articleRepo    domain.ArticleRepository
	contextTimeout time.Duration
}

// NewArticleUseCase - UseCase qatlamini initsializatsiya qilish
func NewArticleUseCase(a domain.ArticleRepository, timeout time.Duration) domain.ArticleUseCase {
	return &articleUseCase{
		articleRepo:    a,
		contextTimeout: timeout,
	}
}

// Fetch - Barcha maqolalarni olish biznes mantig'i
func (a *articleUseCase) Fetch(c context.Context) ([]domain.Article, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	return a.articleRepo.Fetch(ctx)
}
