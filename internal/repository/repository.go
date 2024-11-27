package repository

import (
	"Anastasia/songs/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Songs interface {
	Songs(filters models.Songs, page, pageSize int) ([]models.Songs, error)
	SongByID(id int) (string, error)
	DeleteSong(id int) error
	UpdateSong(song models.Songs) error
	CreateSong(song models.Songs) error
}

type Repo struct {
	Songs
}

func NewRepo(db *pgxpool.Pool) *Repo {
	repo := &Repo{
		Songs: NewSongRepo(db),
	}
	return repo
}
