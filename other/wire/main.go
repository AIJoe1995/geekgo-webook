package wire

import (
	"fmt"
	"geekgo-webook/other/wire/repository"
	"geekgo-webook/other/wire/repository/dao"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("dsn"))
	if err != nil {
		panic(err)
	}
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	fmt.Println(repo)

	InitRepository()
}
