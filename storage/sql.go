package storage

import (
	"github.com/jinzhu/gorm"
	"log"
)

type sqlStorage struct {
	db *gorm.DB
}

func (s *sqlStorage) GetTracks(userId int) []Track {
	var tracks []Track
	s.db.Where("user_id = ?", userId).Find(&tracks)
	if s.db.Error != nil {
		log.Println(s.db.Error)
	}

	return tracks
}

func (s *sqlStorage) GetEvents(trackId string) ([]Event, error) {
	var events []Event
	s.db.Where("track_id = ?", trackId).Find(&events)
	if s.db.Error != nil {
		log.Println(s.db.Error)
	}

	return events, nil
}

func (s *sqlStorage) AddTrack(userId int, number string) error {
	track := Track{
		UserId: userId,
		Number: number,
	}
	s.db.Create(&track)

	return s.db.Error
}

func (s *sqlStorage) getTrack(userId int, number string) *Track {
	track := Track{}
	notFound := s.db.Where("user_id = ? and number = ?", userId, number).Find(&track).RecordNotFound()
	if s.db.Error != nil {
		log.Println(s.db.Error)
	}
	if notFound == true {
		return nil
	}

	return &track
}

func (s *sqlStorage) SetHistory(trackId string, events []Event) error {
	s.db.Where("track_id = ?", trackId).Delete(&Event{})
	for _, event := range events {
		event.TrackId = trackId
		s.db.Create(&event)
		if s.db.Error != nil {
			log.Println(s.db.Error)
		}
	}

	return s.db.Error
}

func (s *sqlStorage) GetAllTracks() []Track {
	var tracks []Track
	s.db.Find(&tracks)
	if s.db.Error != nil {
		log.Println(s.db.Error)
	}

	return tracks
}

func NewSql(db *gorm.DB) Storage {
	return &sqlStorage{db: db}
}
