package data

import (
	"time"
)

// User .
type User struct {
	UID           int64  `json:"uid" xorm:"int(11) autoincr pk 'uid'"`           // 用户 ID
	SpreadUID     int64  `json:"spread_uid" xorm:"int(11) 'spread_uid'"`         // 推广ID
	Nickname      string `json:"nickname" xorm:"varchar(32) 'nickname'"`         // 昵称
	Sex           int    `json:"sex"`                                            // 性别
	Realname      string `json:"realname"`                                       // 真实姓名
	Phone         string `json:"phone"`                                          // 手机号
	HeadImageURL  string `json:"headimgurl" xorm:"text 'headimgurl'"`            // 头像
	Country       string `json:"country"`                                        // 国家
	Province      string `json:"province"`                                       // 省份
	City          string `json:"city"`                                           // 城市
	Area          string `json:"area"`                                           // 区县
	Address       string `json:"address"`                                        // 详细地址
	PointsBalance int64  `json:"points_balance" xorm:"int(10) 'points_balance'"` // 积分余额
	PointsEarning int64  `json:"points_earning" xorm:"int(10) 'points_earning'"` // 积分消费合计
	PointsExpend  int64  `json:"points_expend" xorm:"int(10) 'points_expend'"`   // 积分支出合计
	PointsContribute int64 `json:"points_contribute" xorm:"int(10) 'points_contribute'"` // 积分贡献
	WxOpenID      string `json:"wx_openid" xorm:"varchar(64) 'wx_openid'"`       // 微信开放平台ID
	WxUnionID     string `json:"wx_unionid" xorm:"varchar(64) 'wx_unionid'"`     // 微信unionID
	Created       int64  `json:"created" xorm:"int(10) created"`
	Updated       int64  `json:"updated" xorm:"int(10) updated"`
}

// PointsLog .
// 积分记录表
type PointsLog struct {
	LogID         int64  `json:"log_id" xorm:"int(10) autoincr 'log_id'"`
	UID           int64  `json:"uid" xorm:"int(10) 'uid'"`
	RelationUID   int64  `json:"relation_uid" xorm:"char(10) 'relation_uid'"`
	RelationLogID int64  `json:"relation_log_id" xorm:"int(10) 'relation_log_id'"`
	FriendlyIntro string `json:"friendly_intro" xorm:"text 'friendly_intro'"`
	Total         int64  `json:"total" xorm:"int(10) 'total'"`
	Type          int    `json:"type" xorm:"tinyint(1) 'type'"`
	Created       int64  `json:"created" xorm:"int(10) created 'created'"`
}

// PointsEarningLog .
// 积分收入表
type PointsEarningLog struct {
	EarningID  int64 `json:"earning_id" xorm:"int(10) 'earning_id'"`
	UID        int64 `json:"uid" xorm:"int(10) 'uid'"`
	RelationID int64 `json:"relation_id" xorm:"int(10) 'relation_id'"`
	Total      int64 `json:"total" xorm:"int(10) 'total'"`
	Type       int   `json:"type" xorm:"tinyint(1) 'type'"`
	Created    int64 `json:"created" xorm:"int(10) created 'created'"`
}

// PointsExpendLog .
// 积分消费表
type PointsExpendLog struct {
	ExpendID   int64     `json:"earning_id" xorm:"int(10) 'earning_id'"`
	UID        int64     `json:"uid" xorm:"int(10) 'uid'"`
	RelationID int64     `json:"relation_id" xorm:"int(10) 'relation_id'"`
	Total      int64     `json:"total" xorm:"int(10) 'total'"`
	Type       int       `json:"type" xorm:"tinyint(1) 'type'"`
	Created    time.Time `json:"created" xorm:"int(10) created 'created'"`
}

// Manager .
// 管理员用户表
type Manager struct {
	ManagerID  int64     `json:"manager_id" xorm:"int(10) unsigned autoincr 'manager_id'"` // 管理ID
	Name       string    `json:"name" xorm:"varchar(64) 'name'"`
	Phone      string    `json:"phone" xorm:"char(11) 'phone'"`
	Email      string    `json:"email" xorm:"varchar(255) 'email'"`
	Passwd     string    `json:"-" xorm:"char(40) 'passwd'"`
	Created    time.Time `json:"created" xorm:"int(10) created 'created'"`
	Updated    time.Time `json:"updated" xorm:"int(10) updated 'updated'"`
	IsDisabled int       `json:"is_disabled" xorm:"tinyint(1) 'is_disabled'"`
}
