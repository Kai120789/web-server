package storage

import (
	"context"
	"fmt"
	"web-server/internal/dto"
	"web-server/internal/models"

	"github.com/jackc/pgx/v5"
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

func (d *Storage) GetAllNotes() (*[]models.Note, error) {
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

	return &notes, nil

}

func (d *Storage) AddNote(body dto.Dto) (*models.Note, error) {
	query := `INSERT INTO notes (title, content) VALUES ($1, $2) RETURNING id`

	var id uint

	err := d.db.QueryRow(context.Background(), query, body.Title, body.Content).Scan(&id)
	if err != nil {
		return nil, err
	}

	noteRet, err := d.GetNote(id)
	if err != nil {
		return nil, err
	}

	return noteRet, nil
}

func (d *Storage) GetNote(id uint) (*models.Note, error) {
	query := `SELECT * FROM notes WHERE id=$1`
	row := d.db.QueryRow(context.Background(), query, id)

	var note models.Note
	err := row.Scan(&note.ID, &note.Title, &note.Content)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &note, nil
}

func (d *Storage) UpdateNote(body dto.Dto, id uint) (*models.Note, error) {
	query := `UPDATE notes SET title=$1, content=$2 WHERE id=$3`
	_, err := d.db.Exec(context.Background(), query, body.Title, body.Content, id)
	if err != nil {
		return nil, err
	}

	boardRet, err := d.GetNote(id)
	if err != nil {
		return nil, err
	}

	return boardRet, nil
}

func (d *Storage) DeleteNote(id uint) error {
	query := `DELETE FROM notes WHERE id=$1`
	_, err := d.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}
