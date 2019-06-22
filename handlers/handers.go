package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/maxp007/tech-db-forum/database"
	"github.com/maxp007/tech-db-forum/models"
	"net/http"
)

func PostForumCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PostForumCreate")

	w.Header().Add("Content-type", "application/json")

	var ForumRequest models.ForumRequest

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&ForumRequest)
	if err != nil {
		fmt.Println("PostForumCreate, decode err,", err)
	}
	forum, code := database.MethodCreateOrGetForum(&ForumRequest)
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
	fmt.Println("GetForumDetails")
	w.Header().Add("Content-type", "application/json")

	vars := mux.Vars(r)
	slug, found := vars["slug"]
	if !found {
		fmt.Print("Didn't find `slug`.", slug)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	forum, code := database.MethodGetForumDetails(slug)
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
	fmt.Println("PostThreadCreate")
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
	thread, code := database.MethodCreateOrGetThread(&ThreadModel)
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
	fmt.Println("PostThreadCreate")
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

	threads, code := database.MethodGetThreads(forum_slug, param_limit, param_since, param_desc)
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

func GetForumUsers(w http.ResponseWriter, r *http.Request) {

}

func GetPostDetails(w http.ResponseWriter, r *http.Request) {

}

func PostPostUpdate(w http.ResponseWriter, r *http.Request) {

}

//--------------------------SERVICE-------------------------------
func PostServiceClear(w http.ResponseWriter, r *http.Request) {

}

func GetServiceStatus(w http.ResponseWriter, r *http.Request) {

}

//--------------------------END SERVICE-------------------------------

func PostPostCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PostThreadCreate")
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

	postsSlice, code := database.MethodCreatePost(PostsSlice, thread_slug_or_id)

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
		bytes, err := json.Marshal(&postsSlice)
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

func GetThreadInfo(w http.ResponseWriter, r *http.Request) {

}

func PostThreadUpdate(w http.ResponseWriter, r *http.Request) {

}

func GetThreadPosts(w http.ResponseWriter, r *http.Request) {

}

func PostThreadVote(w http.ResponseWriter, r *http.Request) {

}

//--------------------------CREATE USER-------------------------------

func PostUserCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PostUserCreate")
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

	users, violationflag := database.MethodCreateOrGetUser(&NewUserData)

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
	user, violationflag := database.MethodGetUserProfile(path_nickname)
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
	fmt.Println("PostUserUpdate")
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

	user, violationflag := database.MethodUpdateUserProfile(&UpdatedUserData)

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
