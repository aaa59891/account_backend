package pagination

type Pagination struct {
	Current  int
	Next     int
	Previous int
	Total    int
}

type Page struct {
	CurrentPage int `form:"page"`
}
