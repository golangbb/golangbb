package models

import "errors"

var ErrEmptyUserID = errors.New("empty UserID/AuthorID not allowed")
var ErrEmptyUserName = errors.New("empty UserName not allowed")
var ErrEmptyPassword = errors.New("empty Password not allowed")
var ErrEmptyName = errors.New("empty Name not allowed")
var ErrEmptyTitle = errors.New("empty Title not allowed")
var ErrEmptyTopicID = errors.New("empty TopicID not allowed")
var ErrDiscussionWithoutSinglePost = errors.New("a Discussion must be created with a single Post")

func Models() []interface{} {
	return []interface{}{
		&Discussion{}, &Email{}, &Group{}, &Post{}, &Topic{}, &User{},
	}
}
