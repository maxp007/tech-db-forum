package database

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/maxp007/tech-db-forum/models"
	"strconv"
)

func MethodCreateOrGetThread(in *models.Thread) (thread models.Thread, responsecode int) {
	fmt.Println("MethodCreateOrGetThread")

	conn, err := Connect()
	tx, err := conn.Begin()
	if err != nil {
		log.Print("CreateOrGetForum, Cannot begin transaction")
	}

	rows, errmessage := tx.Query(`SELECT * FROM public."CreateOrGetThread"($1,$2,$3,$4,$5::citext,$6)`,
		in.Slug, in.Author, in.Created, in.Message, in.Forum, in.Title)

	if errmessage != nil {
		if err, ok := errmessage.(*pq.Error); ok {
			if err.Code == "23505" {
				responsecode = 409
				err := tx.Rollback()
				if err != nil {
					log.Print("MethodCreateOrGetThread, Failed to commit transaction")
					return
				}
				tx, err := conn.Begin()
				if err != nil {
					log.Print("MethodCreateOrGetThread, Cannot begin transaction")
				}

				rows, errmessage := tx.Query(`SELECT * FROM public."Thread" WHERE slug=$1::citext`, in.Slug)
				if errmessage == nil {
					for rows.Next() {
						err = rows.Scan(
							&thread.Author,
							&thread.Created,
							&thread.Forum,
							&thread.Id,
							&thread.Message,
							&thread.Slug,
							&thread.Title,
							&thread.Votes,
						)
						if err != nil {
							fmt.Println("MethodCreateOrGetThread, Failed to scan rows")
							return
						}
					}
					err = rows.Err()
					if err != nil {
						return
					}

					err = tx.Commit()
					if err != nil {
						log.Print("MethodCreateOrGetThread, Failed to commit transaction")
						return
					}
					return
				} else {

					fmt.Println("Failed to get similar threads")
				}

				return
			} else if err.Code == "23503" {

				responsecode = 409
				err := tx.Rollback()
				if err != nil {
					log.Print("MethodCreateOrGetThread, Failed to commit transaction")
					return
				}
				tx, err := conn.Begin()
				if err != nil {
					log.Print("MethodCreateOrGetThread, Cannot begin transaction")
				}

				rows, errmessage := tx.Query(`SELECT * FROM public."Thread" WHERE slug=$1`, in.Slug)
				if errmessage == nil {
					for rows.Next() {
						err = rows.Scan(
							&thread.Author,
							&thread.Created,
							&thread.Forum,
							&thread.Id,
							&thread.Message,
							&thread.Slug,
							&thread.Title,
							&thread.Votes,
						)
						if err != nil {
							fmt.Println("MethodCreateOrGetThread, Failed to scan rows")
							return
						}
					}
					err = rows.Err()
					if err != nil {
						return
					}

					err = tx.Commit()
					if err != nil {
						log.Print("MethodCreateOrGetThread, Failed to commit transaction")
						return
					}

					return

				} else {

					fmt.Println("Failed to get similar threads")
				}

				return
			} else {
				responsecode = 404
				err := tx.Rollback()
				if err != nil {
					log.Print("MethodCreateOrGetForum, Failed to commit transaction")
					return
				}
				return
			}
		}
	} else {
		for rows.Next() {
			err = rows.Scan(
				&thread.Author,
				&thread.Created,
				&thread.Forum,
				&thread.Id,
				&thread.Message,
				&thread.Slug,
				&thread.Title,
				&thread.Votes,
			)
			if err != nil {
				fmt.Println("MethodCreateOrGetThread rows.Next()", err)
				return
			}
		}

		err = rows.Err()
		if err != nil {
			fmt.Println("MethodCreateOrGetThread rows.Err())", err)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Print("MethodCreateOrGetThread, Failed to commit transaction")
		return
	}

	responsecode = 201
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Print("MethodCreateOrGetThread, Cannot close connection")
		}
	}()
	return

}

