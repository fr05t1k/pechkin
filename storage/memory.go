package storage

import "errors"

type memory struct {
	tracks  map[int][]Track
	history map[string][]Event
}

var NoEventsError = errors.New("no events")

func (m *memory) GetTracks(userId int) (tracks []Track) {
	tracks = m.tracks[userId]
	return
}

func (m *memory) GetEvents(userId int, trackId string) ([]Event, error) {
	events, ok := m.history[trackId]
	if !ok {
		return nil, NoEventsError
	}
	return events, nil
}

func (m *memory) SetHistory(trackId string, events []Event) error {
	m.history[trackId] = events
	return nil
}

func (m *memory) AddTrack(userId int, trackId string) error {
	m.tracks[userId] = append(m.tracks[userId], Track{Id: trackId})
	return nil
}

func NewMemory() *memory {
	return &memory{
		tracks:  make(map[int][]Track),
		history: make(map[string][]Event),
	}
}
