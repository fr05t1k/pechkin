package storage

import (
	"time"
)

type Track struct {
	ID        int    `gorm:"primary_key"`
	UserId    int    `gorm:"unique_index:idx_user_id_number"`
	Number    string `gorm:"index:idx_number;unique_index:idx_user_id_number"`
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

type Storage interface {
	GetTracks(userId int) []Track
	GetEvents(trackId string) ([]Event, error)
	AddTrack(userId int, trackId string) error
	SetHistory(trackId string, events []Event) error
	GetAllTracks() []Track
}
