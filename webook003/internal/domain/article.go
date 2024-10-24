package domain

import "time"

type Article struct {
	Id      int64
	Title   string
	Content string
	Author  Author
	Status  ArticleStatus
	Ctime   time.Time
	Utime   time.Time
}

// Abstract 函数用于获取文章的摘要
func (a Article) Abstract() string {
	// 将文章内容转换为 rune 类型
	str := []rune(a.Content)
	// 只取部分作为摘要
	if len(str) > 128 {
		str = str[:128]
	}
	// 返回摘要
	return string(str)
}

type Author struct {
	Id   int64
	Name string
}

type ArticleStatus uint8

func (s ArticleStatus) ToUint8() uint8 {
	return uint8(s)
}

const (
	// ArticleStatusUnknown 这是一个未知状态
	ArticleStatusUnknown = iota

	// ArticleStatusUnpublished 未发表
	ArticleStatusUnpublished

	// ArticleStatusPublished 已发表
	ArticleStatusPublished

	// ArticleStatusPrivate 仅自己可见
	ArticleStatusPrivate
)
