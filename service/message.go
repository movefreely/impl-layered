package service

import (
	"errors"
	"in-server/model"
	"in-server/util/sqlinit"
)

type Message struct{}

func (m *Message) NewestMessage(fromId int64) ([]*model.Message, error) {
	msg := make([]*model.Message, 0)
	sqlinit.DB.Table("latest_message").Where("from_id = ? or to_id = ?", fromId, fromId).Order("create_at desc").Find(&msg)
	return msg, nil
}

func (m *Message) HistoryMessage(fromId, toId uint64, page int) ([]*model.Message, error) {
	msg := make([]*model.Message, 0)
	if page == 0 {
		sqlinit.DB.Table("messages").Where("(from_id = ? and to_id = ?) or (from_id = ? and to_id = ?)", fromId, toId, toId, fromId).Limit(20).Order("id desc").Find(&msg)
	} else {
		sqlinit.DB.Table("messages").Where("(from_id = ? and to_id = ?) or (from_id = ? and to_id = ?)", fromId, toId, toId, fromId).Order("id desc").Offset(page * 20).Limit(20).Find(&msg)
	}
	return msg, nil
}

func (m *Message) UpdateMessage(message *model.Message) error {
	tempMsg := &model.Message{}
	sqlinit.DB.Save(&message)
	sqlinit.DB.Table("latest_message").Where("(from_id = ? and to_id = ?) or (from_id = ? and to_id = ?)", message.FromId, message.ToId, message.ToId, message.FromId).First(&tempMsg)
	if tempMsg.Id == 0 {
		sqlinit.DB.Table("latest_message").Save(&message)
	} else {
		sqlinit.DB.Table("latest_message").Where("(from_id = ? and to_id = ?) or (from_id = ? and to_id = ?)", message.FromId, message.ToId, message.ToId, message.FromId).Updates(map[string]interface{}{"content": message.Content, "create_at": message.CreateAt, "media": message.Media, "url": message.Url, "is_ad": message.IsAd})
	}
	return nil
}

func (m *Message) UidNotOnline(message *model.Message) error {
	user := &model.User{}
	sqlinit.DB.Where("id = ?", message.ToId).Find(&user)
	if user.ID == 0 {
		return errors.New("user not exist")
	} else {
		err := m.UpdateMessage(message)
		if err != nil {
			return errors.New("in UidNotOnline: update message error")
		}
	}
	return nil
}