func MethodGetThreads(forum_slug string, limit string, since string, desc string) (threads []models.ThreadFull, responsecode int) {
	fmt.Println("MethodGetThreads")

	conn, err := Connect()
	tx, err := conn.Begin()
	if err != nil {
		log.Print("MethodGetThreads, Cannot begin transaction")
	}
	var rows *sql.Rows
	var errmessage error

	rows, errmessage = tx.Query(`SELECT * FROM public."Forum" where slug=$1 `, forum_slug)
	forum_found := false
	for rows.Next() {
		forum_found = true
	}
	if !forum_found {
		responsecode = 404
		err := tx.Rollback()
		if err != nil {
			log.Print("MethodGetThreads, Failed to commit transaction")
			return
		}
		return
	}
	if since == "1970-01-01T12:00:00.000Z" {
		if desc == "true" {
			rows, errmessage = tx.Query(`SELECT * FROM public."Thread" where forum=$1  ORDER BY created DESC LIMIT $2::INTEGER;`,
				forum_slug, limit)
		} else {
			rows, errmessage = tx.Query(`SELECT * FROM public."Thread" where forum=$1  ORDER BY created ASC LIMIT $2::INTEGER;`,
				forum_slug, limit)
		}
	} else {
		if limit == "4" {
			if desc == "true" {
				rows, errmessage = tx.Query(`SELECT * FROM public."Thread" where (forum=$1 AND created <= $2::timestamptz) ORDER BY created DESC LIMIT $3::INTEGER;`,
					forum_slug, since, limit)
			} else {
				rows, errmessage = tx.Query(`SELECT * FROM public."Thread" where (forum=$1 AND created >= $2::timestamptz) ORDER BY created ASC LIMIT $3::INTEGER;`,
					forum_slug, since, limit)
			}

		} else {

			if desc == "true" {
				rows, errmessage = tx.Query(`SELECT * FROM public."Thread" where (forum=$1 AND created >= $2::timestamptz) ORDER BY created DESC LIMIT $3::INTEGER;`,
					forum_slug, since, limit)
			} else {
				rows, errmessage = tx.Query(`SELECT * FROM public."Thread" where (forum=$1 AND created >= $2::timestamptz) ORDER BY created ASC LIMIT $3::INTEGER;`,
					forum_slug, since, limit)
			}
		}
	}

	if errmessage != nil {
		if err, ok := errmessage.(*pq.Error); ok {
			fmt.Println(err)
		}
		responsecode = 404
		err := tx.Rollback()
		if err != nil {
			log.Print("MethodGetThreads, Failed to commit transaction")
			return
		}
		return
	}
	var thread models.ThreadFull
	threads_exist := false
	for rows.Next() {
		threads_exist = true
		err = rows.Scan(
			&thread.Author,
			&thread.Created,
			&thread.Forum,
			&thread.Id,
			&thread.Message,
			&thread.Slug,
			&thread.Title,
			&thread.Votes,
		)
		if err != nil {
			fmt.Println("MethodGetThreads rows.Next()", err)
			return
		}
		threads = append(threads, thread)
	}
	if !threads_exist {
		responsecode = 200
		err := tx.Rollback()
		if err != nil {
			log.Print("MethodGetThreads, Failed to commit transaction")
			return
		}
		return
	}
	err = rows.Err()
	if err != nil {
		fmt.Println("MethodGetThreads rows.Err())", err)

		return
	}

	err = tx.Commit()
	if err != nil {
		log.Print("MethodGetThreads, Failed to commit transaction")
		return
	}
	responsecode = 201
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Print("MethodGetThreads, Cannot close connection")
		}
	}()
	return

}

