package api

import (
	"Anastasia/songs/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// @Summary		Get all songs
// @Description	Get a list of all songs with optional filters
// @Tags			songs
// @Accept			json
// @Produce		json
// @Param			group		query		string	false	"Group filter"
// @Param			song		query		string	false	"Song filter"
// @Param			releaseDate	query		string	false	"Release date filter"
// @Param			text		query		string	false	"Text filter"
// @Param			link		query		string	false	"Link filter"
// @Param			page		query		int		false	"Page number"
// @Param			pageSize	query		int		false	"Page size"
// @Success		200			{array}		models.Songs
// @Failure		500			{object}	string
// @Router			/songs [get]
func (api *API) songsHandler(w http.ResponseWriter, r *http.Request) {
	filters := models.Songs{
		Group:       r.URL.Query().Get("group"),
		Song:        r.URL.Query().Get("song"),
		ReleaseDate: r.URL.Query().Get("releaseDate"),
		Text:        r.URL.Query().Get("text"),
		Link:        r.URL.Query().Get("link"),
	}

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	pageSizeStr := r.URL.Query().Get("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10
	}

	logrus.WithFields(logrus.Fields{
		"filters":  filters,
		"page":     page,
		"pageSize": pageSize,
	}).Info("Fetching songs")

	songs, err := api.srv.Songs.Songs(filters, page, pageSize)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch songs")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(songs)
	if err != nil {
		logrus.WithError(err).Error("Failed to encode songs to JSON")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// @Summary		Get a song by ID
// @Description	Get a lyrics of the song by its ID with optional verse number
// @Tags			songs
// @Accept			json
// @Produce		json
// @Param			id		path		int	true	"Song ID"
// @Param			verse	query		int	false	"Verse number"
// @Success		200		{string}	string
// @Failure		400		{object}	string
// @Failure		500		{object}	string
// @Router			/songs/{id} [get]
func (api *API) songByIDHandler(w http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["id"]
	verse := r.URL.Query().Get("verse")
	verseNum, err := strconv.Atoi(verse)
	if err != nil {
		verseNum = 0
	}
	id, err := strconv.Atoi(s)
	if err != nil {
		logrus.WithError(err).Error("Invalid song ID")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logrus.WithField("id", id).Info("Fetching song by ID")

	lyrics, err := api.srv.SongByID(id)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch song by ID")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lyricsPaginated := strings.Split(lyrics, "\n\n")

	// Если значение verseNum некорректно, метод будет выводить полный текст песни
	if verseNum == 0 || verseNum > len(lyricsPaginated) {
		err = json.NewEncoder(w).Encode(lyrics)
	} else {
		err = json.NewEncoder(w).Encode(lyricsPaginated[verseNum-1])
	}
	if err != nil {
		logrus.WithError(err).Error("Failed to encode lyrics to JSON")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// @Summary		Delete a song by ID
// @Description	Delete a song by its ID
// @Tags			songs
// @Accept			json
// @Produce		json
// @Param			id	path	int	true	"Song ID"
// @Success		204	"No Content"
// @Failure		400	{object}	string
// @Failure		500	{object}	string
// @Router			/songs/{id} [delete]
func (api *API) deleteSongHandler(w http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(s)
	if err != nil {
		logrus.WithError(err).Error("Invalid song ID")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logrus.WithField("id", id).Info("Deleting song")

	err = api.srv.DeleteSong(id)
	if err != nil {
		logrus.WithError(err).Error("Failed to delete song")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// @Summary		Update a song by ID
// @Description	Update a song by its ID
// @Tags			songs
// @Accept			json
// @Produce		json
// @Param			id		path		int				true	"Song ID"
// @Param			song	body		models.Songs	true	"Song object"
// @Success		200		{object}	models.Songs
// @Failure		400		{object}	string
// @Failure		500		{object}	string
// @Router			/songs/{id} [patch]
func (api *API) updateSongHandler(w http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(s)
	if err != nil {
		logrus.WithError(err).Error("Invalid song ID")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var song models.Songs

	err = json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		logrus.WithError(err).Error("Failed to decode song data")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	song.ID = id
	logrus.WithField("song", song).Info("Updating song")

	err = api.srv.UpdateSong(song)
	if err != nil {
		logrus.WithError(err).Error("Failed to update song")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// @Summary		Create a new song
// @Description	Create a new song
// @Tags			songs
// @Accept			json
// @Produce		json
// @Param			song	body		models.Songs	true	"Song object"
// @Success		201		{object}	models.Songs
// @Failure		400		{object}	string
// @Failure		500		{object}	string
// @Router			/songs [post]
func (api *API) createSongHandler(w http.ResponseWriter, r *http.Request) {
	var song models.Songs

	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		logrus.WithError(err).Error("Failed to decode song data")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	encodedGroup := url.QueryEscape(song.Group)
	encodedSong := url.QueryEscape(song.Song)

	apiURL := os.Getenv("EXTERNAL_API_URL") + fmt.Sprintf("?group=%s&song=%s", encodedGroup, encodedSong)
	logrus.WithField("apiURL", apiURL).Info("Requesting data from external API")

	resp, err := http.Get(apiURL)
	if err != nil {
		logrus.WithError(err).Error("Failed to get data from external API")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	logrus.WithField("status", resp.StatusCode).Info("Received response from external API")

	if resp.StatusCode != http.StatusOK {
		logrus.WithField("status", resp.StatusCode).Error("External API returned non-OK status")
		http.Error(w, "External API returned non-OK status", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Error("Failed to read response body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.WithField("responseBody", string(body)).Info("Response body from external API")

	if err := json.Unmarshal(body, &song); err != nil {
		logrus.WithError(err).Error("Failed to unmarshal song data")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.WithField("song", song).Info("Creating song")

	err = api.srv.CreateSong(song)
	if err != nil {
		logrus.WithError(err).Error("Failed to create song")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
