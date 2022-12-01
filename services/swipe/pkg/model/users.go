package users

import (
	"strconv"
	"strings"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	Email        string
	Firstname    string
	Birthday     string
	About        string
	Photos       pq.StringArray
	Tags         pq.StringArray
	Gender       int
	Distance     int
	ShowMe       int
	FilterByTags int
}

type Photo struct {
	User  string
	Photo string
}

type Tag struct {
	User string
	Tag  string
}

// Get users from DB for swiper
func GetUsers(db *gorm.DB, t *User) ([]User, error) {
	query := "SELECT * FROM users WHERE distance <= " + strconv.Itoa(t.Distance)
	if t.ShowMe != 2 {
		query += " AND gender = " + strconv.Itoa(t.ShowMe)
	}

	users, err := GetUsersFromQuery(db, query)
	if err != nil {
		return nil, err
	}

	users, err = getPhotosByUsers(db, users)
	if err != nil {
		return nil, err
	}

	users, err = getTagsByUsers(db, users)
	if err != nil {
		return nil, err
	}

	if t.FilterByTags == 1 {
		users = FilterUsersByTags(t, users)
	}

	return users, nil
}

func GetUsersFromQuery(db *gorm.DB, query string) ([]User, error) {
	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User
	for rows.Next() {
		db.ScanRows(rows, &users)
	}

	return users, nil
}

func getPhotosByUsers(db *gorm.DB, users []User) ([]User, error) {
	var emails []string
	for _, user := range users {
		emails = append(emails, "'"+user.Email+"'")
	}
	emailStrings := strings.Join(emails, ", ")

	query := "SELECT * FROM photos WHERE user IN (" + emailStrings + ")"

	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var photos []Photo
	for rows.Next() {
		db.ScanRows(rows, &photos)
	}

	for i, user := range users {
		var arr []string

		for _, photo := range photos {
			if user.Email == photo.User {
				arr = append(arr, photo.Photo)
			}
		}

		users[i].Photos = arr
	}

	return users, nil
}

func getTagsByUsers(db *gorm.DB, users []User) ([]User, error) {
	var emails []string
	for _, user := range users {
		emails = append(emails, "'"+user.Email+"'")
	}
	emailStrings := strings.Join(emails, ", ")

	query := "SELECT * FROM tags WHERE user IN (" + emailStrings + ")"

	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		db.ScanRows(rows, &tags)
	}

	for i, user := range users {
		var arr []string

		for _, tag := range tags {
			if user.Email == tag.User {
				arr = append(arr, tag.Tag)
			}
		}

		users[i].Tags = arr
	}

	return users, nil
}

func FilterUsersByTags(t *User, users []User) []User {
	var arr []User
	for _, user := range users {
		containsTags := contains(t.Tags, user.Tags)

		if containsTags {
			arr = append(arr, user)
		}
	}

	return arr
}

func contains(s []string, e []string) bool {
	for _, a := range s {
		for _, b := range e {
			if a == b {
				return true
			}
		}
	}
	return false
}
