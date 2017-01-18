package data

import (
	"encoding/json"
	"os"
	// .
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var (
	// Db .
	Db *xorm.Engine
	// BaseConf .
	BaseConf *BaseConfig
)

func init() {
	BaseConf = &BaseConfig{}
}

// LoadBaseConf .
// 从文件中读取配置
func LoadBaseConf(filepath string) (err error) {
	fi, err := os.Open(filepath)
	if err != nil {
		return
	}
	err = json.NewDecoder(fi).Decode(BaseConf)
	return
}

// ConnectMySQL .
// 连接数据库
func ConnectMySQL() (err error) {
	Db, err = xorm.NewEngine("mysql", BaseConf.MySQL.Dsn)
	if err != nil {
		return
	}
	if err = Db.Ping(); err != nil {
		return
	}
	Db.SetMaxOpenConns(BaseConf.MySQL.MaxConn)
	Db.SetMaxIdleConns(BaseConf.MySQL.MaxIdle)
	Db.ShowSQL(BaseConf.MySQL.ShowSQL)
	return
}
