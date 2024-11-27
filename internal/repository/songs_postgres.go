package repository

import (
	"Anastasia/songs/internal/models"
	"context"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type SongRepo struct {
	db *pgxpool.Pool
}

// Создаёт новый экземпляр репозитория
func NewSongRepo(db *pgxpool.Pool) *SongRepo {
	return &SongRepo{
		db: db,
	}
}

// Получение данных библиотеки с фильтрацией по всем полям и пагинацией
func (s *SongRepo) Songs(filters models.Songs, page, pageSize int) ([]models.Songs, error) {
	logrus.WithFields(logrus.Fields{
		"filters":  filters,
		"page":     page,
		"pageSize": pageSize,
	}).Debug("Fetching songs with filters")

	rows, err := s.db.Query(context.Background(), `
		SELECT s.id, s.name, g.name, s.release_date, s.text, s.link
		FROM songs s
		INNER JOIN groups g ON s.group_id = g.id
		WHERE s.name ILIKE $1 AND
		g.name ILIKE $2 AND
		s.release_date ILIKE $3 AND
		s.text ILIKE $4 AND
		s.link ILIKE $5
		LIMIT $6
		OFFSET $7
	`, "%"+filters.Song+"%", "%"+filters.Group+"%", "%"+filters.ReleaseDate+"%", "%"+filters.Text+"%", "%"+filters.Link+"%",
		pageSize, (page-1)*pageSize)

	if err != nil {
		logrus.WithError(err).Error("Failed to query songs")
		return nil, err
	}
	defer rows.Close()

	songs := []models.Songs{}
	for rows.Next() {
		var song models.Songs
		err := rows.Scan(
			&song.ID,
			&song.Song,
			&song.Group,
			&song.ReleaseDate,
			&song.Text,
			&song.Link,
		)
		if err != nil {
			logrus.WithError(err).Error("Failed to scan song row")
			return nil, err
		}
		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		logrus.WithError(err).Error("Error occurred while iterating over rows")
		return nil, err
	}

	logrus.WithField("songs", songs).Debug("Fetched songs successfully")
	return songs, nil
}

// Получение текста песни с пагинацией по куплетам
func (s *SongRepo) SongByID(id int) (string, error) {
	logrus.WithField("id", id).Debug("Fetching song by ID")

	row := s.db.QueryRow(context.Background(), `
		SELECT s.text
		FROM songs s
		INNER JOIN groups g ON s.group_id = g.id
		WHERE s.id = $1
	`, id)

	var lyrics string
	err := row.Scan(&lyrics)
	if err != nil {
		logrus.WithError(err).Error("Failed to scan lyrics")
		return "", err
	}

	logrus.WithField("lyrics", lyrics).Debug("Fetched lyrics successfully")
	return lyrics, nil
}

// Удаление песни
func (s *SongRepo) DeleteSong(id int) error {
	logrus.WithField("id", id).Debug("Deleting song")

	row := s.db.QueryRow(context.Background(), `
		DELETE FROM songs
		WHERE id = $1
		RETURNING group_id;
	`, id)

	var groupId int
	err := row.Scan(&groupId)
	if err != nil {
		logrus.WithError(err).Error("Failed to delete song")
		return err
	}

	err = s.checkGroup(groupId)
	if err != nil {
		return err
	}

	logrus.WithField("groupId", groupId).Debug("Song deleted successfully")
	return nil
}

// Изменение данных песни
func (s *SongRepo) UpdateSong(song models.Songs) error {
	logrus.WithField("song", song).Debug("Updating song")

	var groupId int
	if song.Group != "" {
		row := s.db.QueryRow(context.Background(), `
			INSERT INTO groups (name)
			VALUES ($1)
			ON CONFLICT (name)
			DO NOTHING
			RETURNING id;
		`, song.Group)

		err := row.Scan(&groupId)
		if err != nil {
			if err.Error() == "no rows in result set" {
				err = s.db.QueryRow(context.Background(), `
					SELECT id FROM groups WHERE name = $1
				`, song.Group).Scan(&groupId)
				if err != nil {
					logrus.WithError(err).Error("Failed to get new group ID")
					return err
				}
			} else {
				logrus.WithError(err).Error("Failed to insert group")
				return err
			}
		}
	}

	row := s.db.QueryRow(context.Background(), `
	SELECT group_id FROM songs
	WHERE id = $1
`, song.ID)

	query := "UPDATE songs SET "
	var args []interface{}
	argIndex := 1

	if song.Song != "" {
		query += "name = $" + strconv.Itoa(argIndex)
		args = append(args, song.Song)
		argIndex++
	}

	if groupId != 0 {
		if argIndex > 1 {
			query += ", "
		}
		query += "group_id = $" + strconv.Itoa(argIndex)
		args = append(args, groupId)
		argIndex++
	}

	if song.ReleaseDate != "" {
		if argIndex > 1 {
			query += ", "
		}
		query += "release_date = $" + strconv.Itoa(argIndex)
		args = append(args, song.ReleaseDate)
		argIndex++
	}

	if song.Text != "" {
		if argIndex > 1 {
			query += ", "
		}
		query += "text = $" + strconv.Itoa(argIndex)
		args = append(args, song.Text)
		argIndex++
	}

	if song.Link != "" {
		if argIndex > 1 {
			query += ", "
		}
		query += "link = $" + strconv.Itoa(argIndex)
		args = append(args, song.Link)
		argIndex++
	}

	query += " WHERE id = $" + strconv.Itoa(argIndex)
	args = append(args, song.ID)

	_, err := s.db.Exec(context.Background(), query, args...)
	if err != nil {
		logrus.WithError(err).Error("Failed to update song")
		return err
	}

	var currentGroupId int
	err = row.Scan(&currentGroupId)
	if err != nil {
		logrus.WithError(err).Error("Failed to get current group ID")
	}

	err = s.checkGroup(currentGroupId)
	if err != nil {
		return err
	}

	logrus.WithField("song", song).Debug("Song updated successfully")
	return nil
}

// Добавление новой песни
func (s *SongRepo) CreateSong(song models.Songs) error {
	logrus.WithField("song", song).Debug("Creating song")

	row := s.db.QueryRow(context.Background(), `
		INSERT INTO groups (name)
		VALUES ($1)
		ON CONFLICT (name)
		DO NOTHING
		RETURNING id;
	`, song.Group)

	var groupId int
	err := row.Scan(&groupId)
	if err != nil {
		if err.Error() == "no rows in result set" {
			err = s.db.QueryRow(context.Background(), `
                SELECT id FROM groups WHERE name = $1
            `, song.Group).Scan(&groupId)
			if err != nil {
				logrus.WithError(err).Error("Failed to get group ID")
				return err
			}
		} else {
			logrus.WithError(err).Error("Failed to insert group")
			return err
		}
	}

	_, err = s.db.Exec(context.Background(), `
		INSERT INTO songs (name, group_id, release_date, text, link)
		VALUES ($1, $2, $3, $4, $5)
	`, song.Song, groupId, song.ReleaseDate, song.Text, song.Link)
	if err != nil {
		logrus.WithError(err).Error("Failed to insert song")
		return err
	}

	logrus.WithField("song", song).Debug("Song created successfully")
	return nil
}

func (s *SongRepo) checkGroup(groupId int) error {
	row := s.db.QueryRow(context.Background(), `
		SELECT COUNT(id) FROM songs
		WHERE group_id = $1
	`, groupId)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logrus.WithError(err).Error("Failed to count songs in group")
		return err
	}

	// Во избежание хранения избыточной информации в таблице groups удаляем неиспользуемые строки таблицы
	if count == 0 {
		_, err = s.db.Exec(context.Background(), `
			DELETE FROM groups
			WHERE id = $1
		`, groupId)
		if err != nil {
			logrus.WithError(err).Error("Failed to delete group")
			return err
		}
	}
	return nil
}
