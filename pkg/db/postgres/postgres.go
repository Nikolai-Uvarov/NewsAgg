// Пакет для взаимоействия с бд postgres
package postgres

import (
	"NewsAgg/pkg/db/obj"
	"context"
	"log"
	"math"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	DB  *pgxpool.Pool
	ctx context.Context
}

func init() {

}

// создает новое подключение к БД
func New() *DB {
	db := new(DB)
	db.ctx = context.Background()
	// Подключение к БД
	dbpass := os.Getenv("dbpass")
	dbhost := os.Getenv("dbhost")
	var err error
	db.DB, err = pgxpool.Connect(db.ctx, "postgres://postgres:"+dbpass+"@"+dbhost)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	return db
}

// Сохраняет пост, представленный объектом obj.Post, в БД
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

// Возвращает n новейших по дате создания постов из БД
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

// Возвращает пост из БД по его ID
func (db *DB) GetPostByID(id int) (obj.Post, error) {
	rows, err := db.DB.Query(db.ctx, `SELECT * FROM news WHERE id=($1);`, id)

	if err != nil {
		return obj.Post{}, err
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
			return obj.Post{}, err
		}

		posts = append(posts, post)
	}
	// проверить rows.Err()
	return posts[0], rows.Err()
}

func (db *DB) SearchPost(str string, p int) ([]obj.Post, *obj.Pagination, error) {

	str = "%" + str + "%"

	//считаем количество страниц
	rows, err := db.DB.Query(db.ctx,
		`SELECT count(*) from news 
		WHERE title ILIKE ($1);`,
		str)

	if err != nil {
		return nil, nil, err
	}

	var pag = obj.Pagination{
		Page: p,
		Of:   0,
		PostsPerPage:  obj.PostsPerPage,
	}

	for rows.Next() {
		var count int

		err = rows.Scan(
			&count)

		if err != nil {
			return nil, nil, err
		}

		pag.Of = count/obj.PostsPerPage + int(math.Pow(float64(count%obj.PostsPerPage),0))
	}

	rows, err = db.DB.Query(db.ctx,
		`SELECT * from news 
		WHERE title ILIKE ($1)  
		ORDER BY pubtime DESC 
		OFFSET ($2) LIMIT ($3);`,
		str, (p-1)*obj.PostsPerPage, obj.PostsPerPage)

	if err != nil {
		return nil, nil, err
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
			return nil, nil, err
		}

		posts = append(posts, post)
	}
	// проверить rows.Err()
	return posts, &pag, rows.Err()
}
