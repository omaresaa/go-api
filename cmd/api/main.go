package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/omaresaa/go-api/internal/handlers"
	"github.com/omaresaa/go-api/internal/storage"
)

func main() {
	storage := storage.NewJSONStorage("data/tasks.json")
	taskHandler := handlers.NewTaskHandler(storage)
	r := chi.NewRouter()

	// middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", taskHandler.GetAllTasks)
		r.Post("/", taskHandler.CreateTask)
		r.Get("/{id}", taskHandler.GetTask)
		r.Put("/{id}", taskHandler.UpdateTask)
		r.Delete("/{id}", taskHandler.DeleteTask)
	})

	log.Println("Server starting on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
