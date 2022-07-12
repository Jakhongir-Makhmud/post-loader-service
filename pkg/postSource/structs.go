package postSource


type DraftPost struct {
	Meta Meta   `json:"meta"`
	Data []Data `json:"data"`
}

type Pagination struct {
	Total int   `json:"total"`
	Pages int   `json:"pages"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
type Meta struct {
	Pagination Pagination `json:"pagination"`
}

type Data struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}