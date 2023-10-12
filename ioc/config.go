package ioc

import (
	"context"
	"github.com/spf13/viper"
)

// 提供统一配置接口
type Configer interface {
	GetString(ctx context.Context, key string) (string, error)
	MustGetString(ctx context.Context, key string) string
	GetStringOrDefault(ctc context.Context, key string, def string) string
	//Unmarshal()

}

// 初始化读取配置模块 适配器

type ViperConfigerAdapter struct {
	v *viper.Viper
}

func (m *myConfiger) GetString(ctx context.Context, key string) (string, error) {
	//TODO implement me
	panic("implement me")
}

type myConfiger struct {
}

func (m *myConfiger) MustGetString(ctx context.Context, key string) string {
	str, err := m.GetString(ctx, key)
	if err != nil {
		panic(err)
	}
	return str
}

func (m *myConfiger) GetStringOrDefault(ctx context.Context, key string, def string) string {
	str, err := m.GetString(ctx, key)
	if err != nil {
		return def
	}
	return str
}
