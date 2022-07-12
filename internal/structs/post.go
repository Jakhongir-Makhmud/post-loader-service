package structs

type Post struct {
	Id     int    `json:"id" db:"post_id"`
	Title  string `json:"title" db:"title"`
	Body   string `json:"body" db:"body"`
}

