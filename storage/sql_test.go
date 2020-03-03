package storage

import (
	"github.com/fr05t1k/pechkin/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func createDb(t *testing.T) (db *gorm.DB) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if !assert.NoError(t, err) {
		return
	}
	db.LogMode(false)
	db.AutoMigrate(&Track{}, &Event{}, &User{})

	return db
}

// db without migrations
func createBrokenDb(t *testing.T) (db *gorm.DB) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if !assert.NoError(t, err) {
		return
	}
	db.LogMode(false)
	return db
}
func Test_sqlStorage_SetHistory(t *testing.T) {
	type args struct {
		trackId string
		events  []Event
	}
	tests := []struct {
		name    string
		db      *gorm.DB
		args    args
		wantErr bool
	}{
		{
			name: "one",
			db:   createDb(t),
			args: args{
				trackId: "ABC",
				events: []Event{{
					TrackId:     "ABC",
					EventAt:     time.Date(2016, 12, 31, 0, 0, 0, 0, time.UTC),
					Description: "test one",
				}},
			},
			wantErr: false,
		},
		{
			name: "two",
			db:   createDb(t),
			args: args{
				trackId: "ABC",
				events: []Event{
					{
						TrackId:     "ABC",
						EventAt:     time.Date(2016, 12, 31, 0, 0, 0, 0, time.UTC),
						Description: "test one",
					},
					{
						TrackId:     "ABC",
						EventAt:     time.Date(2016, 12, 31, 0, 0, 0, 0, time.UTC),
						Description: "test two",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			db:   createBrokenDb(t),
			args: args{
				trackId: "ABC",
				events: []Event{{
					TrackId:     "ABC",
					EventAt:     time.Date(2016, 12, 31, 0, 0, 0, 0, time.UTC),
					Description: "test one",
				}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSql(tt.db, log.NewNopLogger())
			if err := s.SetHistory(tt.args.trackId, tt.args.events); (err != nil) != tt.wantErr {
				t.Errorf("SetHistory() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == false {
				events, err := s.GetEvents(tt.args.trackId)

				if err != nil {
					t.Error(err)
					return
				}

				// check length
				if !assert.Equal(t, len(tt.args.events), len(events)) {
					t.Error("amount of events is not correct")
					return
				}

				// check fields
				for i, event := range events {
					assert.Equal(t, tt.args.events[i].Description, event.Description)
					assert.Equal(t, tt.args.events[i].TrackId, event.TrackId)
					assert.Equal(t, tt.args.events[i].EventAt, event.EventAt)
				}
			}
		})
	}
}

func Test_sqlStorage_GetAllTracks(t *testing.T) {
	s := NewSql(createDb(t), log.NewNopLogger())

	err := s.AddTrack(1, "ABC", "name a")
	assert.NoError(t, err)

	err = s.AddTrack(1, "EFG", "name b")
	assert.NoError(t, err)

	tracks := s.GetAllTracks()

	// check length
	if !assert.Len(t, tracks, 2) {
		return
	}
	// check necessary fields
	assert.Equal(t, 1, tracks[0].UserId)
	assert.Equal(t, "ABC", tracks[0].Number)
	assert.Equal(t, "name a", tracks[0].Name)
	assert.Equal(t, 1, tracks[1].UserId)
	assert.Equal(t, "EFG", tracks[1].Number)
	assert.Equal(t, "name b", tracks[1].Name)

	s = NewSql(createBrokenDb(t), log.NewNopLogger())

	tracks = s.GetAllTracks()

	assert.Equal(t, []Track{}, tracks)
}

func Test_sqlStorage_GetTracks(t *testing.T) {
	tests := []struct {
		name   string
		userId int
		want   []Track
		db     *gorm.DB
	}{
		{
			name:   "first",
			userId: 1,
			want: []Track{
				{
					UserId: 1,
					Number: "ABC",
				},
			},
			db: createDb(t),
		},
		{
			name:   "second",
			userId: 2,
			want: []Track{
				{
					UserId: 2,
					Number: "EFG",
				},
				{
					UserId: 2,
					Number: "XYZ",
				},
			},
			db: createDb(t),
		},
		{
			name:   "empty",
			userId: 3,
			want:   []Track{},
			db:     createDb(t),
		},
		{
			name:   "error",
			userId: 3,
			want:   []Track{},
			db:     createBrokenDb(t),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSql(tt.db, log.NewNopLogger())
			assert.NoError(t, s.AddTrack(1, "ABC", ""))
			assert.NoError(t, s.AddTrack(2, "EFG", ""))
			assert.NoError(t, s.AddTrack(2, "XYZ", ""))

			tracks := s.GetTracks(tt.userId)
			// check length
			if !assert.Equal(t, len(tracks), len(tt.want)) {
				return
			}

			// check necessary fields
			for i, track := range tracks {
				assert.Equal(t, tt.want[i].Number, track.Number)
				assert.Equal(t, tt.want[i].UserId, track.UserId)
			}
		})
	}
}

func Test_sqlStorage_getTrack(t *testing.T) {
	db := createDb(t)
	s := NewSql(db, log.NewNopLogger())
	assert.NoError(t, s.AddTrack(1, "ABC", ""))
	assert.NoError(t, s.AddTrack(2, "EFG", ""))
	assert.NoError(t, s.AddTrack(2, "XYZ", ""))

	brokenDb := createBrokenDb(t)
	type args struct {
		userId int
		number string
	}
	tests := []struct {
		name string
		args args
		db   *gorm.DB
		want *Track
	}{
		{
			name: "found",
			db:   db,
			args: args{
				userId: 1,
				number: "ABC",
			},
			want: &Track{
				UserId: 1,
				Number: "ABC",
			},
		},
		{
			name: "found 2",
			db:   db,
			args: args{
				userId: 2,
				number: "EFG",
			},
			want: &Track{
				UserId: 2,
				Number: "EFG",
			},
		},
		{
			name: "not found",
			db:   db,
			args: args{
				userId: 1,
				number: "EFG",
			},
			want: nil,
		},
		{
			name: "error",
			db:   brokenDb,
			args: args{
				userId: 1,
				number: "EFG",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSql(tt.db, log.NewNopLogger())

			track := s.getTrack(tt.args.userId, tt.args.number)
			if track == nil {
				assert.Equal(t, tt.want, track)
				return
			}
			assert.Equal(t, tt.want.UserId, track.UserId)
			assert.Equal(t, tt.want.Number, track.Number)
		})
	}
}

func Test_sqlStorage_GetEvents(t *testing.T) {
	db := createDb(t)
	s := NewSql(db, log.NewNopLogger())
	err := s.SetHistory(
		"ABC",
		[]Event{
			{Description: "abc 1"},
			{Description: "abc 2"},
		},
	)
	assert.NoError(t, err)
	brokenDb := createBrokenDb(t)
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		trackId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Event
		wantErr bool
	}{
		{
			name: "error",
			fields: fields{
				db: brokenDb,
			},
			args: args{
				trackId: "ABC",
			},
			want:    []Event{},
			wantErr: true,
		},
		{
			name: "found",
			fields: fields{
				db: db,
			},
			args: args{
				trackId: "ABC",
			},
			want: []Event{
				{
					Description: "abc 1",
				},
				{
					Description: "abc 2",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSql(tt.fields.db, log.NewNopLogger())
			got, err := s.GetEvents(tt.args.trackId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !assert.Equal(t, len(tt.want), len(got)) {
				return
			}

			for i, event := range got {
				assert.Equal(t, tt.want[i].Description, event.Description)
			}
		})
	}
}

func Test_sqlStorage_Remove(t *testing.T) {
	s := NewSql(createDb(t), log.NewNopLogger())

	assert.NoError(t, s.AddTrack(1, "123", "test"))
	assert.NoError(t, s.AddTrack(2, "456", "test"))

	assert.Len(t, s.GetAllTracks(), 2)

	assert.NoError(t, s.Remove("123"))
	tracks := s.GetAllTracks()
	assert.Len(t, tracks, 1)
	assert.Equal(t, "456", tracks[0].Number)

	assert.NoError(t, s.Remove("456"))
	assert.Len(t, s.GetAllTracks(), 0)

	s = NewSql(createBrokenDb(t), log.NewNopLogger())

	assert.Error(t, s.Remove("456"))
}

func Test_sqlStorage_GetTrackForUser(t *testing.T) {
	s := NewSql(createDb(t), log.NewNopLogger())

	assert.NoError(t, s.AddTrack(1, "ABC", "test"))
	assert.NoError(t, s.AddTrack(1, "EFG", "test"))

	// Try to get first record
	track, err := s.GetTrackForUser("ABC", 1)
	assert.NoError(t, err, "cannot get tracking number")
	assert.Equal(t, "ABC", track.Number, "wrong tracking number")

	// Try to get not existed record
	_, err = s.GetTrackForUser("HJK", 1)
	assert.Equal(t, NotFound, err)

	s = NewSql(createBrokenDb(t), log.NewNopLogger())

	_, err = s.GetTrackForUser("HJK", 1)
	assert.Error(t, err, "broken db should return an error")
	assert.NotEqual(t, NotFound, err, "error should not be Not Found if we have problem with database")
}

func Test_sqlStorage_countTracks(t *testing.T) {
	store := NewSql(createDb(t), log.NewNopLogger())

	assert.NoError(t, store.AddTrack(1, "ABC", "test 1"))
	assert.NoError(t, store.AddTrack(1, "EFG", "test 2"))
	assert.NoError(t, store.AddTrack(1, "XYZ", "test 3"))
	assert.NoError(t, store.AddTrack(2, "ABC2", "test2 1"))

	brokenStore := NewSql(createBrokenDb(t), log.NewNopLogger())

	tests := []struct {
		name      string
		store     *sqlStorage
		userId    int
		wantCount int
		wantErr   bool
	}{
		{
			name:      "expected 3",
			store:     store,
			userId:    1,
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:      "expected 1",
			store:     store,
			userId:    2,
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "expected 0",
			store:     store,
			userId:    3,
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "broken db",
			store:     brokenStore,
			userId:    1,
			wantCount: 0,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCount, err := tt.store.countTracks(tt.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("countTracks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCount != tt.wantCount {
				t.Errorf("countTracks() gotCount = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}

func Test_sqlStorage_IsLimitExceeded(t *testing.T) {
	db := createDb(t)
	store := NewSql(db, log.NewNopLogger())

	customLimit := 2
	userWithCustomLimit := User{
		ID:         4,
		TrackLimit: customLimit,
	}
	db.Create(&userWithCustomLimit)

	assert.NoError(t, store.AddTrack(1, "A", "test"))
	assert.NoError(t, store.AddTrack(1, "B", "test"))
	assert.NoError(t, store.AddTrack(1, "C", "test"))
	assert.NoError(t, store.AddTrack(1, "E", "test"))
	assert.NoError(t, store.AddTrack(1, "F", "test"))

	assert.NoError(t, store.AddTrack(2, "A", "test"))

	assert.NoError(t, store.AddTrack(4, "A", "test"))
	assert.NoError(t, store.AddTrack(4, "B", "test"))

	tests := []struct {
		name    string
		store   *sqlStorage
		userId  int
		want    bool
		wantErr bool
	}{
		{
			name:    "ok",
			store:   store,
			userId:  2,
			want:    false,
			wantErr: false,
		},
		{
			name:    "exceeded",
			store:   store,
			userId:  1,
			want:    true,
			wantErr: false,
		},
		{
			name:    "custom",
			store:   store,
			userId:  4,
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.store.IsLimitExceeded(tt.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsLimitExceeded() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsLimitExceeded() got = %v, want %v", got, tt.want)
			}
		})
	}
}
