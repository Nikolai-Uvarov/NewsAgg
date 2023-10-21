// objects to interact with db data-model
package obj

// тип obj.Post для операций с БД
// от obj.Post отличается типом ID - int, присваеваемый в нашей БД
// вместо string из источника RSS
type Post struct {
	ID      int    // номер записи
	Title   string // заголовок публикации
	Content string // содержание публикации
	PubTime int64  // время публикации
	Link    string // ссылка на источник
}

type DB interface {
	SavePost(Post) error
	GetTopPosts(int) ([]Post, error)
	GetPostByID(int) (Post, error)
	SearchPost(string, int) ([]Post, *Pagination, error)
}

type Pagination struct {
	Page int  //текущая страница
	Of int    //всего страниц
	PostsPerPage int   //PostsPerPage
}

//параметр пагинации - количество постов на страницу
const PostsPerPage = 15

//структура для ответа на запрос
type Answer struct{
	Posts []Post
	Pagination Pagination
}
