package services

import (
	"Anastasia/songs/internal/models"
	"Anastasia/songs/internal/repository"
)

type SongService struct {
	repo *repository.Repo
}

// Создаёт новый экземпляр сервиса
func NewSongService(repo *repository.Repo) *SongService {
	return &SongService{
		repo: repo,
	}
}

// Получение данных библиотеки с фильтрацией по всем полям и пагинацией
func (s *SongService) Songs(filters models.Songs, page, pageSize int) ([]models.Songs, error) {
	return s.repo.Songs.Songs(filters, page, pageSize)
}

// Получение текста песни с пагинацией по куплетам
func (s *SongService) SongByID(id int) (string, error) {
	return s.repo.Songs.SongByID(id)
}

// Удаление песни
func (s *SongService) DeleteSong(id int) error {
	return s.repo.Songs.DeleteSong(id)
}

// Изменение данных песни
func (s *SongService) UpdateSong(song models.Songs) error {
	return s.repo.Songs.UpdateSong(song)
}

// Добавление новой песни
func (s *SongService) CreateSong(song models.Songs) error {
	return s.repo.Songs.CreateSong(song)
}
