package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"test-kittiphob/backend/internal/handler"
	"test-kittiphob/backend/internal/repository"
)

func main() {
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatal(err)
	}

	repo, err := repository.New("data/queue.db")
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	queue := handler.NewQueueHandler(repo)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/queues", func(r chi.Router) {
		r.Post("/", queue.Issue)
		r.Get("/current", queue.Current)
		r.Post("/reset", queue.Reset)
		r.Get("/{id}", queue.GetTicket)
	})

	log.Println("API running on http://localhost:4000")
	log.Fatal(http.ListenAndServe(":4000", r))
}
