package dto

import (
	database "github.com/lauralunddahl/DevOps-GroupF/program/db"
	log "github.com/sirupsen/logrus"
)

type Follower struct {
	WhoId  int
	WhomId int
}

func IsFollowing(whoId int, whomId int) bool {
	follower := getFollower(whoId, whomId)
	return follower.WhoId != 0
}

func FollowUser(whoId int, whomId int) {
	follower := Follower{WhoId: whoId, WhomId: whomId}
	result := database.DB.Create(&follower)
	if result.Error != nil {
		log.Println("FollowUser")
		log.Error(result.Error)
	}
}

func getFollower(whoId int, whomId int) Follower {
	follower := Follower{}
	res := database.DB.Where("who_id = ? and whom_id = ?", whoId, whomId).First(&follower)
	if res.Error != nil {
		log.Println("getFollower")
		log.Error(res.Error)
	}
	return follower
}

func GetFollowers(whoId int, limit int) []Follower {
	var followers []Follower
	res := database.DB.Table("followers").Where("who_id = ?", whoId).Limit(limit).Scan(&followers)
	if res.Error != nil {
		log.Println("GetFollowers")
		log.Error(res.Error)
	}
	return followers
}

func UnfollowUser(whoId int, whomId int) {
	follower := getFollower(whoId, whomId)
	result := database.DB.Where("who_id = ? and whom_id = ?", whoId, whomId).Delete(&follower)
	if result.Error != nil {
		log.Println("UnfollowUser")
		log.Error(result.Error)
	}
}
