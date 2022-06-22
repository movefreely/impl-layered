package service

import (
	"errors"
	"in-server/model"
	"in-server/util"
	"in-server/util/sqlinit"
	"time"
)

type UserService struct{}

// UserRegister 注册
func (u *UserService) UserRegister(email, pwd, nickname string) (user model.User, err error) {
	register := model.User{}
	sqlinit.DB.Where("email = ?", email).Find(&register)
	// 账号已创建
	if register.ID > 0 {
		return register, errors.New("账号已创建")
	}
	// 创建账号
	register.Email = email
	register.Password = util.Md5(pwd)
	register.Nickname = nickname
	register.CreateAt = time.Now()
	register.Online = "0"

	sqlinit.DB.Save(&register)

	//err = sqlinit.DB.Table(util.Md5(email)[0:16]).AutoMigrate(&model.Message{})
	if err != nil {
		return model.User{}, err
	}

	return register, nil
}

// UserLogin 登录
func (u *UserService) UserLogin(email, pwd string) (user model.User, err error) {
	login := model.User{}
	sqlinit.DB.Where("email = ?", email).Find(&login)
	// 账号不存在
	if login.ID == 0 {
		return login, errors.New("账号不存在")
	}
	// 密码错误
	if login.Password != util.Md5(pwd) {
		return login, errors.New("密码错误")
	}
	//db.Model(&User{}).Where("active = ?", true).Update("name", "hello")
	sqlinit.DB.Model(&model.User{}).Where("id = ?", login.ID).Update("online", "1")
	return login, nil
}

// UserInfo 查询用户信息
func (u *UserService) UserInfo(userid uint64) (model.User, error) {
	user := model.User{}
	sqlinit.DB.Where("id = ?", userid).Find(&user)
	if user.ID == 0 {
		return user, errors.New("账号不存在")
	} else {
		return user, nil
	}
}

// ChangeAvatar 修改头像
func (u *UserService) ChangeAvatar(userid uint64, avatar string) error {
	err := sqlinit.DB.Model(&model.User{}).Where("id = ?", userid).Update("avatar", avatar).Error
	if err != nil {
		return err
	}
	return nil
}
