package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/maxp007/tech-db-forum/database"
	"github.com/maxp007/tech-db-forum/models"
	"net/http"
	"strings"
)

func PostForumCreate(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("PostForumCreate")

	w.Header().Add("Content-type", "application/json")

	var ForumRequest models.ForumRequest

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&ForumRequest)
	if err != nil {
		fmt.Println("PostForumCreate, decode err,", err)
	}
	forum, code := database.GetInstance().MethodCreateOrGetForum(&ForumRequest)
	if code == 409 {
		fmt.Println(`PostForumCreate, MethodCreateOrGetForum error`, err)
		w.WriteHeader(http.StatusConflict)

		bytes, err := json.Marshal(forum)
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}
		return
	}
	if code == 404 {
		fmt.Println(`PostForumCreate, MethodCreateOrGetForum error`, err)
		w.WriteHeader(http.StatusNotFound)
		response := models.Error{fmt.Sprintln("Can't find user with nickname ", ForumRequest.User)}
		bytes, err := json.Marshal(response)
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}
		return
	} else {
		w.WriteHeader(http.StatusCreated)

		bytes, err := json.Marshal(forum)
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}

		return
	}

}

func GetForumDetails(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("GetForumDetails")
	w.Header().Add("Content-type", "application/json")

	vars := mux.Vars(r)
	slug, found := vars["slug"]
	if !found {
		fmt.Print("Didn't find `slug`.", slug)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	forum, code := database.GetInstance().MethodGetForumDetails(slug)
	if code != 404 {
		w.WriteHeader(200)
		bytes, err := json.Marshal(&forum)
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}
	} else {
		w.WriteHeader(404)
		bytes, err := json.Marshal(&models.Error{fmt.Sprintf("cant find forum with slug", slug)})
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}
	}

}

func PostThreadCreate(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("PostThreadCreate")
	w.Header().Add("Content-type", "application/json")

	vars := mux.Vars(r)
	forum_slug, found := vars["slug"]
	if !found {
		fmt.Print("Didn't find `slug`.", forum_slug)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var ThreadModel models.Thread

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&ThreadModel)
	if err != nil {
		fmt.Println("PostThreadCreate, decode err,", err)
	}
	thread_slug := ThreadModel.Slug
	ThreadModel.Forum = forum_slug
	thread, code := database.GetInstance().MethodCreateOrGetThread(&ThreadModel)
	if code == 201 {
		if thread_slug == "" {
			thread.Slug = ""
		}
		w.WriteHeader(201)
		bytes, err := json.Marshal(thread)
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}
		return
	}
	if code == 404 {
		w.WriteHeader(404)
		bytes, err := json.Marshal(&models.Error{fmt.Sprintf("cant find forum with slug", forum_slug)})
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}
		return
	}
	if code == 409 {
		w.WriteHeader(409)
		bytes, err := json.Marshal(thread)
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}
		return
	}

}

func GetForumThreads(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("GetForumThreads")
	w.Header().Add("Content-type", "application/json")

	vars := mux.Vars(r)
	forum_slug, found := vars["slug"]
	if !found {
		fmt.Print("Didn't find `slug`.", forum_slug)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	param_limit := r.URL.Query().Get("limit")
	if param_limit == "" {
		param_limit = "10000000"
	}

	param_since := r.URL.Query().Get("since")
	if param_since == "" {
		param_since = "1970-01-01T12:00:00.000Z"
	}

	param_desc := r.URL.Query().Get("desc")
	if param_desc == "" {
		param_desc = "false"
	}

	threads, code := database.GetInstance().MethodGetThreads(forum_slug, param_limit, param_since, param_desc)
	if code != 404 {
		w.WriteHeader(200)
		if len(threads) == 0 {

			_, err := w.Write([]byte("[]"))
			if err != nil {
				fmt.Printf("failed to write response ")
			}
		} else {
			bytes, err := json.Marshal(&threads)
			if err != nil {
				fmt.Printf("failed to unmarshal response ")
			}
			_, err = w.Write(bytes)
			if err != nil {
				fmt.Printf("failed to write response ")
			}
		}

	} else {
		w.WriteHeader(404)
		bytes, err := json.Marshal(&models.Error{fmt.Sprintf("cant find forum with slug", forum_slug)})
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}
	}
}

