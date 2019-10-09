package storage

import (
	"reflect"
	"testing"
	"time"
)

func Test_memory_AddTrack(t *testing.T) {
	type fields struct {
		tracks  map[int][]Track
		history map[Track][]Event
	}
	type args struct {
		userId  int
		trackId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "add",
			fields:  fields{},
			args:    args{1, "TEST"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMemory()
			if err := m.AddTrack(tt.args.userId, tt.args.trackId); (err != nil) != tt.wantErr {
				t.Errorf("AddTrack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_memory_GetTracks(t *testing.T) {
	type fields struct {
		tracks  map[int][]Track
		history map[string][]Event
	}
	type args struct {
		userId int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantTracks []Track
	}{
		{
			name:   "not found",
			fields: fields{},
			args: args{
				userId: 1,
			},
			wantTracks: nil,
		},
		{
			name: "exists",
			fields: fields{
				tracks:  map[int][]Track{1: {{Number: "abc"}}},
				history: nil,
			},
			args: args{
				userId: 1,
			},
			wantTracks: []Track{{Number: "abc"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &memory{
				tracks:  tt.fields.tracks,
				history: tt.fields.history,
			}
			if gotTracks := m.GetTracks(tt.args.userId); !reflect.DeepEqual(gotTracks, tt.wantTracks) {
				t.Errorf("GetTracks() = %v, want %v", gotTracks, tt.wantTracks)
			}
		})
	}
}

func Test_memory_SetHistory(t *testing.T) {
	type fields struct {
		tracks  map[int][]Track
		history map[string][]Event
	}
	type args struct {
		trackId string
		events  []Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "simple set",
			fields: fields{},
			args: args{
				trackId: "abc",
				events: []Event{
					{
						EventAt:     time.Now(),
						Description: "test",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMemory()
			if err := m.SetHistory(tt.args.trackId, tt.args.events); (err != nil) != tt.wantErr {
				t.Errorf("SetHistory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
