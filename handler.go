package main

import (
	"fmt"
	"github.com/fr05t1k/pechkin/parser"
	"github.com/fr05t1k/pechkin/storage"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"strings"
	"time"
)

func MakeAddHandler(b *tb.Bot, store storage.Storage) func(m *tb.Message) {
	return func(m *tb.Message) {
		err := store.AddTrack(m.Sender.ID, m.Payload)
		if err != nil {
			_, _ = b.Send(m.Sender, "Cannot add this tracking number")
			return
		}

		_, _ = b.Send(m.Sender, fmt.Sprintf("%s Added", m.Payload))
		RunUpdate(b, m.Payload, *m.Sender, store)
	}
}
func MakeListHandler(b *tb.Bot, store storage.Storage) func(m *tb.Message) {
	return func(m *tb.Message) {
		tracks := store.GetTracks(m.Sender.ID)
		if len(tracks) == 0 {
			_, _ = b.Send(m.Sender, "You dont have tracking numbers")
			return
		}
		var trackIds []string
		for _, track := range tracks {
			trackIds = append(trackIds, track.Id)
		}
		_, _ = b.Send(m.Sender, fmt.Sprintf("Here is your tracking numbers:\n%s", strings.Join(trackIds, "\n")))
	}
}

func MakeHistoryHandler(b *tb.Bot, store storage.Storage) func(m *tb.Message) {
	return func(m *tb.Message) {
		events, err := store.GetEvents(m.Sender.ID, m.Payload)
		if err != nil {
			_, _ = b.Send(m.Sender, "No history for this tracking number")
		}
		_, _ = b.Send(m.Sender, ToMessage(m.Payload, events))
	}

}

func RunUpdate(b *tb.Bot, track string, sender tb.User, store storage.Storage) {
	p := parser.NewCyprusPost()
	go func() {
		for {
			<-time.After(10 * time.Second)
			events, err := p.Parse(track)
			if err != nil {
				log.Println(err)
			}

			existedEvents, err := store.GetEvents(sender.ID, track)

			if err != nil {
				log.Println(err)
			} else {
				if len(events) != len(existedEvents) {
					_, _ = b.Send(&sender, ToMessage(track, events))
				}
			}

			err = store.SetHistory(track, events)
			if err != nil {
				log.Println(err)
			}
			log.Println(track, "updated")

		}
	}()
}

func ToMessage(track string, events []storage.Event) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Here is your history for tracking number %s\n", track))
	for _, event := range events {
		builder.WriteString("______________________\n")
		builder.WriteString("At ")
		builder.WriteString(event.When.String())
		builder.WriteString("\n")
		builder.WriteString(strings.Join(event.Description, "\n"))
	}

	return builder.String()
}