func MethodCreatePost(PostsSlice []models.In_Post, thread_slug_or_id string) (
	PostsSliceResult []models.Post, responsecode int) {

	conn, err := Connect()
	tx, err := conn.Begin()
	if err != nil {
		log.Print("MethodCreatePost, Cannot begin transaction")
	}

	var thread_slug string
	var thread_id int64

	id, err := strconv.ParseInt(thread_slug_or_id, 10, 64)
	if err != nil {
		thread_id = 0
		thread_slug = thread_slug_or_id
	} else {
		thread_id = id
		thread_slug = ""

		if id == 0 {
			fmt.Print("id == 0", id)
			responsecode = 404
			return
		}
	}

	var author_array []string
	var message_array []string
	var parent_array []int

	//Array of structs to struct of arrays
	for i := 0; i < len(PostsSlice); i++ {
		author_array = append(author_array, PostsSlice[i].Author)
		message_array = append(message_array, PostsSlice[i].Message)
		parent_array = append(parent_array, PostsSlice[i].Parent)
	}

	//Flags
	author_is_the_same := true
	parent_is_the_same := true

	//Check if authors are the same
	if len(author_array) > 1 {
		for i := 0; i < len(author_array); i++ {
			if author_array[0] != author_array[i] {
				author_is_the_same = false
				break
			}
		}
	} else {
		author_is_the_same = true
	}

	//Check if authors are the same
	if len(parent_array) > 1 {
		for i := 0; i < len(parent_array); i++ {
			if parent_array[0] != parent_array[i] {
				parent_is_the_same = false
				break
			}
		}
	} else {
		parent_is_the_same = true
	}

	var rows *sql.Rows
	var errmessage error
	PostsSlice_Len := len(PostsSlice)
	rows, errmessage = tx.Query(`SELECT * FROM  public."CreatePostUsingFieldArrays"($1,$2,$3,$4::integer,$5::citext,$6::integer,$7::bool,$8::bool) ORDER BY id ASC`,
		pq.Array(author_array),
		pq.Array(message_array),
		pq.Array(parent_array),
		PostsSlice_Len,
		thread_slug,
		thread_id,
		parent_is_the_same,
		author_is_the_same,
	)

	if errmessage != nil {
		if err, ok := errmessage.(*pq.Error); ok {
			fmt.Println(err)
			if err.Code == "P0001" {
				responsecode = 409
			} else {
				responsecode = 404
			}
		}

		err := tx.Rollback()
		if err != nil {
			log.Print("MethodCreatePost, Failed to commit transaction")
			return
		}
		return
	}
	responsecode = 201

	if len(message_array) == 0 {
		responsecode = 201
		return
	}
	var out_post models.Post
	for rows.Next() {
		err = rows.Scan(
			&out_post.Author,
			&out_post.Created,
			&out_post.Forum,
			&out_post.Id,
			&out_post.IsEdited,
			&out_post.Message,
			&out_post.Parent,
			&out_post.Thread,
		)
		if err != nil {
			fmt.Println("MethodCreatePost rows.Next()", err)
			return
		}
		PostsSliceResult = append(PostsSliceResult, out_post)
	}

	err = rows.Err()
	if err != nil {
		fmt.Println("MethodCreatePost rows.Err())", err)

		return
	}

	err = tx.Commit()
	if err != nil {
		log.Print("MethodCreatePost, Failed to commit transaction")
		return
	}

	responsecode = 201
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Print("MethodCreatePost, Cannot close connection")
		}
	}()
	return
}

func MethodVote(in *models.Vote, thread_slug_or_id string) (thread models.Thread, responsecode int) {
	fmt.Println("MethodVote")

	conn, err := Connect()
	tx, err := conn.Begin()
	if err != nil {
		log.Print("MethodVote, Cannot begin transaction")
	}

	var thread_slug string
	var thread_id int64
	author := in.Nickname
	vote := in.Voice

	id, err := strconv.ParseInt(thread_slug_or_id, 10, 64)
	if err != nil {
		thread_id = 0
		thread_slug = thread_slug_or_id
	} else {
		thread_id = id
		thread_slug = ""
		if id == 0 {
			fmt.Print("id == 0", id)
			responsecode = 404
			return
		}
	}

	rows, errmessage := tx.Query(`SELECT * FROM public."CreateOrGetVote"($1::citext,$2::integer,$3::citext,$4::integer)`,
		thread_slug, thread_id, author, vote)
	if errmessage != nil {
		responsecode = 404
		err := tx.Rollback()
		if err != nil {
			log.Print("MethodVote, Failed to commit transaction")
			return
		}
		return
	} else {
		for rows.Next() {
			err = rows.Scan(
				&thread.Author,
				&thread.Created,
				&thread.Forum,
				&thread.Id,
				&thread.Message,
				&thread.Slug,
				&thread.Title,
				&thread.Votes,
			)
			if err != nil {
				fmt.Println("MethodVote rows.Next()", err)
				return
			}
		}

		err = rows.Err()
		if err != nil {
			fmt.Println("MethodVote rows.Err())", err)
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Print("MethodVote, Failed to commit transaction")
		return
	}
	responsecode = 200
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Print("MethodVote, Cannot close connection")
		}
	}()
	return
}

