package demo

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"testing"
)

func TestInitViperReader(t *testing.T) {
	viper.SetConfigType("yaml")
	cfg := `
db.mysql:
  dsn: "root:root@tcp(localhost:13316)/webook"

redis:
  addr: "localhost:6379"
`
	// // NewReader returns a new Reader reading from b.
	//func NewReader(b []byte) *Reader
	// bytes.NewReader把字符串转成io.Reader
	err := viper.ReadConfig(bytes.NewReader([]byte(cfg)))
	// viper.ReadConfig(in io.Reader)
	if err != nil {
		panic(err)
	}
}

func TestViperRemote(t *testing.T) {
	err := viper.AddRemoteProvider("etcd3",
		"http://127.0.0.1:12379",
		"/webook")

	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	err = viper.WatchRemoteConfig()
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

func TestViperReadCMDArgument(t *testing.T) {

	//// String defines a string flag with specified name, default value, and usage string.
	//// The return value is the address of a string variable that stores the value of the flag.
	// Golang 的标准库提供了 flag 包来处理命令行参数 pflag是第三方包
	cfile := pflag.String("config",
		"../config/dev.yaml", "指定配置文件路径") // cfile 是 *string "../config/dev.yaml"
	// flag provided but not defined: -config
	// pflag.String提供了value默认值 要从命令行读取 设置了IDE的program argument
	pflag.Parse()
	//// Parse parses the command-line flags from os.Args[1:].  Must be called
	//// after all flags are defined and before flags are accessed by the program.
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 比较好的设计，它会在 in 里面告诉你变更前的数据，和变更后的数据
		// 更好的设计是，它会直接告诉你差异。
		fmt.Println(in.Name, in.Op)
		fmt.Println(viper.GetString("db.dsn"))
	})

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

}

func TestViper(t *testing.T) {

	// 获取当前进程的工作目录
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Current working directory:", dir)
	//Current working directory: C:\Users\Oasis\go\src\geekgo-webook\demo

	// 将相对文件路径转换成绝对文件路径
	absPath, err := filepath.Abs("test.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Absolute file path:", absPath)

	viper.SetConfigFile("../config/dev.yaml")
	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println("viper all keys: ", viper.AllKeys()) // [db.mysql.dsn redis.addr]
	dbstr := viper.GetString("db.mysql.dsn")
	fmt.Println("GetString db.mysql.dsn", dbstr)
	redisStr := viper.GetString("redis.addr")
	fmt.Println("GetString redis.addr", redisStr)

	type DBConfig struct {
		DSN string `yaml:"dsn"`
	}
	dbconfig := DBConfig{
		DSN: "root:root@tcp(localhost:3306)/mysql",
	}
	viper.UnmarshalKey("db.mysql", &dbconfig)

	type RedisConfig struct {
		Addr string `yaml:"addr"`
	}
	redisconfig := RedisConfig{
		Addr: "localhost:637x",
	}
	viper.UnmarshalKey("redis", &redisconfig)

	fmt.Printf("UnmarshalKey db.mysql %v\n", dbconfig)
	fmt.Printf("UnmarshalKey redis %v\n", redisconfig)

}
