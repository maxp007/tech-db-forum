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
	}

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

	rows, errmessage = tx.Query(`SELECT * FROM public."Forum" where slug=$1::citext `, forum_slug)
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
			rows, errmessage = tx.Query(`SELECT * FROM public."Thread" where forum=$1::citext  ORDER BY created DESC LIMIT $2::INTEGER;`,
				forum_slug, limit)
		} else {
			rows, errmessage = tx.Query(`SELECT * FROM public."Thread" where forum=$1::citext  ORDER BY created ASC LIMIT $2::INTEGER;`,
				forum_slug, limit)
		}
	} else {
		if limit == "4" {
			if desc == "true" {
				rows, errmessage = tx.Query(`SELECT * FROM public."Thread" where (forum=$1::citext AND created <= $2::timestamptz) ORDER BY created DESC LIMIT $3::INTEGER;`,
					forum_slug, since, limit)
			} else {
				rows, errmessage = tx.Query(`SELECT * FROM public."Thread" where (forum=$1::citext AND created >= $2::timestamptz) ORDER BY created ASC LIMIT $3::INTEGER;`,
					forum_slug, since, limit)
			}

		} else {

			if desc == "true" {
				rows, errmessage = tx.Query(`SELECT * FROM public."Thread" where (forum=$1::citext AND created >= $2::timestamptz) ORDER BY created DESC LIMIT $3::INTEGER;`,
					forum_slug, since, limit)
			} else {
				rows, errmessage = tx.Query(`SELECT * FROM public."Thread" where (forum=$1::citext AND created >= $2::timestamptz) ORDER BY created ASC LIMIT $3::INTEGER;`,
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

	if len(PostsSlice) == 0 {
		responsecode = 201
		err = tx.Commit()
		if err != nil {
			log.Print("MethodCreatePost, Failed to commit transaction")
			return
		}
		return
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

	rows, errmessage = tx.Query(`SELECT * FROM  public."CreatePostUsingFieldArrays"($1::citext[],$2::text[],$3::citext[],$4::integer,$5::citext,$6::integer,$7::bool,$8::bool)`,
		pq.Array(author_array),
		pq.Array(message_array),
		pq.Array(parent_array),
		len(PostsSlice),
		thread_slug,
		thread_id,
		parent_is_the_same,
		author_is_the_same,
	)

	if errmessage != nil {
		if err, ok := errmessage.(*pq.Error); ok {
			fmt.Println(err)
		}
		responsecode = 404
		err := tx.Rollback()
		if err != nil {
			log.Print("MethodCreatePost, Failed to commit transaction")
			return
		}
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
