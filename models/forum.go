package models

type Forum struct {
	Posts   int64  `json:"posts, int"`
	Slug    string `json:"slug, string"`
	Threads int32  `json:"threads, int"`
	Title   string `json:"title, string"`
	User    string `json:"user,string"`
}
