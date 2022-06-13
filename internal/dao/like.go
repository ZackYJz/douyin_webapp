package dao

import "go_webapp/internal/model"

func (d *Dao) InsertLike(like *model.Like) (int, error) {
	tx := d.db.Create(like)
	return int(tx.RowsAffected), tx.Error
}

func (d *Dao) DeleteLike(like *model.Like) (int, error) {
	tx := d.db.Where("user_id = ?", like.UserId).Where("video_id = ?", like.VideoId).Unscoped().Delete(like)
	return int(tx.RowsAffected), tx.Error
}

func (d *Dao) QueryLike(likeId int64) (l *model.Like, err error) {
	err = d.db.First(&l, likeId).Error
	return
}

func (d *Dao) CountMyLikedVideo(userId string) (int, error) {
	var count int64
	err := d.db.Model(&model.Like{}).Where("user_id = ?", userId).Count(&count).Error
	return int(count), err
}
