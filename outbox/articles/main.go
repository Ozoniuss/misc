package main

import (
	"fmt"
	"net/http"
	"strconv"
)

var storage ArticleStorage

func LikeArticle(articleId int) error {

	// load the aggregate state and apply bussiness logic
	article, err := storage.Get(articleId)
	if err != nil {
		return fmt.Errorf("could not retrieve article")
	}

	article.Like()
	events := article.getLikedEvents()

	// commit the aggregate's state and domain events atomically to the
	// database
	err = storage.UpdateArticleAndInsertLikedEvents(article, events)
	if err != nil {
		return fmt.Errorf("could not update article state after like: %s", err.Error())
	}
	return nil
}

type handler struct {
	storage ArticleStorage
}

func NewHandler(storage ArticleStorage) handler {
	return handler{
		storage: storage,
	}
}

func main() {
	// Initialize the storage. Based on config, this could instantiate
	// different types of storage, e.g. postgres
	storage := NewInMemoryArticles()
	requestHandler := NewHandler(storage)

	http.HandleFunc("POST /articles/{id}/like", requestHandler.HandleLikeArticle)
	// other handlers
}

func (h *handler) HandleLikeArticle(w http.ResponseWriter, r *http.Request) {
	articleIdStr := r.PathValue("id")
	articleId, err := strconv.Atoi(articleIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = LikeArticle(articleId)

	// staff level engineer error handling
	if err != nil {
		fmt.Printf("got error: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ":(")
		return
	}
}