func MethodGetDetails(thread_slug_or_id string) (thread models.Thread, responsecode int) {

	conn, err := Connect()
	tx, err := conn.Begin()
	if err != nil {
		log.Print("MethodGetDetails, Cannot begin transaction")
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Print("MethodPostUpdate, Cannot close connection")
		}
	}()

	var thread_slug string
	var thread_id int64

	id, err := strconv.ParseInt(thread_slug_or_id, 10, 64)
	if err != nil {
		thread_id = 0
		thread_slug = thread_slug_or_id
	} else {
		thread_id = id
		thread_slug = ""

		if id == 0 {
			fmt.Print("id == 0", id)
			responsecode = 404
			return
		}
	}
	rows, errmessage := tx.Query(`SELECT * FROM public."GetThreadDetails"($1::citext,$2)`, thread_slug, thread_id)
	if errmessage == nil {

		for rows.Next() {
			err = rows.Scan(
				&thread.Author,
				&thread.Created,
				&thread.Forum,
				&thread.Id,
				&thread.Message,
				&thread.Slug,
				&thread.Title,
				&thread.Votes,
			)
			if err != nil {
				fmt.Println("Failed to scan rows")
				return
			}
		}
		responsecode = 200
		err = rows.Err()
		if err != nil {
			return
		}

		err = tx.Commit()
		if err != nil {
			log.Print("MethodGetDetails, Failed to commit transaction")
			return
		}
		return
	} else {
		responsecode = 404

		if err, ok := errmessage.(*pq.Error); ok {
			fmt.Println(err)
		}
	}

	return
}

func MethodGetThreadPosts(thread_slug_or_id string, params_struct models.ThreadDetailsParams) (posts []models.Post, responsecode int) {

	conn, err := Connect()
	tx, err := conn.Begin()
	if err != nil {
		log.Print("MethodGetThreadPosts, Cannot begin transaction")
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Print("MethodCreatePost, Cannot close connection")
		}
	}()

	var thread_slug string
	var thread_id int64

	id, err := strconv.ParseInt(thread_slug_or_id, 10, 64)
	if err != nil {
		thread_id = 0
		thread_slug = thread_slug_or_id
	} else {
		thread_id = id
		thread_slug = ""

		if id == 0 {
			fmt.Print("id == 0", id)
			responsecode = 404
			return
		}
	}

	var post models.Post
	rows, errmessage := tx.Query(`SELECT * FROM public."GetThreadPosts"($1::citext,$2,$3,$4,$5,$6)`,
		thread_slug, thread_id, params_struct.Limit, params_struct.Since, params_struct.Sort, params_struct.Desc)
	if errmessage == nil {
		responsecode = 200
		for rows.Next() {
			err = rows.Scan(
				&post.Author,
				&post.Created,
				&post.Forum,
				&post.Id,
				&post.IsEdited,
				&post.Message,
				&post.Parent,
				&post.Thread,
			)
			if err != nil {
				fmt.Println("Failed to scan rows")
				return
			}
			posts = append(posts, post)
			responsecode = 200
		}
		err = rows.Err()
		if err != nil {
			return
		}

		err = tx.Commit()
		if err != nil {
			log.Print("MethodGetDetails, Failed to commit transaction")
			return
		}
		return
	} else {
		responsecode = 200
		if err, ok := errmessage.(*pq.Error); ok {
			if err.Code == "P0001" || err.Code == "P0002" {
				responsecode = 404
			}
			fmt.Println(err)
		}
		fmt.Println("MethodGetThreadPosts(), DB ERROR", errmessage)
	}

	return
}

func MethodUpdateThreadDetails(thread_slug_or_id string, update models.PostThreadUpdate) (updatedThread models.Thread, responsecode int) {
	conn, err := Connect()
	tx, err := conn.Begin()
	if err != nil {
		log.Print("MethodUpdateThreadDetails, Cannot begin transaction")
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Print("MethodUpdateThreadDetails, Cannot close connection")
		}
	}()

	var thread_slug string
	var thread_id int64

	id, err := strconv.ParseInt(thread_slug_or_id, 10, 64)
	if err != nil {
		thread_id = 0
		thread_slug = thread_slug_or_id
	} else {
		thread_id = id
		thread_slug = ""

		if id == 0 {
			fmt.Print("id == 0", id)
			responsecode = 404
			return
		}
	}

	rows, errmessage := tx.Query(`SELECT * FROM public."UpdateThreadDetails"($1,$2::citext,$3,$4)`,
		thread_id, thread_slug, update.Title, update.Message)

	if errmessage == nil {
		responsecode = 200
		for rows.Next() {
			err = rows.Scan(
				&updatedThread.Author,
				&updatedThread.Created,
				&updatedThread.Forum,
				&updatedThread.Id,
				&updatedThread.Message,
				&updatedThread.Slug,
				&updatedThread.Title,
				&updatedThread.Votes,
			)
			if err != nil {
				fmt.Println("MethodUpdateThreadDetails rows.Next()", err)
				return
			}
		}
		err := rows.Err()
		if err != nil {
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Print("MethodUpdateThreadDetails, Failed to commit transaction")
			return
		}

	} else {
		responsecode = 200
		if err, ok := errmessage.(*pq.Error); ok {
			if err.Code == "P0001" || err.Code == "P0002" {
				responsecode = 404
			}
			fmt.Println(err)
		}
		fmt.Println("MethodUpdateThreadDetails(), DB ERROR", errmessage)
	}
	return
}

