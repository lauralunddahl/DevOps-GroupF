package main

import "dto"

func main() {
	//user := dto.GetUserID("laulu")
	dto.GetUsers()
	timeline := dto.GetPrivateTimeline(1)
	println(timeline)
}
