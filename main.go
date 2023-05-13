package main

import (
	"log"
	"os"
	"strings"

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

	b.Use(lt.Middleware("ru"))

	b.Handle("/start", func(c tele.Context) error {
		return c.Send(lt.Text(c, "start", c.Sender()))
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		text := c.Text()

		if strings.Count(text, "\n") > 1 ||
			strings.Contains(text, "https://") {
			return nil
		}

		return c.Send(
			lt.Text(c, "query", text),
			lt.Markup(c, "query", text),
		)
	})

	b.Handle(tele.OnQuery, func(c tele.Context) error {
		search, err := yt.Search(c.Data(), 10)
		if err != nil {
			return err
		}

		var results tele.Results
		for _, v := range search {
			switch v.Id.Kind {
			case "youtube#video":
				r := lt.Result(c, "search_video", v)
				results = append(results, r)
			case "youtube#channel":
				r := lt.Result(c, "search_channel", v)
				results = append(tele.Results{r}, results...)
			}
		}

		return c.Answer(&tele.QueryResponse{Results: results})
	})

	b.Start()
}
