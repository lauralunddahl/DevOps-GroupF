package dto

import database "db"

type Follower struct {
	WhoId  int
	WhomId int
}

func IsFollowing(who_id int, whom_id int) bool {
	follower := getFollower(who_id, whom_id)
	return follower.WhoId != 0
}

func FollowUser(who_id int, whom_id int) {
	follower := Follower{WhoId: who_id, WhomId: whom_id}
	result := database.DB.Create(&follower)
	if result.Error != nil {
		print(result.Error)
	}
}

func getFollower(who_id int, whom_id int) Follower {
	follower := Follower{}
	database.DB.Where("who_id = ? and whom_id = ?", who_id, whom_id).First(&follower)
	return follower
}

func UnfollowUser(who_id int, whom_id int) {
	follower := getFollower(who_id, whom_id)
	result := database.DB.Delete(&follower)
	if result.Error != nil {
		print(result.Error)
	}
}
