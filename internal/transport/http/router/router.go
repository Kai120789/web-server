package router

import (
	"net/http"
	"web-server/internal/transport/http/handler"

	"github.com/go-chi/chi/v5"
)

type Router interface {
	GetAllNotes(w http.ResponseWriter, r *http.Request)
	AddNote(w http.ResponseWriter, r *http.Request)
	GetNote(w http.ResponseWriter, r *http.Request)
	UpdateNote(w http.ResponseWriter, r *http.Request)
	DeleteNote(w http.ResponseWriter, r *http.Request)
}

func New(h *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Route("/api/notes", func(r chi.Router) {
		r.Get("/", h.GetAllNotes)
		r.Post("/", h.AddNote)
		r.Get("/{id}", h.GetNote)
		r.Put("/{id}", h.UpdateNote)
		r.Delete("/{id}", h.DeleteNote)
	})

	return r
}
