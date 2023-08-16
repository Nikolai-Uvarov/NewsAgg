package main

import (
	"NewsAgg/pkg/db/dbmock"
	"NewsAgg/pkg/db/obj"
	"NewsAgg/pkg/rss"
	"fmt"
	"time"
)

func main() {
	c := rss.Listen("https://habr.com/ru/rss/best/daily/?fl=ru", time.Minute)

	db := dbmock.New()

	var stop chan int

	go func() {
		for {
			var a string
			fmt.Scanln(&a)
			if a == "" {
				stop <- 1
			}
		}
	}()

	var dbp obj.Post

loop:
	for i := 0; i < 10; i++ {
		select {
		case p := <-c:
			dbp = obj.RssToObjConvert(p)
			db.SavePost(dbp)
		case <-stop:
			break loop
		}
	}

	posts := db.GetTopPosts(10)

	for _, p := range posts {
		fmt.Println(p.ID, p.Title, p.PubTime)
	}

}
