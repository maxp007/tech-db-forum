package models

type Status struct {
	Forum  int32 `json:"forum, int"`
	Post   int64 `json:"post, int"`
	Thread int32 `json:"thread, int"`
	User   int32 `json:"user, int"`
}
