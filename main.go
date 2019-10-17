package main

import (
	"github.com/fr05t1k/pechkin/storage"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"os"
	"time"
)

func main() {

	logger := logrus.New()

	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	b, err := tb.NewBot(tb.Settings{
		Token:   os.Getenv("PECHKIN_BOT_TOKEN"),
		Updates: 0,
		Poller:  &tb.LongPoller{Timeout: 10 * time.Second},
		Reporter: func(e error) {
			logger.Error(e)
		},
	})

	if err != nil {
		logger.Fatal(err)
		return
	}

	store := createStore(logger)
	handler := NewHandler(logger, store, b)

	runUpdates(b, store, logger, 10*time.Second)

	b.Handle("/add", handler.AddHandler)
	b.Handle("/remove", handler.RemoveHandler)
	b.Handle("/list", handler.ListHandler)
	b.Handle("/history", handler.HistoryHandler)
	b.Handle("/start", handler.StartHandler)
	b.Handle("/help", handler.HelpHandler)
	b.Start()
}

func createStore(logger logrus.FieldLogger) storage.Storage {
	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.WithField("err", err).Error("cannot to connect to the storage")
	}
	db.AutoMigrate(
		&storage.Track{},
		&storage.Event{},
		&storage.User{},
	)

	return storage.NewSql(db, logger)
}
