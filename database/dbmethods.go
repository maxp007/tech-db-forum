package database

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/maxp007/tech-db-forum/models"
	lg "log"
)

var log lg.Logger

func MethodCreateOrGetForum(in *models.ForumRequest) (forum models.ForumResponse, responsecode int) {
	conn, err := Connect()
	tx, err := conn.Begin()
	if err != nil {
		log.Print("CreateOrGetForum, Cannot begin transaction")
	}

	rows, errmessage := tx.Query(`SELECT * FROM public."CreateOrGetForum"($1,$2,$3)`, in.Title, in.Slug, in.User)
	if errmessage != nil {
		if err, ok := errmessage.(*pq.Error); ok {
			if err.Code == "23505" {
				responsecode = 409
				err := tx.Rollback()
				if err != nil {
					log.Print("MethodCreateOrGetForum, Failed to commit transaction")
					return
				}
				tx, err := conn.Begin()
				if err != nil {
					log.Print("CreateOrGetForum, Cannot begin transaction")
				}

				rows, errmessage := tx.Query(`SELECT * FROM public."Forum" WHERE slug=$1`, in.Slug)
				if errmessage == nil {
					for rows.Next() {
						err = rows.Scan(
							&forum.Posts,
							&forum.Slug,
							&forum.Threads,
							&forum.Title,
							&forum.User,
						)
						if err != nil {
							fmt.Println("Failed to scan rows")
							return
						}
					}

					err = rows.Err()
					if err != nil {
						return
					}

					err = tx.Commit()
					if err != nil {
						log.Print("CreateOrGetForum, Failed to commit transaction")
						return
					}

					return

				} else {

					fmt.Println("Failed to get similar users")
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
			&forum.Posts,
			&forum.Slug,
			&forum.Threads,
			&forum.Title,
			&forum.User)
		if err != nil {
			fmt.Println("CreateOrGetForum rows.Next()", err)
			return
		}
	}

	err = rows.Err()
	if err != nil {
		fmt.Println("CreateOrGetForum rows.Err())", err)

		return
	}

	err = tx.Commit()
	if err != nil {
		log.Print("CreateOrGetForum, Failed to commit transaction")
		return
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Print("CreateOrGetForum, Cannot close connection")
		}
	}()
	return
}

func MethodGetForumDetails(slug string) (forum models.ForumRequest, responsecode int) {
	conn, err := Connect()

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Print("CreateOrGetForum, Cannot close connection")
		}
	}()

	rows, errmessage := conn.Query(`SELECT * FROM public."Forum" WHERE slug=$1::citext`, slug)
	if errmessage == nil {
		responsecode = 404
		for rows.Next() {

			err = rows.Scan(
				&forum.Posts,
				&forum.Slug,
				&forum.Threads,
				&forum.Title,
				&forum.User,
			)
			if err != nil {
				fmt.Println("Failed to scan rows")
				return
			}
			responsecode = 200
		}
		err = rows.Err()
		if err != nil {
			return
		}
		return
	}
	return
}

func MethodCreateOrGetUser(in *models.User) (users []models.User, violationflag int) {
	conn, err := Connect()
	tx, err := conn.Begin()
	if err != nil {
		log.Print("CreateOrGetForum, Cannot begin connection")
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Print("CreateOrGetForum, Cannot close connection")
		}
	}()

	rows, errmessage := tx.Query(`SELECT * FROM public."InsertUser"($1,$2,$3,$4)`, in.About, in.Email, in.Fullname, in.Nickname)
	if errmessage == nil {
		violationflag = 0
		var user models.User

		for rows.Next() {
			err := rows.Scan(
				&user.About,
				&user.Email,
				&user.Fullname,
				&user.Nickname,
			)
			if err != nil {
				fmt.Println("Failed to scan rows")
				return
			}
			users = append(users, user)
		}

		err := rows.Err()
		if err != nil {
			return
		}

		err = tx.Commit()
		if err != nil {
			log.Print("CreateOrGetForum, Failed to commit transaction")
			return
		}

		return

	} else {
		violationflag = 1

		err_ := tx.Rollback()
		if err_ != nil {
			log.Print("MethodCreateOrGetUser, Failed to Rollback transaction")
		}

		tx, err := conn.Begin()

		if err != nil {
			log.Print("CreateOrGetForum, Cannot begin transaction")
		}

		rows, errmessage := tx.Query(`SELECT * FROM public."GetSimilarUsers"($1,$2)`, in.Email, in.Nickname)
		if errmessage == nil {
			violationflag = 1
			var user models.User
			for rows.Next() {
				err = rows.Scan(
					&user.About,
					&user.Email,
					&user.Fullname,
					&user.Nickname,
				)
				if err != nil {
					fmt.Println("Failed to scan rows")
					return
				}
				users = append(users, user)
			}

			err = rows.Err()
			if err != nil {
				return
			}

			err = tx.Commit()
			if err != nil {
				log.Print("CreateOrGetForum, Failed to commit transaction")
				return
			}

			return

		} else {

			fmt.Println("Failed to get similar users")
		}

		return
	}
}

