package main

import (
	"fmt"
	"github.com/fr05t1k/pechkin/parser"
	"github.com/fr05t1k/pechkin/storage"
	"github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
	"time"
)

const startMessage = `Hello! I can help you with tracking your packages.`

const availableCommands = `
Here is a list of available commands:
/add <tracking number> <name> - add tracking number
/list - show all your tracking numbers
/history <tracking number> - show all events for given tracking number
`

type Handler struct {
	logger logrus.FieldLogger
	store  storage.Storage
	bot    *tb.Bot
}

func NewHandler(log logrus.FieldLogger, store storage.Storage, bot *tb.Bot) *Handler {
	return &Handler{
		logger: log,
		store:  store,
		bot:    bot,
	}
}

func (h *Handler) StartHandler(m *tb.Message) {
	_, _ = h.bot.Send(m.Sender, fmt.Sprintf("%s\n%s", startMessage, availableCommands))
}

func (h *Handler) HelpHandler(m *tb.Message) {
	_, _ = h.bot.Send(m.Sender, availableCommands)
}

func (h *Handler) AddHandler(m *tb.Message) {
	params := strings.Split(m.Payload, " ")
	if len(params) != 2 {
		_, _ = h.bot.Send(m.Sender, "Please provide track number and name. Example: /add RB12312412CY Watch")
		return
	}
	err := h.store.AddTrack(m.Sender.ID, params[0], params[1])
	if err != nil {
		_, _ = h.bot.Send(m.Sender, "Cannot add this tracking number")
		return
	}

	_, _ = h.bot.Send(m.Sender, fmt.Sprintf("%s Added", m.Payload))
	h.RunUpdate(h.bot, m.Payload, m.Sender, h.store)
}

func (h *Handler) ListHandler(m *tb.Message) {
	tracks := h.store.GetTracks(m.Sender.ID)
	if len(tracks) == 0 {
		_, _ = h.bot.Send(m.Sender, "You dont have tracking numbers")
		return
	}
	var trackIds []string
	for _, track := range tracks {
		trackIds = append(trackIds, track.Number+" "+track.Name)
	}
	_, _ = h.bot.Send(m.Sender, fmt.Sprintf("Here is your tracking numbers:\n%s", strings.Join(trackIds, "\n")))
}

func (h *Handler) HistoryHandler(m *tb.Message) {
	events, err := h.store.GetEvents(m.Payload)
	if err != nil {
		_, _ = h.bot.Send(m.Sender, "No history for this tracking number")
	}
	_, _ = h.bot.Send(m.Sender, ToMessage(m.Payload, events))
}

func (h *Handler) RunUpdate(b *tb.Bot, track string, sender *tb.User, store storage.Storage) {
	p := parser.NewCyprusPost()
	go func() {
		for {
			<-time.After(10 * time.Second)
			events, err := p.Parse(track)
			if err != nil {
				h.logger.WithFields(logrus.Fields{"track": track, "err": err}).Error("error parsing")
			}

			existedEvents, err := store.GetEvents(track)

			if err != nil {
				h.logger.WithFields(logrus.Fields{"track": track, "err": err}).Error("cannot get events")
			} else {
				if len(events) != len(existedEvents) {
					_, _ = b.Send(sender, ToMessage(track, events))
				}
			}

			err = store.SetHistory(track, events)
			if err != nil {
				h.logger.WithFields(logrus.Fields{"track": track, "err": err}).Error("error settings history")
			}

			h.logger.WithFields(logrus.Fields{"trackId": track}).Info("track updated")
		}
	}()
}

func ToMessage(track string, events []storage.Event) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Here is your history for tracking number %s\n", track))
	for _, event := range events {
		builder.WriteString("______________________\n")
		builder.WriteString("At ")
		builder.WriteString(event.EventAt.String())
		builder.WriteString("\n")
		builder.WriteString(event.Description)
	}

	return builder.String()
}
