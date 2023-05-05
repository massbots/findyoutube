package main

import (
	"log"
	"os"

	"go.massbots.xyz/findyoutube/youtube"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/layout"
)

func main() {
	lt, err := layout.New("bot.yml", templateFuncs)
	if err != nil {
		log.Fatal(err)
	}

	b, err := tele.NewBot(lt.Settings())
	if err != nil {
		log.Fatal(err)
	}

	yt, err := youtube.NewClient(os.Getenv("YT_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	b.Handle("/start", func(c tele.Context) error {
		return c.Send(lt.TextLocale("ru", "start"))
	})

	b.Handle(tele.OnQuery, func(c tele.Context) error {
		search, err := yt.Search(c.Text(), 10)
		if err != nil {
			return err
		}

		var results tele.Results
		for _, v := range search {
			switch v.Id.Kind {
			case "youtube#video":
				r := lt.ResultLocale("ru", "search_video", v)
				results = append(results, r)
			case "youtube#channel":
				r := lt.ResultLocale("ru", "search_channel", v)
				results = append(tele.Results{r}, results...)
			}
		}

		return nil
	})

	b.Start()
}
