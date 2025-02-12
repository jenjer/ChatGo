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

func SetID(ID string) {
	currentUser.id = ID
}

func GetID()(string) {
	return currentUser.id
}
