package models

type Post struct {
	Description string `json:"description,int"`
	Posts       int64  `json:"posts, int"`
	Slug        string `json:"slug, string"`
	Threads     int64  `json:"threads, int"`
	Title       string `json:"title, string"`
	User        string `json:"user, string"`
}
