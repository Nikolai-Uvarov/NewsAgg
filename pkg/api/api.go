package api

import (
	"NewsAgg/pkg/db/obj"
	"encoding/json"

	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// API приложения.
type API struct {
	r  *mux.Router // маршрутизатор запросов
	db obj.DB      // база данных
}

// Конструктор API.
func New(db obj.DB) *API {
	api := API{}
	api.db = db
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

// Router возвращает маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.r
}

// Регистрация методов API в маршрутизаторе запросов.
func (api *API) endpoints() {
	// получить n последних новостей
	api.r.HandleFunc("/news/{n}", api.posts).Methods(http.MethodGet, http.MethodOptions)
	// получить новость по postID, строке search, в том числе с указанием страницы
	api.r.HandleFunc("/news", api.postWithFilters).Methods(http.MethodGet, http.MethodOptions)
	// веб-приложение
	api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
	//заголовок ответа
	api.r.Use(api.HeadersMiddleware)
}

// HeadersMiddleware устанавливает заголовки ответа сервера.
func (api *API) HeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// posts возвращает n новейших новостей в зависимости от параметра пути n
func (api *API) posts(w http.ResponseWriter, r *http.Request) {
	// Считывание параметра {n} из пути запроса.
	// Например, /news/10.
	s := mux.Vars(r)["n"]
	n, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Получение данных из БД.
	posts, err := api.db.GetTopPosts(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Отправка данных клиенту в формате JSON.
	json.NewEncoder(w).Encode(posts)
	// Отправка клиенту статуса успешного выполнения запроса
	w.WriteHeader(http.StatusOK)
}

// postByID возвращает пост по postID
// при указании параметра search возвращает посты по вхождению строки в заголовке
// page - номер возвращаемой страницы 
func (api *API) postWithFilters(w http.ResponseWriter, r *http.Request) {
	
	// Считывание параметра  строки запроса.
	idParam := r.URL.Query().Get("postID")
	if idParam!=""{
		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Получение данных из БД.
		post, err := api.db.GetPostByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Отправка данных клиенту в формате JSON.
		json.NewEncoder(w).Encode(post)
		// Отправка клиенту статуса успешного выполнения запроса
		w.WriteHeader(http.StatusOK)
		return
	}

	// Если нет page - в переменной будет пустая строка
	pageParam := r.URL.Query().Get("page")

	var page int
	var err error

	if pageParam != "" {
		page, err = strconv.Atoi(pageParam)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		page = 1
	}

	// Считывание параметра  строки запроса.
	str := r.URL.Query().Get("search")

	// Получение данных из БД.
	posts, err := api.db.SearchPost(str,page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Отправка данных клиенту в формате JSON.
	json.NewEncoder(w).Encode(posts)
	// Отправка клиенту статуса успешного выполнения запроса
	w.WriteHeader(http.StatusOK)
}

