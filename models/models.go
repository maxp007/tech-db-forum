package models

type Error struct {
	Message string `json:"message"`
}

type ForumRequest struct {
	Posts   int64  `json:"posts,omitempty"`
	Slug    string `json:"slug"`
	Threads int32  `json:"threads,omitempty"`
	Title   string `json:"title"`
	User    string `json:"user"`
}

type ForumResponse struct {
	Posts   int64  `json:"posts"`
	Slug    string `json:"slug"`
	Threads int    `json:"threads"`
	Title   string `json:"title"`
	User    string `json:"user"`
}

type In_Post struct {
	Author  string `json:"author"`
	Message string `json:"message"`
	Parent  int    `json:"parent"`
}

type Post struct {
	Author   string `json:"author"`
	Created  string `json:"created"`
	Forum    string `json:"forum"`
	Id       int    `json:"id"`
	IsEdited string `json:"is_edited"`
	Message  string `json:"message"`
	Parent   int    `json:"parent"`
	Thread   int    `json:"thread"`
}

type User struct {
	About    string `json:"about"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname,omitempty"`
}

type PostFull struct {
	Author User
	Forum  ForumResponse
	Post   Thread
}

type PostUpdate struct {
	description string
	message     string
}

type Status struct {
	Forum  int32 `json:"forum"`
	Post   int64 `json:"post"`
	Thread int32 `json:"thread"`
	User   int32 `json:"user"`
}

type Thread struct {
	Author  string `json:"author"`
	Created string `json:"created"`
	Forum   string `json:"forum"`
	Id      int32  `json:"id"`
	Message string `json:"message"`
	Slug    string `json:"slug,omitempty"`
	Title   string `json:"title"`
	Votes   int32  `json:"-"`
}

type ThreadFull struct {
	Author  string `json:"author"`
	Created string `json:"created"`
	Forum   string `json:"forum"`
	Id      int32  `json:"id"`
	Message string `json:"message"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Votes   int32  `json:"votes"`
}

type ThreadNoSlug struct {
	Author  string `json:"author"`
	Created string `json:"created"`
	Forum   string `json:"forum"`
	Id      int32  `json:"id"`
	Message string `json:"message"`
	Slug    string `json:"-"`
	Title   string `json:"title"`
	Votes   int32  `json:"-"`
}