func MethodGetUserProfile(nickname string) (user models.User, violationflag int) {
	conn, err := Connect()

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Print("CreateOrGetForum, Cannot close connection")
		}
	}()

	row := conn.QueryRow(`SELECT * FROM public."User" Where nickname=$1`, nickname)
	err = row.Scan(
		&user.About,
		&user.Email,
		&user.Fullname,
		&user.Nickname,
	)
	if err != nil {
		violationflag = 1
	} else {
		violationflag = 0
	}

	return
}

func MethodUpdateUserProfile(userprofile *models.User) (user models.User, violationflag int) {
	conn, err := Connect()
	tx, err := conn.Begin()
	if err != nil {
		log.Print("CreateOrGetForum, Cannot begin transaction")
		return
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Print("CreateOrGetForum, Cannot close connection")
		}
	}()

	rows, errmessage := tx.Query(
		`SELECT * FROM "UpdateUserProfile"($1,$2,$3,$4)`,
		userprofile.About, userprofile.Email, userprofile.Fullname, userprofile.Nickname)

	if errmessage != nil {
		if err, ok := errmessage.(*pq.Error); ok {
			if err.Code == "23505" {
				violationflag = 2
				err := tx.Rollback()
				if err != nil {
					log.Print("CreateOrGetForum, Failed to commit transaction")
					return
				}
				return
			} else {
				violationflag = 1
				err := tx.Rollback()
				if err != nil {
					log.Print("CreateOrGetForum, Failed to commit transaction")
					return
				}
				return
			}
		}

	} else {
		violationflag = 0
		for rows.Next() {
			err := rows.Scan(
				&user.About,
				&user.Email,
				&user.Fullname,
				&user.Nickname,
			)
			if err != nil {
				fmt.Println("Failed to scan rows")
				return
			}
		}

	}

	err = tx.Commit()
	if err != nil {
		log.Print("CreateOrGetForum, Failed to commit transaction")
		return
	}

	return
}

func MethodGetForumUsers(forum_slug string, limit string, since string, desc string) (users []models.User, responsecode int) {
	conn, err := Connect()
	if err != nil {
		log.Print("MethodGetForumUsers, Cannot open connection")
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Print("MethodGetForumUsers, Cannot close connection")
		}
	}()

	rows, errmessage := conn.Query(`SELECT * FROM "GetForumUsers"($1,$2,$3,$4)`, forum_slug, limit, since, desc)
	if errmessage == nil {
		if err, ok := errmessage.(*pq.Error); ok {
			if err.Code == "P0002" {
				responsecode = 404
			} else {
				responsecode = 404

			}
		}
		var user models.User
		for rows.Next() {
			err := rows.Scan(
				&user.About,
				&user.Email,
				&user.Fullname,
				&user.Nickname,
			)
			if err != nil {
				fmt.Println("Failed to scan rows")
				return
			}
			users = append(users, user)
		}
		err := rows.Err()
		if err != nil {
			return
		}

		responsecode = 200
		return
	} else {
		responsecode = 404
		return
	}
}

func ServiceCleanData() (responsecode int) {
	conn, err := Connect()

	if err != nil {
		log.Print("ServiceCleanData, Cannot begin connection")
	}

	_, errmessage := conn.Query(`SELECT * FROM "clearalldata"()`)
	if errmessage == nil {
		responsecode = 200
	} else {
		responsecode = 404
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Print("ServiceCleanData, Cannot close connection")
		}
	}()
	return
}
