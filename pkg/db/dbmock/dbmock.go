package dbmock

import (
	"NewsAgg/pkg/db/obj"
	"sort"
)

type DB struct {
	posts  []obj.Post
	nextid int
}

func New() *DB {
	db := new(DB)
	db.nextid = 1
	return db
}

func (db *DB) SavePost(p obj.Post) {
	p.ID = db.nextid
	db.posts = append(db.posts, p)
	db.nextid++
}

func (db *DB) GetTopPosts(n int) []obj.Post {
	sort.Slice(db.posts, func(i, j int) bool { return db.posts[i].PubTime > db.posts[j].PubTime })
	return db.posts[:n]
}
