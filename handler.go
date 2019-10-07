package main

import (
	"fmt"
	"github.com/fr05t1k/pechkin/parser"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"strings"
	"time"
)

var tracks map[int][]string
var history map[string][]parser.Event

func MakeAddHandler(b *tb.Bot) func(m *tb.Message) {
	tracks = make(map[int][]string)
	return func(m *tb.Message) {
		tracks[m.Sender.ID] = append(tracks[m.Sender.ID], m.Payload)
		_, _ = b.Send(m.Sender, fmt.Sprintf("%s Added", m.Payload))
		RunUpdate(b, m.Payload, m.Sender)
	}
}
func MakeListHandler(b *tb.Bot) func(m *tb.Message) {
	return func(m *tb.Message) {
		_, _ = b.Send(m.Sender, fmt.Sprintf("Here is your tracking numbers:\n%s", strings.Join(tracks[m.Sender.ID], "\n")))
	}
}

func MakeHistoryHandler(b *tb.Bot) func(m *tb.Message) {
	return func(m *tb.Message) {
		events, ok := history[m.Payload]
		if !ok {
			_, _ = b.Send(m.Sender, "No history for this tracking number")
		} else {
			_, _ = b.Send(m.Sender, ToMessage(m.Payload, events))
		}
	}

}

func RunUpdate(b *tb.Bot, track string, sender tb.Recipient) {
	p := parser.NewCyprusPost()
	if history == nil {
		history = make(map[string][]parser.Event)
	}

	go func() {
		for {
			<-time.After(10 * time.Second)
			events, err := p.Parse(track)
			if err != nil {
				log.Println(err)
			}

			if len(events) != len(history[track]) {
				_, _ = b.Send(sender, ToMessage(track, events))
			}

			history[track] = events
			log.Println(track, "updated")

		}
	}()
}

func ToMessage(track string, events []parser.Event) string {
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
