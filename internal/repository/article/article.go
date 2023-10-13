package article

import (
	"context"
	"geekgo-webook/internal/domain"

	"geekgo-webook/internal/repository/dao/article"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
}

// 在repo层面同步数据 dao层区分为AuthorDao readerDao 制作库线上库
// service不做区别 统一调用repo的Sync方法

type CachedArticleRepository struct {
	dao       article.ArticleDAO
	authordao article.ArticleAuthorDAO
	readerdao article.ArticleReaderDAO

	// repository层面操作事务 需要组合gorm.DB 来开启事务 这样会耦合dao操作的东西
	db *gorm.DB
}

func NewArticleRepositoryV1(dao article.ArticleDAO, authordao article.ArticleAuthorDAO,
	readerdao article.ArticleReaderDAO) ArticleRepository {
	return &CachedArticleRepository{
		dao:       dao,
		authordao: authordao,
		readerdao: readerdao,
	}
}

func NewArticleRepository(dao article.ArticleDAO) ArticleRepository {
	return &CachedArticleRepository{
		dao: dao,
	}
}

// SyncV2尝试在repository解决事务 这需要在CachedArticleRepository里面增加 db *gorm.DB
// 确保保存到制作库线上库同时成功 同时失败
// 如果执行事务中间panic掉了，既没有回滚也没有提交，事务会一直挂在数据库，最好 defer.Rollback()
func (c CachedArticleRepository) SyncV2(ctx context.Context, art domain.Article) (int64, error) {
	// 开启事务
	tx := c.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	defer tx.Rollback() // 如果提交成功了 再defer回滚 会返回error
	// 利用 tx 来构建 DAO
	author := article.NewAuthorDAO(tx)
	reader := article.NewReaderDAO(tx)

	// 把正常的 制作库新增更新和线上库新增或更新的代码拷过来

	var (
		id  = art.Id
		err error
	)
	artn := c.domainToEntity(art)
	// 应该先保存到制作库，再保存到线上库
	if id > 0 {
		err = author.UpdateByID(ctx, art)
	} else {
		id, err = author.Insert(ctx, artn)
	}
	if err != nil {
		return id, err
	}
	// 操作线上库了，保存数据，同步过来
	// 考虑到，此时线上库可能有，可能没有，你要有一个 UPSERT 的写法
	// INSERT or UPDATE
	// 如果数据库有，那么就更新，不然就插入
	err = reader.Upsert(ctx, artn)
	if err != nil {
		return id, err
	}
	// 成功 提交
	tx.Commit()
	return id, err

}

// syncv1是非事务实现
func (c CachedArticleRepository) SyncV1(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	artn := c.domainToEntity(art)
	// 应该先保存到制作库，再保存到线上库
	if id > 0 {
		err = c.authordao.UpdateByID(ctx, art)
	} else {
		id, err = c.authordao.Insert(ctx, artn)
	}
	if err != nil {
		return id, err
	}
	// 操作线上库了，保存数据，同步过来
	// 考虑到，此时线上库可能有，可能没有，你要有一个 UPSERT 的写法
	// INSERT or UPDATE
	// 如果数据库有，那么就更新，不然就插入
	err = c.readerdao.Upsert(ctx, artn)
	return id, err
}

// sync这里调用dao.Sync dao.Sync处理事务 同库不同表
// dao 操作Article 和PublishArticle 两个表
func (c *CachedArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Sync(ctx, c.domainToEntity(art))
}

func (c CachedArticleRepository) Update(ctx context.Context, art domain.Article) error {
	return c.dao.UpdateById(ctx, c.domainToEntity(art))
}

func (c CachedArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	//调用dao.Create来创建文章
	id, err := c.dao.Insert(ctx, c.domainToEntity(art))
	return id, err
}

func (c CachedArticleRepository) domainToEntity(art domain.Article) article.Article {
	return article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}
