package conf

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/blackdog1987/annual-card/lib/data"
	"github.com/blackdog1987/annual-card/module"
)

type ConfService struct {
	Key   string `xorm:"char(10) 'key'"`
	Value string `xorm:"text 'value'"`
	Name  string `xorm:"varchar(64) 'name'"`
}

func (ConfService) TableName() string {
	return "config"
}

func init() {
	module.Conf = &ConfService{}
}

func (s *ConfService) SetObject(key string, val interface{}) (ok bool) {
	v, err := json.Marshal(val)
	if err != nil {
		fmt.Println(err)
		return
	}
	return s.set(key, string(v))
}

func (s *ConfService) GetObject(key string, val interface{}) (has bool) {
	has, v := s.get(key)
	if has {
		json.Unmarshal([]byte(v), val)
	}
	return
}

//
func (s *ConfService) Set(key string, val interface{}) (ok bool) {
	return s.set(key, val)
}

//
func (s *ConfService) Get(key string) (has bool, val string) {
	return s.get(key)
}

// GetInt64 .
func (s *ConfService) GetInt64(key string) (has bool, val int64) {
	has, v := s.get(key)
	if has {
		val, _ = strconv.ParseInt(v, 10, 64)
	}
	return
}

func (s *ConfService) GetFloat64 (key string) (has bool, val float64) {
	has, v := s.get(key)
	if has {
		val, _ = strconv.ParseFloat(v, 10)
	}
	return
}

func (s *ConfService) set(key string, val interface{}) (ok bool) {
	key = strings.ToUpper(key)
	_, err := data.Db.Exec("INSERT INTO `config` VALUES(?, ?, ?) ON DUPLICATE KEY UPDATE `value`= ?", key, val, key, val)
	fmt.Println(err)
	return err == nil
}

func (s *ConfService) get(key string) (has bool, val string) {
	key = strings.ToUpper(key)
	var c = &ConfService{}
	has, _ = data.Db.Where("`key` = ?", key).Get(c)
	return has, c.Value
}
