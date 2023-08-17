//Authored by Nikolai Uvarov
//Allowed for free copying
//Package rss implements a simple mechanism for listening of certain rss news feed

package rss

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"NewsAgg/pkg/db/obj"

	"github.com/SlyMarbo/rss"
)

var config = "config.json"

// для переопределения минимального интервала обновления rss ленты
// при попытке обновить rss ленту с меньшим интервалом, чем DefaultRefreshInterval, SlyMarbo/rss генерит ошибку
var minRefreshInterval = time.Second * 30

func init() {
	rss.DefaultRefreshInterval = minRefreshInterval
}

type Post struct {
	ID      string // номер записи
	Title   string // заголовок публикации
	Content string // содержание публикации
	PubTime int64  // время публикации
	Link    string // ссылка на источник
}

func Listen(url string, period time.Duration) <-chan Post {
	c := make(chan Post)

	go func() {
		defer close(c)

		f, err := rss.Fetch(url)
		if err != nil {
			log.Printf("Error to fetch rss %v : %v", url, err)
		}

		var posts []Post

		//первый раз отправляем в канал все посты из feed
		for k, i := range f.Items {
			posts = append(posts, itemToPost(i))
			c <- posts[k]
		}

		for {
			time.Sleep(period)

			//запросить обновление RSS
			err := f.Update()
			if err != nil {
				log.Printf("Error to update rss %v : %v", url, err)
			}

			//отправляем в канал только новые посты из feed
			for _, i := range f.Items {

				post := itemToPost(i)

				isnew := true

				for _, p := range posts {
					if p.ID == post.ID {
						isnew = false
					}
				}

				if isnew {
					posts = append(posts, post)
					c <- post
				}

			}
		}
	}()

	return c
}

// Формирует объект Post, забирая из Item нужные поля
func itemToPost(i *rss.Item) (p Post) {
	p.ID = i.ID
	p.Content = i.Summary
	p.Link = i.Link
	p.PubTime = i.Date.Unix()
	p.Title = i.Title
	return p
}

// Конвертирует тип пакета rss в тип для взаимодействия с БД сервиса
func RssToObjConvert(p Post) obj.Post {
	return obj.Post{
		ID:      0,
		Title:   p.Title,
		Content: p.Content,
		PubTime: p.PubTime,
		Link:    p.Link,
	}
}

// RSSMultiplex возвращает общий канал, в который будут попадать сообщения от всех
// источников
func RSSMultiplex(channels ...<-chan Post) <-chan Post {
	var wg sync.WaitGroup

	multiplexedChan := make(chan Post)
	multiplex := func(c <-chan Post) {
		defer wg.Done()
		for i := range c {
			multiplexedChan <- i
		}
	}
	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}
	// Горутина, которая закроет канал
	go func() {
		wg.Wait()
		close(multiplexedChan)
	}()
	return multiplexedChan
}

// обработчик rss-канала, сохраняющий посты в базу
func Collect(db obj.DB) {

	ttu, urls := readConfig()

	var cs []<-chan Post

	//запуск прослушивания rss лент
	for _, u := range urls {
		cs = append(cs, Listen(u, time.Minute*time.Duration(ttu)))
	}

	//мультиплексирование всех прослушиваемых лент в 1 канал
	c := RSSMultiplex(cs...)

	//сохранение постов в БД
	go func() {
		for p := range c {
			dbp := RssToObjConvert(p)
			db.SavePost(dbp)
		}
	}()

}

// возвращает массив строк с адресами rss лент и частоту обновления из config
func readConfig() (int, []string) {
	f, err := os.OpenFile(config, os.O_RDONLY, 0777)
	if err != nil {
		log.Fatal("Cannot read config file: ", err)
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		log.Fatal("Cannot read config file: ", err)
	}

	conf := struct {
		Rss            []string
		Request_period int
	}{}

	err = json.Unmarshal(buf, &conf)
	if err != nil {
		log.Fatal("Cannot parse config: ", err)
	}

	return conf.Request_period, conf.Rss
}
