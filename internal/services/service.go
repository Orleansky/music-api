package services

import (
	"Anastasia/songs/internal/models"
	"Anastasia/songs/internal/repository"
)

type Songs interface {
	Songs(filters models.Songs, page, pageSize int) ([]models.Songs, error)
	SongByID(id int) (string, error)
	DeleteSong(id int) error
	UpdateSong(song models.Songs) error
	CreateSong(song models.Songs) error
}

type Service struct {
	Songs
}

func NewService(repo *repository.Repo) *Service {
	service := &Service{
		Songs: NewSongService(repo),
	}
	return service
}
