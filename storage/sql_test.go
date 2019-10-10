package storage

import (
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
	db.AutoMigrate(&Track{}, &Event{})

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
			s := NewSql(tt.db)
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
	s := NewSql(createDb(t))

	err := s.AddTrack(1, "ABC")
	assert.NoError(t, err)

	err = s.AddTrack(1, "EFG")
	assert.NoError(t, err)

	tracks := s.GetAllTracks()

	// check length
	if !assert.Len(t, tracks, 2) {
		return
	}
	// check necessary fields
	assert.Equal(t, 1, tracks[0].UserId)
	assert.Equal(t, "ABC", tracks[0].Number)
	assert.Equal(t, 1, tracks[1].UserId)
	assert.Equal(t, "EFG", tracks[1].Number)

	s = NewSql(createBrokenDb(t))

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
			s := NewSql(tt.db)
			assert.NoError(t, s.AddTrack(1, "ABC"))
			assert.NoError(t, s.AddTrack(2, "EFG"))
			assert.NoError(t, s.AddTrack(2, "XYZ"))

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
	s := NewSql(db)
	assert.NoError(t, s.AddTrack(1, "ABC"))
	assert.NoError(t, s.AddTrack(2, "EFG"))
	assert.NoError(t, s.AddTrack(2, "XYZ"))

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
			s := NewSql(tt.db)

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
	s := NewSql(db)
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
			s := NewSql(tt.fields.db)
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
