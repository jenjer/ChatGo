package Login

import (
	"fmt"
	"sync"
	"github.com/jenjer/ChatGo/internal"

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

func tryLogin(id string, pw string) {
}

func Login() {
	var ID string
	var PW string
	fmt.Print("Input your Login ID : ")
	fmt.Scan(&ID)

	fmt.Print("Input your Login PW : ")
	fmt.Scan(&PW)

	fmt.Printf("ID : %s, PW : %s\n", ID, PW)
}
