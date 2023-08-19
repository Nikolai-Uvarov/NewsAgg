package main

import (
	"NewsAgg/pkg/db/dbmock"
	//"NewsAgg/pkg/db/obj"
	"NewsAgg/pkg/rss"
	"fmt"
	"time"
)

func main() {

	// a := rss.Listen("http://www.kommersant.ru/RSS/news.xml", time.Minute)
	// b := rss.Listen("https://habr.com/ru/rss/best/daily/?fl=ru", time.Minute)
	// d := rss.Listen("https://tass.ru/rss/v2.xml", time.Minute)
	// e := rss.Listen("http://www.polit.ru/rss/index.xml", time.Minute)

	// //c := rss.RSSMultiplex(b, a, e, d)

	db := dbmock.New()

	rss.Collect(db)

	time.Sleep(30 * time.Second)

	posts := db.GetTopPosts(10)

	for _, p := range posts {
		fmt.Println(p.ID, p.Title, p.PubTime, p.Link)
	}

}
