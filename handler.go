package main

import (
	"errors"
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
/remove <tracking number> - remove tracking number
/list - show all your tracking numbers
/history <tracking number> - show all events for given tracking number
`

const cannotDeleteTrack = "I cannot delete your tracking number right now. Please try again later."
const cannotAddTrack = "I cannot add your tracking number right now. Please try again later."
const limitExceeded = "Your limit is exceeded. Please remove one tracking number or contact @spavlovichev to increase the limit."
const cannotFindTrack = "I cannot find you tracking number. Did you add it? Check /list first."
const trackHasBeenDeleted = "Your tracking number has been deleted. Thanks for cleaning :)"
const noTracks = "You dont have tracking numbers"
const noHistory = "No history for this tracking number"

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
	isExceeded, err := h.store.IsLimitExceeded(m.Sender.ID)
	if err != nil {
		_, _ = h.bot.Send(m.Sender, cannotAddTrack)
		return
	}
	if isExceeded {
		_, _ = h.bot.Send(m.Sender, limitExceeded)
		return
	}
	err = h.store.AddTrack(m.Sender.ID, params[0], params[1])
	if err != nil {
		_, _ = h.bot.Send(m.Sender, cannotAddTrack)
		return
	}

	_, _ = h.bot.Send(m.Sender, fmt.Sprintf("%s Added", m.Payload))
}

func (h *Handler) RemoveHandler(m *tb.Message) {
	payload := m.Payload
	if payload == "" {
		_, _ = h.bot.Send(m.Sender, "Please specify you track number. Example /remove RB12345678CY.")
		return
	}
	track, err := h.store.GetTrackForUser(payload, m.Sender.ID)
	if err == storage.NotFound {
		_, _ = h.bot.Send(m.Sender, cannotFindTrack)
		return
	}
	if err != nil {
		_, _ = h.bot.Send(m.Sender, cannotDeleteTrack)
		return
	}
	if track.UserId != m.Sender.ID {
		_, _ = h.bot.Send(m.Sender, cannotFindTrack)
		return
	}
	err = h.store.Remove(payload)
	if err != nil {
		_, _ = h.bot.Send(m.Sender, cannotDeleteTrack)
		return
	}

	_, _ = h.bot.Send(m.Sender, trackHasBeenDeleted)
}

func (h *Handler) ListHandler(m *tb.Message) {
	tracks := h.store.GetTracks(m.Sender.ID)
	if len(tracks) == 0 {
		_, _ = h.bot.Send(m.Sender, noTracks)
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
		_, _ = h.bot.Send(m.Sender, noHistory)
	}
	_, _ = h.bot.Send(m.Sender, MakeHistoryMessage(m.Payload, events))
}

func RunUpdate(b *tb.Bot, track string, store storage.Storage) error {
	p := parser.NewCyprusPost()
	events, err := p.Parse(track)
	if err != nil {
		return errors.New("error parsing")
	}

	existedEvents, err := store.GetEvents(track)
	if err != nil {
		return errors.New("cannot get events")
	}

	logrus.Printf("%d %d %s", len(events), len(existedEvents), track)
	if len(events) <= len(existedEvents) {
		return nil
	}

	err = store.SetHistory(track, events)
	if err != nil {
		return errors.New("error settings history")
	}

	tracks, err := store.GetTrackByNumber(track)
	if err != nil {
		return err
	}

	for i := range tracks {
		_, _ = b.Send(&tb.User{ID: tracks[i].UserId}, MakeNewUpdateMessage(tracks[i], events[len(existedEvents):]))
	}

	return nil
}

func runUpdates(b *tb.Bot, store storage.Storage, logger logrus.FieldLogger, each time.Duration) {
	go func() {
		for {
			<-time.After(each)
			for _, track := range store.GetAllTracks() {
				if err := RunUpdate(b, track.Number, store); err != nil {
					logger.WithFields(logrus.Fields{"trackId": track.Number, "err": err}).Info("cannot update track")
				} else {
					logger.WithFields(logrus.Fields{"trackId": track.Number}).Info("track updated")
				}

			}
		}
	}()
}

func MakeNewUpdateMessage(track storage.Track, events []storage.Event) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("You have new updates for %s (%s) 🎉\n", track.Name, track.Number))
	builder.WriteString(MakeEventsMessage(events))

	return builder.String()
}

func MakeHistoryMessage(track string, events []storage.Event) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Here is your history for tracking number %s\n", track))
	builder.WriteString(MakeEventsMessage(events))

	return builder.String()
}

func MakeEventsMessage(events []storage.Event) string {
	builder := strings.Builder{}
	for _, event := range events {
		builder.WriteString("\n")
		builder.WriteString("⏱️ ")
		builder.WriteString(event.EventAt.String())
		builder.WriteString("\n")
		builder.WriteString(event.Description)
	}

	return builder.String()
}
