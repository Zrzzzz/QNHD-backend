package models

import (
	"gorm.io/gorm"
)

type postFreq struct {
	Id        uint64 `json:"id"`
	FreqLevel uint64 `json:"freq_level"`
}

type postStat struct {
	Id  uint64 `json:"id"`
	Cnt uint64 `json:"cnt"`
}

func RefreshPostFreq() error {
	// 把所有帖子频率+1
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&Post{}).Update("freq_level", gorm.Expr("freq_level + ?", 1)).Error; err != nil {
			return err
		}

		var posts []postFreq
		if err := tx.Model(&Post{}).Select("id", "freq_level").Find(&posts).Error; err != nil {
			return err
		}
		var postFreqMap = make(map[uint64]uint64)
		for _, post := range posts {
			postFreqMap[post.Id] = post.FreqLevel
		}

		// 过去7天的点赞数
		var postLikes []postStat
		if err := tx.Raw(`
		SELECT p.id AS id, COUNT(l.post_id) AS cnt
		FROM qnhd.post p
		LEFT JOIN qnhd.log_post_like l ON p.id = l.post_id
		WHERE l.created_at >= NOW() - INTERVAL '7 days'
		GROUP BY p.id
		`).Scan(&postLikes).Error; err != nil {
			return err
		}

		// 过去7天的评论数，单人单帖单天只算一次
		var postComments []postStat
		if err := tx.Raw(`
		WITH daily_comments AS (
			SELECT post_id, uid, DATE_TRUNC('day', created_at) AS day
			FROM qnhd.floor
			WHERE created_at >= NOW() - INTERVAL '7 days'
			GROUP BY post_id, uid, day
		),
		total_comments AS (
			SELECT post_id, COUNT(*) AS total_comment_count
			FROM daily_comments
			GROUP BY post_id
		)
		SELECT p.id AS id, tc.total_comment_count AS cnt
		FROM qnhd.post p
		INNER JOIN total_comments tc ON p.id = tc.post_id
		ORDER BY total_comment_count DESC
		`).Scan(&postComments).Error; err != nil {
			return err
		}

		// 点赞+1分，评论+5分
		var postScoreMap = make(map[uint64]uint64)
		for _, postLike := range postLikes {
			postScoreMap[postLike.Id] += postLike.Cnt
		}
		for _, postComment := range postComments {
			postScoreMap[postComment.Id] += postComment.Cnt * 5
		}

		// 如果分数大于freq * 10 + 10则将level置0
		var updatePosts []uint64
		for id, score := range postScoreMap {
			if score > postFreqMap[id]*10+10 {
				updatePosts = append(updatePosts, id)
			}
		}

		if err := tx.Model(&Post{}).Where("id IN (?)", updatePosts).Update("freq_level", 0).Error; err != nil {
			return err
		}

		return nil
	})
	return err
}
