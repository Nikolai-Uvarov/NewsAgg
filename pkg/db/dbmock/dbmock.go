// Пакет in-mem database для разработки и тестов
package dbmock

import (
	"NewsAgg/pkg/db/obj"
	"fmt"
	"sort"
)

type DB struct {
	posts  []obj.Post
	nextid int
}

// создает новое подключение к БД
func New() *DB {
	db := new(DB)
	db.nextid = 1
	return db
}

// Сохраняет пост, представленный объектом obj.Post, в БД
func (db *DB) SavePost(p obj.Post) error {
	p.ID = db.nextid
	db.posts = append(db.posts, p)
	db.nextid++
	return nil
}

// Возвращает n новейших по дате создания постов из БД
func (db *DB) GetTopPosts(n int) ([]obj.Post, error) {
	sort.Slice(db.posts, func(i, j int) bool { return db.posts[i].PubTime > db.posts[j].PubTime })
	return db.posts[:n], nil
}

// Возвращает количество постов, сохраненных в  момент вызова в БД (для тестов)
func (db *DB) Len() int {
	return len(db.posts)
}

// Возвращает пост из БД по его ID
func (db *DB) GetPostByID(id int) (obj.Post, error) {

	if id > len(db.posts){
		return obj.Post{}, fmt.Errorf("id not found: id = %v", id)
	}
	return db.posts[id-1], nil
}
