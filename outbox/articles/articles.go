package articles

import (
	"errors"
	"math/rand/v2"
	"time"
)

// Article is an aggregate encapsulating an article's state. It exposes methods
// to mutate the article which ensure an article's consistency.
type Article struct {
	Id          int
	AuthorId    int
	Content     string
	Likes       int
	PublishedAt time.Time
	ModifiedAt  time.Time

	likedEvents []ArticleLikedEvent
}

// Like should be called when an author's article receives a like.
func (a *Article) Like() {
	a.Likes++

	a.likedEvents = append(a.likedEvents, ArticleLikedEvent{
		ArticleId: a.Id,
		Timestamp: time.Now(),
		EventId:   rand.Int(),
	})
}

// GetLikedEvents returns all liked events that were generated when performing
// some business logic through the article aggregate.
func (a *Article) GetLikedEvents() []ArticleLikedEvent {
	return a.likedEvents
}

// UpdateContent should be called when an author updates his article.
func (a *Article) UpdateContent(newContent string) {
	a.Content = newContent
	a.ModifiedAt = time.Now()
}

// ArticleLikedEvent models the payload of an event that is emitted when an
// author's article receives a like.
type ArticleLikedEvent struct {
	EventId   int       `json:"event_id"`
	ArticleId int       `json:"article_id"`
	Timestamp time.Time `json:"timestamp"`
}

// ArticleStorage is a port which represents the required interactions between
// the application and a data source containing all articles
type ArticleStorage interface {
	// Insert inserts a new article, or returns an error if the article already
	// exists.
	Insert(a Article) (Article, error)
	// Get retrieves an article, or returns an error if the article does not
	// exist.
	Get(id int) (Article, error)
	// Update updates an article that matches the provided id, or returns an
	// error if there is no such article.
	Update(a Article) (Article, error)
	// UpdateArticleAndInsertLikedEvents commits the changes to an article's
	// state atomically with the events it emitted.
	UpdateArticleAndInsertLikedEvents(a Article, events []ArticleLikedEvent) error
	// GetArticleLikedEventsFromIndex returns all events starting from a
	// specific index.
	GetArticleLikedEventsFromIndex(i int) ([]ArticleLikedEvent, error)
}

// very simple error model
var (
	ErrArticleNotFound      = errors.New("article not found")
	ErrArticleAlreadyExists = errors.New("article already exists")
)
