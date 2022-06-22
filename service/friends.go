package service

import (
	"errors"
	"fmt"
	"in-server/model"
	"in-server/util/sqlinit"
	"strconv"
)

type FriendServer struct{}

func (f *FriendServer) AddFriend(oneId, otherId uint64) (bool, error) {
	if oneId == otherId {
		return false, errors.New("不能添加自己为好友")
	}
	friend := model.Friend{}
	sqlinit.DB.Where("one_id = ? and other_id = ?", oneId, otherId).First(&friend)
	//fmt.Printf("%+v\n", friend)
	if friend.Id != 0 {
		return false, errors.New("已经发送请求或已添加好友")
	}
	sqlinit.DB.Save(&model.Friend{OneId: oneId, OtherId: otherId})
	return true, nil
}

func (f *FriendServer) FindUserByNickname(name string) []model.User {
	var users []model.User
	sqlinit.DB.Where("nickname LIKE ?", "%"+name+"%").Find(&users)
	return users
}

func (f *FriendServer) GetFriends(selfId string) ([]model.User, error) {
	atoi, err := strconv.Atoi(selfId)
	if err != nil {
		return nil, errors.New("传入ID无法转为整数")
	}
	oneId := uint64(atoi)
	var friends []model.Friend
	sqlinit.DB.Where("one_id = ? or other_id = ?", oneId, oneId).Find(&friends)
	var users []model.User
	for _, friend := range friends {
		var id uint64
		if friend.OneId == oneId {
			id = friend.OtherId
		} else {
			id = friend.OneId
		}
		var user model.User
		sqlinit.DB.Where("id = ? and status = 1", id).First(&user)
		users = append(users, user)
	}
	return users, nil
}

func (f *FriendServer) AgreeFriend(selfId, friendId string) (bool, error) {
	oneId, err := strconvToUint64(selfId)
	otherId, err := strconvToUint64(friendId)
	if err != nil {
		return false, errors.New("传入ID无法转为整数")
	}
	var friend model.Friend
	sqlinit.DB.Where("one_id = ? and other_id = ?", otherId, oneId).First(&friend).Update("status", 1)
	sqlinit.DB.Table("messages").Where("from_id = ? and to_id = ?", otherId, oneId).Delete(model.Message{})
	if friend.Id == 0 {
		return false, errors.New("没有找到该好友")
	}
	return true, nil
}

func (f *FriendServer) RejectFriend(selfId, friendId string) (bool, error) {
	oneId, err := strconvToUint64(selfId)
	otherId, err := strconvToUint64(friendId)
	if err != nil {
		return false, errors.New("传入ID无法转为整数")
	}
	var friend model.Friend
	sqlinit.DB.Where("one_id = ? and other_id = ?", otherId, oneId).First(&friend)
	if friend.Id == 0 {
		return false, errors.New("没有找到该好友")
	}
	sqlinit.DB.Delete(&friend)
	sqlinit.DB.Table("messages").Where("from_id = ? and to_id = ?", otherId, oneId).Delete(model.Message{})
	sqlinit.DB.Table("latest_message").Where("from_id = ? and to_id = ?", otherId, oneId).Delete(model.Message{})
	return true, nil
}

func strconvToUint64(id string) (uint64, error) {
	atoi, err := strconv.Atoi(id)
	fmt.Println(atoi)
	if err != nil {
		return 0, errors.New("传入ID无法转为整数")
	}
	return uint64(atoi), nil
}
