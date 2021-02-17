package main

import "controller"

func main() {
	user := controller.GetUserID("laulu")
	println(user)
}
