package service

import (
	"context"
	"geekgo-webook/internal/domain"
	"geekgo-webook/internal/repository/article"
	artrepomocks "geekgo-webook/internal/repository/article/mocks"
	"geekgo-webook/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_articeService_PublishV1(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (
			author article.ArticleAuthorRepository,
			reader article.ArticleReaderRepository,
		)
		// 返回初始化articleService所需要的 articlerepo readerrepo 实现接口Create Update Save
		art     domain.Article
		wantErr error
		wantId  int64
	}{
		{
			name: "新建发表成功",
			mock: func(ctrl *gomock.Controller) (
				article.ArticleAuthorRepository,
				article.ArticleReaderRepository,
			) {
				author := artrepomocks.NewMockArticleAuthorRepository(ctrl)
				author.EXPECT().Create(gomock.Any(),
					domain.Article{
						Title:   "我的标题",
						Content: "我的内容",
						Author: domain.Author{
							Id: 123,
						},
					}).Return(int64(1), nil) // 注意返回的类型 如果不做 int64(1)转换 测试会失败 missing call
				reader := artrepomocks.NewMockArticleReaderRepository(ctrl)
				reader.EXPECT().Save(gomock.Any(),
					domain.Article{
						Id:      1,
						Title:   "我的标题",
						Content: "我的内容",
						Author: domain.Author{
							Id: 123,
						},
					}).Return(int64(1), nil)
				return author, reader
			},
			art: domain.Article{
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantErr: nil,
			wantId:  1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			author, reader := tc.mock(ctrl)
			svc := NewArticleServiceV1(author, reader, &logger.NopLogger{})
			id, err := svc.PublishV1(context.Background(), tc.art)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantId, id)

		})
	}
}
