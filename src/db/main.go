package main

import (
	"db/user"
)

func main() {
	user := user.GetUserID("laulu")
	println(user)
}
