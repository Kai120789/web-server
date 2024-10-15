package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"web-server/internal/dto"
	"web-server/internal/models"
)

type Storage struct {
	filePath string
}

func New(filePath string) *Storage {
	return &Storage{
		filePath: filePath,
	}
}

func (s *Storage) readNotes() ([]models.Note, error) {
	file, err := os.ReadFile(s.filePath)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	// Проверка на пустой файл
	if len(file) == 0 {
		return []models.Note{}, nil
	}

	var notes []models.Note
	err = json.Unmarshal(file, &notes)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return notes, nil
}

func (s *Storage) writeNotes(notes []models.Note) error {
	data, err := json.Marshal(notes)
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, os.ModePerm)
}

func (s *Storage) GetAllNotes() (*[]models.Note, error) {
	notes, err := s.readNotes()
	if err != nil {
		return nil, err
	}
	return &notes, nil
}

func (s *Storage) AddNote(body dto.Dto) (*models.Note, error) {
	notes, err := s.readNotes()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	id := uint(len(notes) + 1) // простая автоинкрементация
	newNote := models.Note{
		ID:      id,
		Title:   body.Title,
		Content: body.Content,
	}
	notes = append(notes, newNote)

	err = s.writeNotes(notes)
	if err != nil {
		return nil, err
	}

	return &newNote, nil
}

func (s *Storage) GetNote(id uint) (*models.Note, error) {
	notes, err := s.readNotes()
	if err != nil {
		return nil, err
	}

	for _, note := range notes {
		if note.ID == id {
			return &note, nil
		}
	}

	return nil, nil
}

func (s *Storage) UpdateNote(body dto.Dto, id uint) (*models.Note, error) {
	notes, err := s.readNotes()
	if err != nil {
		return nil, err
	}

	for i, note := range notes {
		if note.ID == id {
			notes[i].Title = body.Title
			notes[i].Content = body.Content
			err = s.writeNotes(notes)
			if err != nil {
				return nil, err
			}
			return &notes[i], nil
		}
	}

	return nil, fmt.Errorf("note with id %d not found", id)
}

func (s *Storage) DeleteNote(id uint) error {
	notes, err := s.readNotes()
	if err != nil {
		return err
	}

	for i, note := range notes {
		if note.ID == id {
			notes = append(notes[:i], notes[i+1:]...)
			return s.writeNotes(notes)
		}
	}

	return fmt.Errorf("note with id %d not found", id)
}
