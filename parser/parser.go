package parser

import "time"

type Event struct {
	When        time.Time
	Description []string
}

type Parser interface {
	Parse(track string) ([]Event, error)
}
