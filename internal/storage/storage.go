package storage

import (
	"context"
	"fmt"
	"web-server/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Storage struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func New(Conn *pgxpool.Pool, l *zap.Logger) *Storage {
	return &Storage{
		db:     Conn,
		logger: l,
	}
}

func GetConnect(connectStr string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), connectStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %v", err)
	}

	return db, nil
}

func (d *Storage) GetAllNotes() ([]models.Note, error) {
	query := `SELECT * FROM notes`

	rows, err := d.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		err := rows.Scan(&note.ID, &note.Title, &note.Content)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil

}

func AddNote() {

}

func GetNote() {

}

func UpdateNote() {

}

func DeleteNote() {

}
