package ioc

import (
	"geekgo-webook/internal/repository/dao"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync/atomic"
	"unsafe"
)

// 监听配置变化返回新db的函数
func InitDBV2() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg = Config{
		DSN: "root:root@tcp(localhost:13316)/webook_default",
	}
	err := viper.UnmarshalKey("db.mysql", &cfg)
	db, err := gorm.Open(mysql.Open(cfg.DSN))
	if err != nil {
		// 我只会在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 一旦初始化过程出错，应用就不要启动了
		panic(err)
	}

	// config变化时，这里返回新的db 调用UserDao接口的方法，会调用这个函数 拿到这个函数返回的db.
	dao.NewUserDAOV2(
		func() *gorm.DB {
			viper.OnConfigChange(func(in fsnotify.Event) {
				//oldDB := db
				err := viper.UnmarshalKey("db.mysql", &cfg)
				if err != nil {
					panic(err)
				}
				db, err = gorm.Open(mysql.Open(cfg.DSN))
				pt := unsafe.Pointer(&db)
				atomic.StorePointer(&pt, unsafe.Pointer(&db))
				//oldDB.Close()
			})
			//要用原子操作
			return db
		})

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:1234@tcp(localhost:3306)/webook"))
	if err != nil {
		// 我只会在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 一旦初始化过程出错，应用就不要启动了
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
