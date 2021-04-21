package dto

import (
	database "github.com/lauralunddahl/DevOps-GroupF/program/db"
	log "github.com/sirupsen/logrus"
)

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
		log.Println("FollowUser")
		log.Error(result.Error)
	}
}

func getFollower(who_id int, whom_id int) Follower {
	follower := Follower{}
	res := database.DB.Where("who_id = ? and whom_id = ?", who_id, whom_id).First(&follower)
	if res.Error != nil {
		log.Println("getFollower")
		log.Error(res.Error)
	}
	return follower
}

func GetFollowers(who_id int, limit int) []Follower {
	var followers []Follower
	res := database.DB.Table("followers").Where("who_id = ?", who_id).Limit(limit).Scan(&followers)
	if res.Error != nil {
		log.Println("GetFollowers")
		log.Error(res.Error)
	}
	return followers
}

func UnfollowUser(who_id int, whom_id int) {
	follower := getFollower(who_id, whom_id)
	result := database.DB.Where("who_id = ? and whom_id = ?", who_id, whom_id).Delete(&follower)
	if result.Error != nil {
		log.Println("UnfollowUser")
		log.Error(result.Error)
	}
}
