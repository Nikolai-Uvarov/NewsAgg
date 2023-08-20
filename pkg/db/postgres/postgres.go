// Пакет для взаимоействия с бд postgres
package postgres

import (
	"NewsAgg/pkg/db/obj"
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	DB  *pgxpool.Pool
	ctx context.Context
}

func init() {

}

func New() *DB {
	db := new(DB)
	db.ctx = context.Background()
	// Подключение к БД
	dbpass := os.Getenv("dbpass")
	var err error
	db.DB, err = pgxpool.Connect(db.ctx, "postgres://postgres:"+dbpass+"@192.168.1.35:5432/news")

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	return db
}

func (db *DB) SavePost(p obj.Post) error {

	_, err := db.DB.Exec(db.ctx,
		`INSERT INTO news (title, content, pubtime, link) 
		VALUES (($1), ($2), ($3), ($4))`,
		p.Title, p.Content, p.PubTime, p.Link)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetTopPosts(n int) ([]obj.Post, error) {
	rows, err := db.DB.Query(db.ctx, `SELECT * FROM news ORDER BY pubtime DESC LIMIT ($1);`, n)

	if err != nil {
		return nil, err
	}

	var posts []obj.Post

	for rows.Next() {
		var post obj.Post

		err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.PubTime,
			&post.Link)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}
	// проверить rows.Err()
	return posts, rows.Err()
}
