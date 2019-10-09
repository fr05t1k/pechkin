package main

import (
	"github.com/fr05t1k/pechkin/storage"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"os"
	"time"
)

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:   os.Getenv("PECHKIN_BOT_TOKEN"),
		Updates: 0,
		Poller:  &tb.LongPoller{Timeout: 10 * time.Second},
		Reporter: func(e error) {
			log.Println(e)
		},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	store := createStore()
	for _, track := range store.GetAllTracks() {
		user := tb.User{
			ID: track.ID,
		}
		RunUpdate(b, track.Number, &user, store)
	}
	b.Handle("/add", MakeAddHandler(b, store))
	b.Handle("/list", MakeListHandler(b, store))
	b.Handle("/history", MakeHistoryHandler(b, store))
	b.Start()
}

func createStore() storage.Storage {
	if os.Getenv("DATABASE_URL") != "" {
		db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			log.Fatalln(err)
		}
		db.AutoMigrate(&storage.Track{}, &storage.Event{})

		return storage.NewSql(db)
	}

	return storage.NewMemory()
}
