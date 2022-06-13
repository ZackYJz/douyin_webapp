package dao

import (
	"go_webapp/internal/model"
	"go_webapp/pkg/app"

	"gorm.io/gorm"
)

func (d *Dao) DoFollow(fan *model.Fans) error {
	return d.db.Create(fan).Error
}

func (d *Dao) QueryByTwoIds(myId, publisherId string) (fans *model.Fans, err error) {
	row := d.db.Where("fan_id", myId).Where("publisher_id", publisherId).Find(&fans).RowsAffected
	if row == 0 {
		err = gorm.ErrRecordNotFound
	}
	return
}

func (d *Dao) UpdateIsFriend(fan *model.Fans) error {
	return d.db.Model(&model.Fans{}).Where("id = ?", fan.ID).Update("is_friend", fan.IsFriend).Error
}

func (d *Dao) DeleteFans(fan *model.Fans) error {
	return d.db.Unscoped().Delete(fan).Error
}

func (d *Dao) CountMyFollower(myId string) (int64, error) {
	var count int64
	err := d.db.Model(&model.Fans{}).Where("fan_id = ?", myId).Count(&count).Error
	return count, err
}

func (d *Dao) CountMyFans(myId string) (int64, error) {
	var count int64
	err := d.db.Model(&model.Fans{}).Where("publisher_id = ?", myId).Count(&count).Error
	return count, err
}

func (d *Dao) QueryMyFollowerList(myId string, page, pageSize int) ([]*model.FollowerVO, error) {
	var follows []*model.FollowerVO
	err := d.db.Raw(`
		SELECT
			u.user_id as user_id,
			u.nickname as nickname,
			u.photo as photo
		FROM
			tb_fans f
		LEFT JOIN
			tb_user u
		ON
			f.publisher_id = u.user_id
		WHERE
			f.fan_id = ?
		ORDER BY
			f.created_at
		DESC
		LIMIT ? OFFSET ?
	`, myId, pageSize, app.GetPageOffset(page, pageSize)).Scan(&follows).Error
	return follows, err
}
func (d *Dao) QueryMyFansList(myId string, page, pageSize int) ([]*model.FansVO, error) {
	var fans []*model.FansVO
	err := d.db.Raw(`
		SELECT
			u.user_id as userId,
			u.nickname as nickname,
			u.photo as photo
		FROM
			tb_fans f
		LEFT JOIN
			tb_user u
		ON
			f.fan_id = u.user_id
		WHERE
			f.publisher_id = ?
		ORDER BY
			f.created_at
		DESC
		LIMIT ? OFFSET ?
	`, myId, pageSize, app.GetPageOffset(page, pageSize)).Scan(&fans).Error
	return fans, err
}
