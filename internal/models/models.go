package models

import "errors"

var ErrEmptyUserID = errors.New("empty UserID/AuthorID not allowed")
var ErrEmptyUserName = errors.New("empty UserName not allowed")
var ErrEmptyPassword = errors.New("empty Password not allowed")
var ErrEmptyName = errors.New("empty Name not allowed")

func Models() []interface{} {
	return []interface{}{
		&Discussion{}, &Email{}, &Group{}, &Post{}, &Topic{}, &User{},
	}
}
