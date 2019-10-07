package main

import (
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

	b.Handle("/add", MakeAddHandler(b))
	b.Handle("/list", MakeListHandler(b))
	b.Handle("/history", MakeHistoryHandler(b))
	b.Start()
}
