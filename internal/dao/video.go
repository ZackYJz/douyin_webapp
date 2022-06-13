package dao

import (
	"go_webapp/global"
	"go_webapp/internal/model"
	"go_webapp/pkg/app"
	"go_webapp/pkg/util"
)

func (d *Dao) Insert(v *model.Video) error {
	return d.db.Create(v).Error
}

func (d *Dao) CountVideoIndexList() (int, error) {
	var count int64
	err := d.db.Model(&model.Video{}).Where("is_private", global.NO).Count(&count).Error
	return int(count), err
}

func (d *Dao) GetVideoIndexList(search string, page, pageSize int) ([]*model.Video, error) {
	var videos []*model.Video
	var err error
	if search != "" {
		err = d.db.Where("is_private", global.NO).
			Where("title LIKE ?", "%"+search+"%").
			Limit(pageSize).Offset(app.GetPageOffset(page, pageSize)).
			Order("created_at desc").
			Find(&videos).Error
	} else {
		err = d.db.Where("is_private", global.NO).
			Limit(pageSize).Offset(app.GetPageOffset(page, pageSize)).
			Order("created_at desc").
			Find(&videos).Error
	}
	return videos, err
}

func (d *Dao) GetVlogById(videoId string) (video *model.Video, err error) {
	err = d.db.First(&video, util.StrTo(videoId).MustInt64()).Error
	return
}

func (d *Dao) ChangeVideoStatus(userId, videoId string, status int) error {
	return d.db.Model(&model.Video{}).Where("publisher = ?", userId).Where("id = ?", videoId).
		//TODO 小心 0 没有被更新
		Update("is_private", status).Error
}

func (d *Dao) CountMyVideoList(userId string, isPrivate int) (int, error) {
	var count int64
	err := d.db.Model(&model.Video{}).Where("publisher = ?", userId).Where("is_private", isPrivate).Count(&count).Error
	return int(count), err
}

func (d *Dao) QueryMyVideoList(userId string, isPrivate, page, pageSize int) ([]*model.VideoSimpleVO, error) {
	var videoList []*model.VideoSimpleVO
	err := d.db.Where(map[string]interface{}{"publisher": userId, "is_private": isPrivate}).
		Limit(pageSize).Offset(app.GetPageOffset(page, pageSize)).
		Find(&videoList).Error
	return videoList, err
}

func (d *Dao) QueryMyLikedVideoList(userId string, page, pageSize int) ([]*model.VideoVO, error) {
	var likedVideoList []*model.VideoVO
	err := d.db.Raw(`
		SELECT
			v.id as id,
			v.publisher as publisher,
			v.title as title,
			v.url as url,
			v.cover as cover,
			v.width as width,
			v.height as height,
		v.like_counts as like_counts,
			v.comments_counts as comment_counts,
			v.is_private as is_private
		FROM
			tb_video v
		LEFT JOIN
			tb_like l
		ON 
			v.id = l.video_id
		LEFT JOIN
			tb_user u
		ON
			l.user_id = u.user_id
		WHERE
			u.user_id = ? AND v.is_private = 0
		ORDER BY v.created_at DESC
		LIMIT ? OFFSET ?
	`, userId, pageSize, app.GetPageOffset(page, pageSize)).Scan(&likedVideoList).Error
	if err != nil {
		return nil, err
	}
	return likedVideoList, nil
}

func (d *Dao) CountMyFollowVideo(userId string) (int, error) {
	var count int64
	err := d.db.Raw(`
		SELECT
		 COUNT(*)
		FROM
			tb_video v
		LEFT JOIN
			tb_fans f
		ON v.publisher = f.publisher_id
		LEFT JOIN
			tb_user u
		ON
			f.publisher_id = u.user_id
		WHERE
			v.is_private = 0
			AND
			f.fan_id = ?
	`, userId).Scan(&count).Error
	return int(count), err
}

func (d *Dao) CountMyFriendVideo(userId string) (int, error) {
	var count int64
	err := d.db.Raw(`
		SELECT
			count(*)
		FROM
			tb_video v
		LEFT JOIN
			tb_fans f
		On
			v.publisher = f.fan_id
		LEFT JOIN
			tb_user u
		ON
			f.fan_id = u.user_id
		WHERE
			v.is_private = 0
			AND
			f.publisher_id = ?
			AND
			f.is_friend = 1
	`, userId).Scan(&count).Error
	return int(count), err
}

func (d *Dao) QueryMyFollowVideoList(userId string, page, pageSize int) ([]*model.VideoVO, error) {
	var myFollowVideoList []*model.VideoVO
	err := d.db.Raw(`
		SELECT
		 v.id as id,
		 v.publisher as publisher,
		 u.photo as publisher_face,
		 u.nickname as publisher_name,
		 v.title as title,
		 v.url as url,
		 v.cover as cover,
			v.width as width,
			v.height as height,
			v.like_counts as like_counts,
			v.comments_counts as comment_counts,
			v.is_private as is_private
		FROM
			tb_video v
		LEFT JOIN
			tb_fans f
		ON v.publisher = f.publisher_id
		LEFT JOIN
			tb_user u
		ON
			f.publisher_id = u.user_id
		WHERE
			v.is_private = 0
			AND
			f.fan_id = ?
		ORDER BY v.created_at DESC
		LIMIT ? OFFSET ?
	`, userId, pageSize, app.GetPageOffset(page, pageSize)).Scan(&myFollowVideoList).Error
	if err != nil {
		return nil, err
	}
	return myFollowVideoList, nil
}

func (d *Dao) QueryMyFriendVideoList(userId string, page, pageSize int) ([]*model.VideoVO, error) {
	var myFriendVideoList []*model.VideoVO
	err := d.db.Raw(`
		SELECT
	 v.id as id,
		 v.publisher as publisher,
		 u.photo as publisher_face,
		 u.nickname as publisher_name,
		 v.title as title,
		 v.url as url,
		 v.cover as cover,
			v.width as width,
			v.height as height,
			v.like_counts as like_counts,
			v.comments_counts as comment_counts,
			v.is_private as is_private
		FROM
			tb_video v
		LEFT JOIN
			tb_fans f
		On
			v.publisher = f.fan_id
		LEFT JOIN
			tb_user u
		ON
			f.fan_id = u.user_id
		WHERE
			v.is_private = 0
			AND
			f.publisher_id = ?
			AND
			f.is_friend = 1
		ORDER BY v.created_at DESC
		LIMIT ? OFFSET ?
	`, userId, pageSize, app.GetPageOffset(page, pageSize)).Scan(&myFriendVideoList).Error
	if err != nil {
		return nil, err
	}
	return myFriendVideoList, nil
}