func GetPostDetails(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("GetPostDetails")
	w.Header().Add("Content-type", "application/json")
	vars := mux.Vars(r)
	post_id, found := vars["id"]
	if !found {
		fmt.Print("Didn't find `id`.", post_id)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := r.ParseForm()
	if err != nil {
		fmt.Println("GetPostDetails ParseForm error", err)
	}
	a := r.Form["related"]
	var s []string
	if a != nil {
		s = strings.Split(a[0], ",")
		for i, _ := range s {
			fmt.Println(s[i])
		}
	} else {
		s = make([]string, 0)
	}

	post_full, code := database.GetInstance().MethodGetPostDetails(post_id, s)
	if code != 404 {
		w.WriteHeader(200)

		bytes, err := json.Marshal(&post_full)
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}

	} else {
		w.WriteHeader(404)
		bytes, err := json.Marshal(&models.Error{fmt.Sprintf("cant find post with id", post_id)})
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}
	}

}

func PostPostUpdate(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("PostPostUpdate")
	w.Header().Add("Content-type", "application/json")
	vars := mux.Vars(r)
	post_id, found := vars["id"]
	if !found {
		fmt.Print("Didn't find `id`.", post_id)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var newPostMessage models.PostMessageUpdate

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&newPostMessage)
	if err != nil {
		fmt.Println("PostThreadCreate, decode err,", err)
	}
	post_full, code := database.GetInstance().MethodPostUpdate(post_id, newPostMessage.Message)
	if code != 404 {
		w.WriteHeader(200)

		bytes, err := json.Marshal(&post_full)
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}

	} else {
		w.WriteHeader(404)
		bytes, err := json.Marshal(&models.Error{fmt.Sprintf("cant find post with id", post_id)})
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}
	}
}

func GetForumUsers(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("GetForumUsers")
	w.Header().Add("Content-type", "application/json")

	vars := mux.Vars(r)
	forum_slug, found := vars["slug"]
	if !found {
		fmt.Print("Didn't find `slug`.", forum_slug)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	param_limit := r.URL.Query().Get("limit")
	if param_limit == "" {
		param_limit = "0"
	}

	param_since := r.URL.Query().Get("since")
	if param_since == "" {
		param_since = ""
	}

	param_desc := r.URL.Query().Get("desc")
	if param_desc == "" {
		param_desc = "false"
	}

	users, code := database.GetInstance().MethodGetForumUsers(forum_slug, param_limit, param_since, param_desc)
	if code != 404 {
		w.WriteHeader(200)
		if len(users) == 0 {
			_, err := w.Write([]byte("[]"))
			if err != nil {
				fmt.Printf("failed to write response ")
			}
		} else {
			bytes, err := json.Marshal(&users)
			if err != nil {
				fmt.Printf("failed to unmarshal response ")
			}
			_, err = w.Write(bytes)
			if err != nil {
				fmt.Printf("failed to write response ")
			}
		}

	} else {
		w.WriteHeader(404)
		bytes, err := json.Marshal(&models.Error{fmt.Sprintf("cant find forum with slug", forum_slug)})
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}
	}
}

//--------------------------SERVICE-------------------------------
func PostServiceClear(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PostServiceClear")
	w.Header().Add("Content-type", "application/json")
	code := database.GetInstance().ServiceCleanData()
	if code != 200 {
		w.WriteHeader(500)
		_, err := w.Write([]byte("[]"))
		if err != nil {
			fmt.Printf("failed to write response ")
		}
		fmt.Println("PostServiceClear ERROR")
		return
	} else {
		_, err := w.Write([]byte("[]"))
		if err != nil {
			fmt.Printf("failed to write response ")
		}
		w.WriteHeader(200)
		return
	}

}

func GetServiceStatus(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("GetServiceStatus")
	w.Header().Add("Content-type", "application/json")

	status, code := database.GetInstance().MethodGetServiceStatus()
	if code == 200 {
		w.WriteHeader(200)

		bytes, err := json.Marshal(&status)
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}

	}
	if code == 404 {
		w.WriteHeader(404)
		bytes, err := json.Marshal(&models.Error{fmt.Sprintf("cant find tread ")})
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}
	}
}

//--------------------------END SERVICE-------------------------------

