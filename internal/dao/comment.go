package dao

import (
	"go_webapp/internal/model"
	"go_webapp/pkg/app"
)

func (d *Dao) InsertComment(c *model.Comment) error {
	return d.db.Create(c).Error
}

func (d *Dao) CountComment(videoId string) (int64, error) {
	var count int64
	err := d.db.Model(&model.Comment{}).Where("video_id = ?", videoId).Count(&count).Error
	return count, err
}

func (d *Dao) GetCommentVOList(videoId string, page, pageSize int) ([]*model.CommentVO, error) {
	var comments []*model.CommentVO
	err := d.db.Raw(`
		SELECT
			c.id as comment_id,
			c.video_id as video_id,
			c.publisher as publisher,
			u.nickname as nickname,
			u.photo as face,
			c.father_id as father_id,
			c.user_id as user_id,
			c.content as content,
			c.like_counts as like_counts,
			fu.nickname as replyed_user_nickname,
			c.created_at as create_time
		FROM
			tb_comment as c
		LEFT JOIN
			tb_user as u
		ON 
			c.user_id = u.user_id
		LEFT JOIN
			tb_comment as fc
		ON
			c.father_id = fc.id
		LEFT JOIN
			tb_user as fu
		ON
			fc.user_id = fu.user_id
		WHERE
			c.video_id = ?
		ORDER BY
			c.like_counts DESC,
			c.created_at DESC
		LIMIT ? OFFSET ?
	`, videoId, pageSize, app.GetPageOffset(page, pageSize)).Scan(&comments).Error
	return comments, err
}

func (d *Dao) DeleteComment(id string) error {
	return d.db.Unscoped().Delete(&model.Comment{}, id).Error
}
