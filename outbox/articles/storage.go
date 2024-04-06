package articles

import (
	"slices"
	"sync"
)

// InMemoryArticles is our lazy database attempt. It's also an adapter for our
// articles storage.
type InMemoryArticles struct {
	// we want to seem sophisticated
	mu sync.Mutex

	// a not-so-efficient wannabe postgres database
	articles []Article
	outbox   []ArticleLikedEvent
}

func NewInMemoryArticles() *InMemoryArticles {
	return &InMemoryArticles{
		mu:       sync.Mutex{},
		articles: make([]Article, 0, 64),
		outbox:   make([]ArticleLikedEvent, 0, 64),
	}
}

func (repo *InMemoryArticles) Insert(article Article) (Article, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.articles = append(repo.articles, article)
	return article, nil
}

func (repo *InMemoryArticles) Get(articleId int) (Article, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	idx := slices.IndexFunc(repo.articles, func(a Article) bool {
		return a.Id == articleId
	})
	if idx == -1 {
		return Article{}, ErrArticleNotFound
	}
	return repo.articles[idx], nil
}

func (repo *InMemoryArticles) Update(article Article) (Article, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	idx := slices.IndexFunc(repo.articles, func(a Article) bool {
		return a.Id == article.Id
	})
	if idx == -1 {
		return Article{}, ErrArticleNotFound
	}
	repo.articles[idx] = article
	return article, nil
}

// UpdateArticleAndInsertLikedEvents commits the changes to an article's state
// atomically with the events it emitted.
func (repo *InMemoryArticles) UpdateArticleAndInsertLikedEvents(
	article Article, events []ArticleLikedEvent) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	idx := slices.IndexFunc(repo.articles, func(a Article) bool {
		return a.Id == article.Id
	})
	if idx == -1 {
		return ErrArticleNotFound
	}
	repo.articles[idx] = article

	for _, e := range events {
		repo.outbox = append(repo.outbox, e)
	}
	return nil
}

func (repo *InMemoryArticles) GetArticleLikedEventsFromIndex(index int) ([]ArticleLikedEvent, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	return slices.Clone(repo.outbox[index+1:]), nil
}
