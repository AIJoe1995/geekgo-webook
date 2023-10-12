package demo

import (
	"errors"
	"go.uber.org/zap"
	"testing"
)

func TestInitLoggerWithInfo(t *testing.T) {
	logger, err := zap.NewDevelopment(zap.AddCaller())
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	zap.L().Info("打印文件方法名代码行数")
}

func TestInitLogger(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	zap.L().Info("hello，你搞好了")
	type Demo struct {
		Name string `json:"name"`
	}
	zap.L().Info("这是实验参数",
		zap.Error(errors.New("这是一个 error")),
		zap.Int64("id", 123),
		zap.Any("一个结构体", Demo{Name: "hello"}))

}
