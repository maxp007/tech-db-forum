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
	IsEdited string `json:"-"`
	Message  string `json:"message"`
	Parent   int    `json:"parent,omitempty"`
	Thread   int    `json:"thread"`
}

type PostWithEdited struct {
	Author   string `json:"author"`
	Created  string `json:"created"`
	Forum    string `json:"forum"`
	Id       int    `json:"id"`
	IsEdited bool   `json:"isEdited"`
	Message  string `json:"message"`
	Parent   int    `json:"parent,omitempty"`
	Thread   int    `json:"thread"`
}

type User struct {
	About    string `json:"about"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname,omitempty"`
}

type PostThreadUpdate struct {
	Message string `json:"message,omitempty"`
	Title   string `json:"title,omitempty"`
}

type Status struct {
	User   *int32 `json:"user"`
	Forum  *int32 `json:"forum"`
	Post   *int64 `json:"post"`
	Thread *int32 `json:"thread"`
}

type Thread struct {
	Author  string `json:"author"`
	Created string `json:"created"`
	Forum   string `json:"forum"`
	Id      int32  `json:"id,omitempty"`
	Message string `json:"message"`
	Slug    string `json:"slug,omitempty"`
	Title   string `json:"title"`
	Votes   int32  `json:"votes"`
}

type ThreadFull struct {
	Author  string `json:"author"`
	Created string `json:"created"`
	Forum   string `json:"forum"`
	Id      int32  `json:"id,omitempty"`
	Message string `json:"message"`
	Slug    string `json:"slug,omitempty"`
	Title   string `json:"title"`
	Votes   int32  `json:"votes,omitempty"`
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

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
}

type ThreadDetailsParams struct {
	Limit string
	Since string
	Sort  string
	Desc  string
}

type ThreadDetails struct {
	Message string `json:"message"`
	Title   string `json:"title"`
}

type PostFull struct {
	Author User           `json:"user,omitempty"`
	Forum  ForumResponse  `json:"forum,omitempty"`
	Post   PostWithEdited `json:"post,omitempty"`
	Thread Thread         `json:"thread,omitempty"`
}

type PostFullOnlyPost struct {
	Author User           `json:"-"`
	Forum  ForumResponse  `json:"-"`
	Post   PostWithEdited `json:"post,omitempty"`
	Thread Thread         `json:"-"`
}

type PostFullOnlyPostAndUser struct {
	Author *User           `json:"author,omitempty"`
	Forum  *ForumResponse  `json:"forum,omitempty"`
	Post   *PostWithEdited `json:"post,omitempty"`
	Thread *Thread         `json:"thread,omitempty"`
}

type PostFullOnlyPostAndThread struct {
	Author User           `json:"-"`
	Forum  ForumResponse  `json:"-"`
	Post   PostWithEdited `json:"post,omitempty"`
	Thread Thread         `json:"thread"`
}

type PostMessageUpdate struct {
	Message string `json:"message"`
}
