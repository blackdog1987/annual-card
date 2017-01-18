package manager

import (
	"crypto/md5"
	"encoding/hex"
	"strings"

	"github.com/blackdog1987/annual-card/lib/data"
	"github.com/blackdog1987/annual-card/module"
)

// ManagerService .
type ManagerService struct{}

func init() {
	module.Manager = &ManagerService{}
}

// Login .
// 登陆
func (s *ManagerService) Login(phone, passwd string) (info *data.Manager, ok bool) {
	info = &data.Manager{}
	if has, _ := data.Db.Where("phone = ?", phone).Get(info); !has {
		return
	}
	passwd = strings.ToUpper(passwd)
	md5Ctx := md5.New()

	md5Ctx.Write([]byte(passwd))
	cipherStr := md5Ctx.Sum(nil)
	passwd = hex.EncodeToString(cipherStr)
	if info.Passwd != passwd {
		return
	}
	return info, true
}
