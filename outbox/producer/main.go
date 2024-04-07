package main

import (
	"fmt"
	"net/http"
	"os"
	"outbox/articles"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var storage articles.ArticleStorage

func LikeArticle(articleId int) error {

	// load the aggregate state and apply bussiness logic
	article, err := storage.Get(articleId)
	if err != nil {
		return fmt.Errorf("could not retrieve article")
	}

	article.Like()
	events := article.GetLikedEvents()

	// commit the aggregate's state and domain events atomically to the
	// database
	err = storage.UpdateArticleAndInsertLikedEvents(article, events)
	if err != nil {
		return fmt.Errorf("could not update article state after like: %s", err.Error())
	}
	log.Info().Int("id", events[0].EventId).Msg("received like")
	return nil
}

type handler struct {
	storage articles.ArticleStorage
}

func NewHandler(storage articles.ArticleStorage) handler {
	return handler{
		storage: storage,
	}
}

func main() {

	portstr := os.Args[1]
	if portstr == "" {
		portstr = "13311"
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Initialize the storage. Based on config, this could instantiate
	// different types of storage, e.g. postgres
	storage = articles.NewInMemoryArticles()
	requestHandler := NewHandler(storage)

	a := articles.Article{
		Id:          1,
		AuthorId:    1,
		Content:     "bla",
		Likes:       3,
		PublishedAt: time.Now(),
		ModifiedAt:  time.Now(),
	}
	storage.Insert(a)
	// err := LikeArticle(a.Id)
	// if err != nil {
	// 	log.Error().Err(err).Msg("could not like article")
	// }

	// Initialize poller
	poller, err := articles.NewLikedArticlesPoller(storage, 5*time.Second, "localhost:"+portstr)
	if err != nil {
		fmt.Printf("could not start poller: %s\n", err.Error())
		os.Exit(1)
	}
	// Start the poller in a new goroutine.
	poller.Poll()
	log.Info().Msg("started poller")

	http.HandleFunc("POST /articles/{id}/like", requestHandler.HandleLikeArticle)
	// other handlers

	log.Info().Msg("started server")
	http.ListenAndServe("127.0.0.1:18889", nil)
}

func (h *handler) HandleLikeArticle(w http.ResponseWriter, r *http.Request) {
	articleIdStr := r.PathValue("id")
	articleId, err := strconv.Atoi(articleIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = LikeArticle(articleId)

	if err != nil {
		fmt.Printf("got error: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ":(")
		return
	}
}
