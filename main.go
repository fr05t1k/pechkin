package main

import (
	"github.com/fr05t1k/pechkin/storage"
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

	store := storage.NewMemory()
	b.Handle("/add", MakeAddHandler(b, store))
	b.Handle("/list", MakeListHandler(b, store))
	b.Handle("/history", MakeHistoryHandler(b, store))
	b.Start()
}
