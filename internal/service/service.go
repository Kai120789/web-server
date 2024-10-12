package service

import (
	"web-server/internal/dto"
	"web-server/internal/models"
)

type Service struct {
	storage Storager
}

type Storager interface {
	GetAllNotes() (*[]models.Note, error)
	AddNote(body dto.Dto) (*models.Note, error)
	GetNote(id uint) (*models.Note, error)
	UpdateNote(body dto.Dto, id uint) (*models.Note, error)
	DeleteNote(id uint) error
}

func New(s Storager) *Service {
	return &Service{
		storage: s,
	}
}

func (s *Service) GetAllNotes() (*[]models.Note, error) {
	notes, err := s.storage.GetAllNotes()
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (s *Service) AddNote(body dto.Dto) (*models.Note, error) {
	note, err := s.storage.AddNote(body)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (s *Service) GetNote(id uint) (*models.Note, error) {
	note, err := s.storage.GetNote(id)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (s *Service) UpdateNote(body dto.Dto, id uint) (*models.Note, error) {
	note, err := s.storage.UpdateNote(body, id)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (s *Service) DeleteNote(id uint) error {
	err := s.storage.DeleteNote(id)
	if err != nil {
		return err
	}

	return nil
}
