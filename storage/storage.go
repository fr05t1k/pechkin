package storage

import (
	"errors"
	"time"
)

type Track struct {
	ID        int    `gorm:"primary_key"`
	UserId    int    `gorm:"unique_index:idx_user_id_number"`
	Number    string `gorm:"index:idx_number;unique_index:idx_user_id_number"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Event struct {
	ID          int    `gorm:"primary_key"`
	TrackId     string `gorm:"index:idx_track_id"`
	EventAt     time.Time
	Description string
	CreatedAt   time.Time
}

type User struct {
	ID         int `gorm:"primary_key"`
	TrackLimit int
}

var NotFound = errors.New("track not found")

type Storage interface {
	GetTracks(userId int) []Track
	GetEvents(number string) ([]Event, error)
	AddTrack(userId int, number string, name string) error
	SetHistory(number string, events []Event) error
	GetAllTracks() []Track
	Remove(number string) error
	GetTrackForUser(number string, userId int) (Track, error)
	GetTrackByNumber(number string) (users []Track, err error)
	IsLimitExceeded(userId int) (bool, error)
}
