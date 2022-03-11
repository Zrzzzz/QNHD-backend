package models

import (
	"fmt"
	"qnhd/pkg/segment"
	"strings"

	"gorm.io/gorm"
)

// 转义引号
func escapeString(a string) string {
	return strings.ReplaceAll(a, "'", "''")
}

// 切割字符串数组按照顺序排权重
func geneTokenString(strs ...string) string {
	var (
		cutStrs []string
		tokens  []string
		weights = []string{"A", "B", "C", "D"}
	)
	if len(strs) > 4 {
		cutStrs = strs[:4]
	} else {
		cutStrs = strs
	}
	if len(strs) == 1 {
		t := segment.Cut(strs[0], " ")
		tokens = append(tokens, fmt.Sprintf("to_tsvector('simple', '%s')", t))
	} else {
		for i, s := range cutStrs {
			t := segment.Cut(s, " ")
			tokens = append(tokens, fmt.Sprintf("setweight(to_tsvector('simple', '%s'), '%s')", escapeString(t), weights[i]))
		}
	}
	return strings.Join(tokens, " || ")
}

// 刷新单个post
func flushPostTokens(postId uint64, title, content string) error {
	return db.Model(&Post{}).Where("id = ?", postId).
		Update("tokens", gorm.Expr(geneTokenString(title, content))).Error
}

// 将现有的未生成token的post刷新
func FlushPostsTokens(all bool) error {
	var posts []Post
	db.Find(&posts)
	if all {
		for _, p := range posts {
			if err := db.Model(&Post{}).Where("id = ?", p.Id).
				Update("tokens", gorm.Expr(geneTokenString(p.Title, p.Content))).Error; err != nil {
				return err
			}
			fmt.Println(p.Title, "更新成功")
		}
	} else {
		for _, p := range posts {
			if p.Tokens == "" {
				if err := db.Model(&Post{}).Where("id = ?", p.Id).
					Update("tokens", gorm.Expr(geneTokenString(p.Title, p.Content))).Error; err != nil {
					return err
				}
			}
			fmt.Println(p.Title, "更新成功")
		}
	}
	return nil
}

// 刷新单个tag
func flushTagTokens(tagId uint64, content string) error {
	return db.Model(&Tag{}).Where("id = ?", tagId).
		Update("tokens", gorm.Expr(geneTokenString(content))).Error
}

// 刷新tag的token
func FlushTagsTokens(all bool) error {
	var tags []Tag
	db.Find(&tags)
	for _, p := range tags {
		if all {
			if err := db.Model(&Tag{}).Where("id = ?", p.Id).
				Update("tokens", gorm.Expr(geneTokenString(p.Name))).Error; err != nil {
				return err
			}
		} else {
			if p.Tokens == "" {
				if err := db.Model(&Tag{}).Where("id = ?", p.Id).
					Update("tokens", gorm.Expr(geneTokenString(p.Name))).Error; err != nil {
					return err
				}
			}
		}
		fmt.Println(p.Name, "更新成功")
	}
	return nil
}