func PostPostCreate(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("PostPostCreate")
	w.Header().Add("Content-type", "application/json")

	vars := mux.Vars(r)
	thread_slug_or_id, found := vars["slug_or_id"]
	if !found || thread_slug_or_id == "" {
		fmt.Print("Didn't find `slug_or_id`.", thread_slug_or_id)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var PostsSlice []models.In_Post

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&PostsSlice)
	if err != nil {
		fmt.Println("PostThreadCreate, decode err,", err)
	}

	postsSlice, code := database.GetInstance().MethodCreatePost(PostsSlice, thread_slug_or_id)

	if code == 201 {
		w.WriteHeader(201)
		if len(postsSlice) == 0 {
			_, err = w.Write([]byte("[]"))
			if err != nil {
				fmt.Printf("failed to write response ")
			}
		} else {

			bytes, err := json.Marshal(&postsSlice)
			if err != nil {
				fmt.Printf("failed to unmarshal response ")
			}
			_, err = w.Write(bytes)
			if err != nil {
				fmt.Printf("failed to write response ")
			}
		}

	}
	if code == 404 {
		w.WriteHeader(404)
		bytes, err := json.Marshal(&models.Error{fmt.Sprintf("cant find tread ")})
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}
	}
	if code == 409 {
		w.WriteHeader(409)
		bytes, err := json.Marshal(&models.Error{fmt.Sprintf("Conflict 409 ")})
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}
	}
}

func PostThreadVote(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("PostThreadVote")
	w.Header().Add("Content-type", "application/json")

	vars := mux.Vars(r)
	thread_slug_or_id, found := vars["slug_or_id"]

	if !found || thread_slug_or_id == "" {
		fmt.Print("PostThreadVote, Didn't find `slug_or_id`.", thread_slug_or_id)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var VoteModel models.Vote

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&VoteModel)
	if err != nil {
		fmt.Println("PostThreadCreate, decode err,", err)
	}

	voted_thread, code := database.GetInstance().MethodVote(&VoteModel, thread_slug_or_id)
	if code == 200 {
		w.WriteHeader(200)

		bytes, err := json.Marshal(&voted_thread)
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}
	}
	if code == 404 {
		w.WriteHeader(404)
		bytes, err := json.Marshal(&voted_thread)
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}
	}
}

func GetThreadDetails(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("GetThreadDetails")
	w.Header().Add("Content-type", "application/json")

	vars := mux.Vars(r)
	thread_slug_or_id, found := vars["slug_or_id"]

	if !found || thread_slug_or_id == "" {
		fmt.Print("GetThreadDetails, Didn't find `slug_or_id`.", thread_slug_or_id)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	thread, code := database.GetInstance().MethodGetDetails(thread_slug_or_id)
	if code == 200 {
		w.WriteHeader(200)

		bytes, err := json.Marshal(&thread)
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}

	} else {
		w.WriteHeader(404)
		bytes, err := json.Marshal(&models.Error{fmt.Sprintf("cant find thread_slug_or_id", thread_slug_or_id)})
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}
	}

}

func GetThreadPosts(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("GetThreadDetails")
	w.Header().Add("Content-type", "application/json")

	vars := mux.Vars(r)
	thread_slug_or_id, found := vars["slug_or_id"]

	if !found || thread_slug_or_id == "" {
		fmt.Print("GetThreadDetails, Didn't find `slug_or_id`.", thread_slug_or_id)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	param_limit := r.URL.Query().Get("limit")
	if param_limit == "" {
		param_limit = "0"
	}

	param_since := r.URL.Query().Get("since")
	if param_since == "" {
		param_since = "0"
	}

	param_sort := r.URL.Query().Get("sort")

	if param_sort != "flat" &&
		param_sort != "tree" &&
		param_sort != "parent_tree" {
		param_sort = "flat"
	}

	param_desc := r.URL.Query().Get("desc")
	if param_desc == "" {
		param_desc = "false"
	}

	var params_struct models.ThreadDetailsParams
	params_struct.Desc = param_desc
	params_struct.Limit = param_limit
	params_struct.Since = param_since
	params_struct.Sort = param_sort

	posts, code := database.GetInstance().MethodGetThreadPosts(thread_slug_or_id, params_struct)
	if code == 200 {
		w.WriteHeader(200)
		if len(posts) == 0 {
			_, err := w.Write([]byte("[]"))
			if err != nil {
				fmt.Printf("failed to write response ")
			}
		} else {
			bytes, err := json.Marshal(&posts)
			if err != nil {
				fmt.Printf("failed to unmarshal response ")
			}
			_, err = w.Write(bytes)
			if err != nil {
				fmt.Printf("failed to write response ")
			}
		}

	} else {
		w.WriteHeader(404)
		bytes, err := json.Marshal(&models.Error{fmt.Sprintf("cant find thread_slug_or_id", thread_slug_or_id)})
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}
	}

}

