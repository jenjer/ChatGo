package Globals

import (
	"sync"
)

type User struct {
	id   string
	name string
}

var (
	currentUser *User
	once        sync.Once
)

func getInstance() *User {
	once.Do(func() {
		currentUser = &User{}
	})
	return currentUser
}

func SetID(ID string) {
	user := getInstance()
	user.id = ID
}

func GetID()(string) {
	user := getInstance()
	return user.id
}
