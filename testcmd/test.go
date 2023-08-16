package main

import (
	"NewsAgg/pkg/rss"
	"fmt"
	"time"
)

func main() {
	c := rss.Listen("https://habr.com/ru/rss/best/daily/?fl=ru", time.Minute)

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

loop:
	for {
		select {
		case p := <-c:
			fmt.Println(p.ID, p.Title)
		case <-stop:
			break loop
		}
	}

}
