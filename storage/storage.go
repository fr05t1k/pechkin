package storage

import "time"

type Track struct {
	Id string
}

type Event struct {
	When        time.Time
	Description []string
}

type Storage interface {
	GetTracks(userId int) []Track
	GetEvents(userId int, trackId string) ([]Event, error)
	AddTrack(userId int, trackId string) error
	SetHistory(trackId string, events []Event) error
}