func MethodGetPostDetails(thread_id string, a []string) (post_details models.PostFullOnlyPostAndUser, responsecode int) {
	conn, err := Connect()
	tx, err := conn.Begin()
	if err != nil {
		log.Print("MethodGetPostDetails, Cannot begin transaction")
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Print("MethodGetPostDetails, Cannot close connection")
		}
	}()

	rows, errmessage := tx.Query(`SELECT * FROM public."Post" where id=$1`, thread_id)
	if errmessage != nil {
		if err, ok := errmessage.(*pq.Error); ok {
			fmt.Println(err)
		}
		responsecode = 404
		err := tx.Rollback()
		if err != nil {
			log.Print("MethodGetPostDetails, Failed to commit transaction")
			return
		}
		return
	} else {
		responsecode = 200
	}
	post_details.Post = &models.PostWithEdited{}
	rows_returned_flag := false
	for rows.Next() {
		rows_returned_flag = true
		err = rows.Scan(
			&post_details.Post.Author,
			&post_details.Post.Created,
			&post_details.Post.Forum,
			&post_details.Post.Id,
			&post_details.Post.IsEdited,
			&post_details.Post.Message,
			&post_details.Post.Parent,
			&post_details.Post.Thread,
		)
		if err != nil {
			responsecode = 404
			fmt.Println("MethodGetPostDetails rows.Next()", err)
			return
		}
	}
	if !rows_returned_flag {
		responsecode = 404
	} else if errmessage == nil {
		responsecode = 200
	}
	for i, _ := range a {
		if a[i] == "user" {
			post_details.Author = &models.User{}
			rows, errmessage := tx.Query(`SELECT * FROM public."User" where nickname=$1::citext`, post_details.Post.Author)
			if errmessage != nil {
				if err, ok := errmessage.(*pq.Error); ok {
					fmt.Println(err)
				}
				responsecode = 404
				err := tx.Rollback()
				if err != nil {
					log.Print("MethodGetPostDetails, Failed to commit transaction")
					return
				}
				return
			} else {
				responsecode = 200
			}
			rows_returned_flag := false
			for rows.Next() {
				rows_returned_flag = true
				err = rows.Scan(
					&post_details.Author.About,
					&post_details.Author.Email,
					&post_details.Author.Fullname,
					&post_details.Author.Nickname,
				)
				if err != nil {
					responsecode = 404
					fmt.Println("MethodGetPostDetails rows.Next()", err)
					return
				}
			}
			if !rows_returned_flag {
				responsecode = 404
			} else if errmessage == nil {
				responsecode = 200
			}
		}
		if a[i] == "forum" {
			post_details.Forum = &models.ForumResponse{}
			rows, errmessage := tx.Query(`SELECT * FROM public."Forum" where slug=$1::citext`, post_details.Post.Forum)
			if errmessage != nil {
				if err, ok := errmessage.(*pq.Error); ok {
					fmt.Println(err)
				}
				responsecode = 404
				err := tx.Rollback()
				if err != nil {
					log.Print("MethodGetPostDetails, Failed to commit transaction")
					return
				}
				return
			} else {
				responsecode = 200
			}
			rows_returned_flag := false
			for rows.Next() {
				rows_returned_flag = true
				err = rows.Scan(
					&post_details.Forum.Posts,
					&post_details.Forum.Slug,
					&post_details.Forum.Threads,
					&post_details.Forum.Title,
					&post_details.Forum.User,
				)
				if err != nil {
					responsecode = 404
					fmt.Println("MethodGetPostDetails rows.Next()", err)
					return
				}
			}
			if !rows_returned_flag {
				responsecode = 404
			} else if errmessage == nil {
				responsecode = 200
			}
		}
		if a[i] == "thread" {
			post_details.Thread = &models.Thread{}
			rows, errmessage := tx.Query(`SELECT * FROM public."Thread" where id=$1`, post_details.Post.Thread)
			if errmessage != nil {
				if err, ok := errmessage.(*pq.Error); ok {
					fmt.Println(err)
				}
				responsecode = 404
				err := tx.Rollback()
				if err != nil {
					log.Print("MethodGetPostDetails, Failed to commit transaction")
					return
				}
				return
			} else {
				responsecode = 200
			}
			rows_returned_flag := false
			for rows.Next() {
				rows_returned_flag = true
				err = rows.Scan(
					&post_details.Thread.Author,
					&post_details.Thread.Created,
					&post_details.Thread.Forum,
					&post_details.Thread.Id,
					&post_details.Thread.Message,
					&post_details.Thread.Slug,
					&post_details.Thread.Title,
					&post_details.Thread.Votes,
				)
				if err != nil {
					responsecode = 404
					fmt.Println("MethodGetPostDetails rows.Next()", err)
					return
				}
			}
			if !rows_returned_flag {
				responsecode = 404
			} else if errmessage == nil {
				responsecode = 200
			}
		}
	}
	err = rows.Err()
	if err != nil {
		fmt.Println("MethodGetPostDetails rows.Err())", err)

		return
	}

	err = tx.Commit()
	if err != nil {
		log.Print("MethodGetPostDetails, Failed to commit transaction")
		return
	}

	return
}