func PostThreadUpdate(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("GetThreadPosts")
	w.Header().Add("Content-type", "application/json")

	vars := mux.Vars(r)
	thread_slug_or_id, found := vars["slug_or_id"]

	if !found || thread_slug_or_id == "" {
		fmt.Print("GetThreadDetails, Didn't find `slug_or_id`.", thread_slug_or_id)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var post_update_struct models.PostThreadUpdate

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&post_update_struct)
	if err != nil {
		fmt.Println("PostThreadCreate, decode err,", err)
	}
	new_thread_details, code := database.GetInstance().MethodUpdateThreadDetails(thread_slug_or_id, post_update_struct)
	if code == 200 {
		w.WriteHeader(200)

		bytes, err := json.Marshal(&new_thread_details)
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}

	} else {
		w.WriteHeader(404)
		bytes, err := json.Marshal(&models.Error{fmt.Sprintf("cant find thread_slug_or_id", thread_slug_or_id)})
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}
	}

}

//--------------------------CREATE USER-------------------------------

func PostUserCreate(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("PostUserCreate")
	w.Header().Add("Content-type", "application/json")

	vars := mux.Vars(r)
	path_nickname, found := vars["nickname"]
	if !found {
		fmt.Print("Didn't find `nickname`.")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var NewUserData models.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&NewUserData)
	NewUserData.Nickname = path_nickname

	if err != nil {
		fmt.Println("PostForumCreate, decode err,", err)
	}

	users, violationflag := database.GetInstance().MethodCreateOrGetUser(&NewUserData)

	if violationflag == 1 {
		w.WriteHeader(http.StatusConflict)
		fmt.Println("Пользователь уже присутсвует в базе данных.")

		bytes, err := json.Marshal(&users)
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}

	} else {
		w.WriteHeader(http.StatusCreated)
		for _, user := range users {
			bytes, err := json.Marshal(user)
			if err != nil {
				fmt.Printf("failed to unmarshal response ")
			}
			_, err = w.Write(bytes)
			if err != nil {
				fmt.Printf("failed to write response ")
			}
		}

	}

}

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetUserProfile")
	w.Header().Add("Content-type", "application/json")

	vars := mux.Vars(r)
	path_nickname, found := vars["nickname"]
	if !found {
		fmt.Print("Didn't find `nickname`.")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	user, violationflag := database.GetInstance().MethodGetUserProfile(path_nickname)
	if violationflag == 1 {
		w.WriteHeader(404)
		error_model := models.Error{fmt.Sprintf("Can't find user with nickname %s", path_nickname)}
		bytes, err := json.Marshal(&error_model)
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}

	} else {
		w.WriteHeader(200)
		bytes, err := json.Marshal(&user)
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}
	}
}

func PostUserUpdate(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("PostUserUpdate")
	w.Header().Add("Content-type", "application/json")

	vars := mux.Vars(r)
	path_nickname, found := vars["nickname"]
	if !found {
		fmt.Print("Didn't find `nickname`.")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var UpdatedUserData models.User

	decoder := json.NewDecoder(r.Body)
	decerr := decoder.Decode(&UpdatedUserData)
	UpdatedUserData.Nickname = path_nickname

	if decerr != nil {
		fmt.Println("PostForumCreate, decode err,", decerr)
		UpdatedUserData.Nickname = path_nickname
		UpdatedUserData.Fullname = ""
		UpdatedUserData.About = ""
		UpdatedUserData.Email = ""
	}

	user, violationflag := database.GetInstance().MethodUpdateUserProfile(&UpdatedUserData)

	if violationflag == 1 {
		w.WriteHeader(404)
		error_model := models.Error{fmt.Sprintf("Can't find user with nickname %s", path_nickname)}
		bytes, err := json.Marshal(&error_model)
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}
	} else if violationflag == 2 {
		w.WriteHeader(409)
		error_model := models.Error{fmt.Sprintf("Can't find user with nickname %s", path_nickname)}
		bytes, err := json.Marshal(&error_model)
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}
	} else {
		w.WriteHeader(200)
		bytes, err := json.Marshal(&user)
		if err != nil {
			fmt.Printf("failed to unmarshal response ")
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("failed to write response ")
		}
	}

}

//--------------------------END CREATE USER-------------------------------
