package services

import (
	"Anastasia/songs/internal/models"
	"Anastasia/songs/internal/repository"
)

type SongService struct {
	repo *repository.Repo
}

func NewSongService(repo *repository.Repo) *SongService {
	return &SongService{
		repo: repo,
	}
}

func (s *SongService) Songs(filters models.Songs, page, pageSize int) ([]models.Songs, error) {
	return s.repo.Songs.Songs(filters, page, pageSize)
}

func (s *SongService) SongByID(id int) (string, error) {
	return s.repo.Songs.SongByID(id)
}
func (s *SongService) DeleteSong(id int) error {
	return s.repo.Songs.DeleteSong(id)
}
func (s *SongService) UpdateSong(song models.Songs) error {
	return s.repo.Songs.UpdateSong(song)
}
func (s *SongService) CreateSong(song models.Songs) error {
	return s.repo.Songs.CreateSong(song)
}