func MethodPostUpdate(post_id string, newMessage string) (post models.PostWithEdited, responsecode int) {
	conn, err := Connect()
	tx, err := conn.Begin()
	if err != nil {
		log.Print("MethodPostUpdate, Cannot begin transaction")
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Print("MethodPostUpdate, Cannot close connection")
		}
	}()

	rows, errmessage := tx.Query(`SELECT * FROM "UpdatePostDetails"($1,$2)`, post_id, newMessage)
	if errmessage != nil {
		if err, ok := errmessage.(*pq.Error); ok {
			if err.Code == "P0001" {
				responsecode = 304
			} else {
				responsecode = 404
			}
			fmt.Println(err)
		}
		responsecode = 404
		err := tx.Rollback()
		if err != nil {
			log.Print("MethodPostUpdate, Failed to commit transaction")
			return
		}
		return
	}
	for rows.Next() {
		err = rows.Scan(
			&post.Author,
			&post.Created,
			&post.Forum,
			&post.Id,
			&post.IsEdited,
			&post.Message,
			&post.Parent,
			&post.Thread,
		)
		if err != nil {
			fmt.Println("MethodPostUpdate rows.Next()", err)
			return
		}
	}

	err = rows.Err()
	if err != nil {
		fmt.Println("MethodPostUpdate rows.Err())", err)

		return
	}
	responsecode = 200

	err = tx.Commit()
	if err != nil {
		log.Print("MethodPostUpdate, Failed to commit transaction")
		return
	}

	return
}

func MethodGetServiceStatus() (status models.Status, responsecode int) {
	conn, err := Connect()
	tx, err := conn.Begin()
	if err != nil {
		log.Print("MethodGetServiceStatus, Cannot begin transaction")
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Print("MethodGetServiceStatus, Cannot close connection")
		}
	}()

	rows, errmessage := tx.Query(`SELECT * FROM "GetServiceStatus"()`)
	if errmessage != nil {
		if err, ok := errmessage.(*pq.Error); ok {

			fmt.Println(err)
		}
		responsecode = 404
		err := tx.Rollback()
		if err != nil {
			log.Print("MethodGetServiceStatus, Failed to commit transaction")
			return
		}
		return
	}
	for rows.Next() {
		err = rows.Scan(
			&status.Forum,
			&status.Post,
			&status.Thread,
			&status.User,
		)
		if err != nil {
			fmt.Println("MethodGetServiceStatus rows.Next()", err)
			return
		}
	}

	err = rows.Err()
	if err != nil {
		fmt.Println("MethodGetServiceStatus rows.Err())", err)

		return
	}
	responsecode = 200
	err = tx.Commit()
	if err != nil {
		log.Print("MethodGetServiceStatus, Failed to commit transaction")
		return
	}

	return

}
