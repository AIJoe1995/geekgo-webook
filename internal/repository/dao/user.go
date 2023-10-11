package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
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
}

type GORMUserDAO struct {
	db *gorm.DB
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
	Phone sql.NullString `gorm:"unique"`
	// 最大问题就是，你要解引用
	// 你要判空
	//Phone *string

	// 往这面加

	// 创建时间，毫秒数
	Ctime int64
	// 更新时间，毫秒数
	Utime int64
}
