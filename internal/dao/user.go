package dao

import (
	"errors"
	"go_webapp/global"
	"go_webapp/internal/model"
	"strconv"
)

func (d *Dao) GetUserByMobile(mobile string) (user *model.User, err error) {
	err = d.db.Where(&model.User{Mobile: mobile}).First(&user).Error
	return
}

func (d *Dao) GetUserByUserId(userId string) (user *model.User, err error) {
	id, _ := strconv.ParseInt(userId, 10, 64)
	err = d.db.Where(&model.User{UserId: id}).First(&user).Error
	return
}

func (d *Dao) CreateUser(user *model.User) (err error) {
	err = d.db.Create(&user).Error
	return
}

func (d *Dao) UpdateImageByUserId(param *model.UpdateUserDTO) (user *model.User, err error) {
	row := d.db.Model(model.User{}).Omit("user_id").Where("user_id", param.UserId).Updates(param).RowsAffected
	if row == 0 {
		return nil, errors.New("修改用户信息失败")
	}
	return d.GetUserByUserId(param.UserId)
}

func (d *Dao) UpdateInfoByUserId(param *model.UpdateUserDTO, typeNum int) (user *model.User, err error) {
	if typeNum == global.SEX {
		d.db.Model(model.User{}).Select("sex").Where("user_id", param.UserId).Updates(param)
		return d.GetUserByUserId(param.UserId)
	}
	if typeNum == global.BIRTHDAY {
		d.db.Model(model.User{}).Select("birthday").Where("user_id", param.UserId).Updates(param)
		return d.GetUserByUserId(param.UserId)
	}
	row := d.db.Model(model.User{}).Omit("user_id").Where("user_id", param.UserId).Updates(param).RowsAffected
	if row == 0 {
		return nil, errors.New("修改用户信息失败")
	}
	return d.GetUserByUserId(param.UserId)
}
