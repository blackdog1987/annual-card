package user

import (
	"github.com/molibei/annual-card/lib/data"
	"github.com/molibei/annual-card/module"
)

type UserService struct{}

func init() {
	module.User = &UserService{}
}

// GetInfoByOpenID .
func (s *UserService) GetInfoByOpenID(openID string) (has bool, info *data.User) {
	info = &data.User{}
	has, _ = data.Db.Where("wx_openid = ?", openID).Get(info)
	return
}

func (s *UserService) GetInfoByUID(uid int64) (has bool, info *data.User) {
	info = &data.User{}
	has, _ = data.Db.Where("uid = ?", uid).Get(info)
	return
}

//
func (s *UserService) Sync(info *data.User) (newinfo *data.User, err error) {
	_, err = data.Db.Insert(info)
	newinfo = info
	return
}

func (s *UserService) SpreadCount (uid int64) (total int64, err error) {
	return data.Db.Where("spread_uid in (select access_id from spread_access where relation_id = ?)", uid).Count(&data.User{})
}

func (s *UserService) SpreadSearch (uid int64, offset, length int) (items []data.User, err error) {
	items = make([]data.User, 0)
	err= data.Db.Where("spread_uid in (select access_id from spread_access where relation_id = ?)", uid).Limit(length,offset).OrderBy("points_contribute DESC").Find(&items)
	return
}