package service

import (
	"context"
	"geekgo-webook/internal/domain"
	"geekgo-webook/internal/repository/article"
	"time"

	//"geekgo-webook/internal/web" // service不应该引用web? 会造成循环依赖？
	"geekgo-webook/pkg/logger"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error) // 需要实现一个Article的domain
	PublishV1(ctx context.Context, art domain.Article) (int64, error)
}

func NewArticleService(repo article.ArticleRepository) ArticleService {
	return &articeService{
		repo: repo,
	}
}

//type articleServiceV1 struct {
// 没有实现接口 会报错 暂时使用articleService来初始化带有author和reader的版本
//	author article.ArticleAuthorRepository
//	reader article.ArticleReaderRepository
//	l      logger.LoggerV1
//}

func NewArticleServiceV1(author_repo article.ArticleAuthorRepository, reader_repo article.ArticleReaderRepository, l logger.LoggerV1) ArticleService {
	return &articeService{
		author: author_repo,
		reader: reader_repo,
		logger: l,
	}

}

type articeService struct {
	repo article.ArticleRepository

	// v1
	author article.ArticleAuthorRepository
	reader article.ArticleReaderRepository
	logger logger.LoggerV1
}

// 测试web的Publish方法时 会mock ArticleService 提供输入输出 所以暂时不需要implement
func (a *articeService) Publish(ctx context.Context, article domain.Article) (int64, error) {
	// TDD service的publish方法
	// service会调用repo 存在线上库和制作库 在repo中区分
	// 制作库
	//id, err := a.repo.Create(ctx, art)
	//// 线上库呢？
	//a.repo.SyncToLiveDB(ctx, art)
	panic("implement me")
}

func (a *articeService) PublishV1(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	// 注入repository层的 AuthorArticleRepository ReaderArticleRepository
	// a.author 操作制作库
	if art.Id > 0 {
		err = a.author.Update(ctx, art)
	} else {
		id, err = a.author.Create(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	// 确保制作库和线上库id相同
	art.Id = id
	// a.reader 操作制作库 可能失败 设计重试
	//更优雅的做法是不在service里循环重试 而是在reader上做一个重试的装饰器 调用带有重试装饰器的save方法
	//id, err = a.reader.Save(ctx, art)

	for i := 0; i < 3; i++ {
		time.Sleep(time.Second * time.Duration(i))
		id, err = a.reader.Save(ctx, art)
		if err == nil {
			break
		}
		a.logger.Error("部分失败，保存到线上库失败",
			logger.Int64("art_id", art.Id),
			logger.Error(err))
	}
	if err != nil {
		a.logger.Error("部分失败，重试彻底失败",
			logger.Int64("art_id", art.Id),
			logger.Error(err))
		// 接入你的告警系统，手工处理一下
		// 走异步，我直接保存到本地文件
		// 走 Canal
		// 打 MQ
	}

	if err != nil {
		a.logger.Error("部分失败，保存到线上库失败", logger.Error(err))

	}
	return id, err

}

func (a *articeService) Save(ctx context.Context, art domain.Article) (int64, error) {
	// 区分新建和修改
	if art.Id > 0 {
		return art.Id, a.repo.Update(ctx, art)
	}
	return a.repo.Create(ctx, art)

}
