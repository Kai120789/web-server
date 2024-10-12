package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"web-server/internal/dto"
	"web-server/internal/models"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	service Handlerer
	logger  *zap.Logger
}

type Handlerer interface {
	GetAllNotes() (*[]models.Note, error)
	AddNote(body dto.Dto) (*models.Note, error)
	GetNote(id uint) (*models.Note, error)
	UpdateNote(body dto.Dto, id uint) (*models.Note, error)
	DeleteNote(id uint) error
}

func New(s Handlerer, l *zap.Logger) Handler {
	return Handler{
		service: s,
		logger:  l,
	}
}

func (h *Handler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.service.GetAllNotes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notes)
}

func (h *Handler) AddNote(w http.ResponseWriter, r *http.Request) {
	var note dto.Dto
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if note.Title == "" {
		http.Error(w, "Note title is empty", http.StatusBadRequest)
		return
	}

	noteRet, err := h.service.AddNote(note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(noteRet)
}

func (h *Handler) GetNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid ID: %s, error: %v", idStr, err), http.StatusBadRequest)
		return
	}

	note, err := h.service.GetNote(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if note == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(note)
}

func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var note dto.Dto
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if note.Title == "" {
		http.Error(w, "Note title is empty", http.StatusBadRequest)
		return
	}

	noteRet, err := h.service.UpdateNote(note, uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(noteRet)
}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteNote(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
