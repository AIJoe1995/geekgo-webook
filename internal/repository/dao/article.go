package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error // 这里要是用dao层定义的article
}

type GORMArticleDAO struct {
	db *gorm.DB
}

func NewGORMArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}

func (G GORMArticleDAO) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.Utime = now
	res := G.db.WithContext(ctx).Model(&art).
		Where("id=? AND author_id = ?", art.Id, art.AuthorId).
		Updates(map[string]any{
			"title":   art.Title,
			"content": art.Content,
			"utime":   art.Utime,
		})

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		//dangerousDBOp.Count(1)
		// 补充一点日志
		return fmt.Errorf("更新失败，可能是创作者非法 id %d, author_id %d",
			art.Id, art.AuthorId)
	}
	return res.Error
}

func (G GORMArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := G.db.WithContext(ctx).Create(&art).Error
	return art.Id, err

}

type Article struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 长度 1024
	Title   string `gorm:"type=varchar(1024)"`
	Content string `gorm:"type=BLOB"`
	// 如何设计索引
	// 在帖子这里，什么样查询场景？
	// 对于创作者来说，是不是看草稿箱，看到所有自己的文章？
	// SELECT * FROM articles WHERE author_id = 123 ORDER BY `ctime` DESC;
	// 产品经理告诉你，要按照创建时间的倒序排序
	// 单独查询某一篇 SELECT * FROM articles WHERE id = 1
	// 在查询接口，我们深入讨论这个问题
	// - 在 author_id 和 ctime 上创建联合索引
	// - 在 author_id 上创建索引

	// 学学 Explain 命令

	// 在 author_id 上创建索引
	AuthorId int64 `gorm:"index"`
	//AuthorId int64 `gorm:"index=aid_ctime"`
	//Ctime    int64 `gorm:"index=aid_ctime"`
	Ctime int64
	Utime int64
}