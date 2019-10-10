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
	err := s.db.Where("user_id = ?", userId).Find(&tracks).Error
	if err != nil {
		log.Println(s.db.Error)
	}

	return tracks
}

func (s *sqlStorage) GetEvents(trackId string) (events []Event, err error) {
	err = s.db.Where("track_id = ?", trackId).Find(&events).Error
	if err != nil {
		log.Println(err)
		return
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
	err := s.db.Where("user_id = ? and number = ?", userId, number).Find(&track).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	if err != nil {
		log.Println(s.db.Error)
		return nil
	}

	return &track
}

func (s *sqlStorage) SetHistory(trackId string, events []Event) error {
	s.db.Where("track_id = ?", trackId).Delete(&Event{})
	for _, event := range events {
		event.TrackId = trackId
		err := s.db.Create(&event).Error
		if err != nil {
			return err
		}
	}

	return s.db.Error
}

func (s *sqlStorage) GetAllTracks() []Track {
	var tracks []Track
	err := s.db.Find(&tracks).Error
	if err != nil {
		log.Println(err)
	}

	return tracks
}

func NewSql(db *gorm.DB) *sqlStorage {
	return &sqlStorage{db: db}
}
