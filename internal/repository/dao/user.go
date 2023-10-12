package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	mysql2 "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync/atomic"
	"time"
	"unsafe"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(ctx context.Context, u User) error
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByWechat(ctx context.Context, openID string) (User, error)
}

type GORMUserDAO struct {
	db *gorm.DB
}

type GORMUserDAOV2 struct {
	db *gorm.DB

	p DBProvider
}

func (G GORMUserDAOV2) Insert(ctx context.Context, u User) error {
	return nil
}

func (G GORMUserDAOV2) FindByPhone(ctx context.Context, phone string) (User, error) {
	return User{}, nil
}

func (G GORMUserDAOV2) FindById(ctx context.Context, id int64) (User, error) {
	return User{}, nil
}

func (G GORMUserDAOV2) FindByEmail(ctx context.Context, email string) (User, error) {
	return User{}, nil
}

func (G GORMUserDAOV2) FindByWechat(ctx context.Context, openID string) (User, error) {
	return User{}, nil
}

func NewUserDAOV2(p DBProvider) UserDAO {
	return &GORMUserDAOV2{
		p: p,
	}
}

type DBProvider func() *gorm.DB

// 需要有可以监听db配置变化 创建新db连接的初始化方法  但这样不够优雅， 最好在初始化db的地方监听配置变化，
// 修改方式： 把监听配置变化的代码移动到初始化db的位置ioc/db.go, 返回一个db, 在GORMUserDAO里组合这个方法
// 然后NewUserDAO里传入这个方法，作为userdao的属性，在使用db的时候调用这个方法，使用这个方法返回的db
func NewUserDAOV1(db *gorm.DB) UserDAO {
	res := &GORMUserDAO{
		db: db,
	}
	viper.OnConfigChange(func(in fsnotify.Event) { // 应该只有mysql变化的时候才去重新创建db
		// viper的onconfigchange 是什么机制 如果其他的配置变了 也存在变化 也会重新初始化mysql吗
		dsn := viper.GetString("db.mysql.dsn")
		db, err := gorm.Open(mysql2.Open(dsn))
		if err != nil {
			panic(err)
		}
		pt := unsafe.Pointer(&res.db)
		//func StorePointer(addr *unsafe.Pointer, val unsafe.Pointer)
		atomic.StorePointer(&pt, unsafe.Pointer(&db)) // atomic.StorePointer要求传入unsafe.Pointer类型

	})
	return res
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

func (dao *GORMUserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	// // WithContext change current instance db's context to ctx
	//func (db *DB) WithContext(ctx context.Context) *DB {
	//	return db.Session(&Session{Context: ctx})
	//}
	// // Create inserts value, returning the inserted data's primary key in value's id
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			// 邮箱冲突
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (dao *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) FindByWechat(ctx context.Context, openID string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("wechat_open_id = ?", openID).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return u, err
}

// dao层返回的是dao层的User结构体
func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err == gorm.ErrRecordNotFound {
		return u, ErrUserNotFound
	}
	return u, err
}

// User 直接对应数据库表结构
//type User struct {
//	Id int64 `gorm:"primaryKey,autoIncrement"`
//	// 全部用户唯一
//	Email    string `gorm:"unique"`
//	Password string
//
//	// 往这面加
//
//	// 创建时间，毫秒数
//	Ctime int64
//	// 更新时间，毫秒数
//	Utime int64
//}

type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 全部用户唯一
	Email    sql.NullString `gorm:"unique"`
	Password string

	// 唯一索引允许有多个空值
	// 但是不能有多个 ""
	Phone         sql.NullString `gorm:"unique"`
	WechatUnionID sql.NullString
	WechatOpenID  sql.NullString `gorm:"unique"`
	// 最大问题就是，你要解引用
	// 你要判空
	//Phone *string

	// 往这面加

	// 创建时间，毫秒数
	Ctime int64
	// 更新时间，毫秒数
	Utime int64
}
