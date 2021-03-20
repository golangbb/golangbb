package models

func Models() []interface{} {
	return []interface{}{
		&Discussion{}, &Email{}, &Group{}, &Post{}, &Topic{}, &User{},
	}
}
