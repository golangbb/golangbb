package models

func Models() []interface{} {
	return []interface{}{
		&Email{}, &Group{}, &User{}, &Topic{}, &Post{},
	}
}
