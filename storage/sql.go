package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type sqlStorage struct {
	logger logrus.FieldLogger
	db     *gorm.DB
}

func (s *sqlStorage) GetTracks(userId int) []Track {
	var tracks []Track
	err := s.db.Where("user_id = ?", userId).Find(&tracks).Error
	if err != nil {
		s.logger.WithFields(logrus.Fields{"userId": userId, "err": s.db.Error}).Error("error getting tracks for user")
	}

	return tracks
}

func (s *sqlStorage) GetEvents(trackId string) (events []Event, err error) {
	err = s.db.Where("track_id = ?", trackId).Find(&events).Error
	if err != nil {
		s.logger.WithFields(logrus.Fields{"trackId": trackId, "err": err}).Error("error getting events for track")
		return
	}

	return events, nil
}

func (s *sqlStorage) AddTrack(userId int, number string, name string) error {
	track := Track{
		UserId: userId,
		Number: number,
		Name:   name,
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
		s.logger.WithFields(logrus.Fields{"trackId": number, "userId": userId, "err": s.db.Error}).Error("error getting track")
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
		s.logger.WithFields(logrus.Fields{"err": err}).Error("error getting all tracks")
	}

	return tracks
}

func NewSql(db *gorm.DB, log logrus.FieldLogger) *sqlStorage {
	return &sqlStorage{
		db:     db,
		logger: log,
	}
}
