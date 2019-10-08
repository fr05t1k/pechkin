package parser

import (
	"github.com/fr05t1k/pechkin/storage"
	"time"
)

type Event struct {
	When        time.Time
	Description []string
}

type Parser interface {
	Parse(track string) ([]storage.Event, error)
}
