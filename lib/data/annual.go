package data

// AnnualCardConf .
type AnnualCardConf struct {
	Name        string `json:"name" form:"name"`
	OriginPrice int64  `json:"origin_price" form:"origin_price"`
	SalePrice   int64  `json:"sale_price" form:"sale_price"`
	Usage       string `json:"usage" form:"usage"`
	Images      string `json:"images" form:"images"`
}

// CouponConf .
type CouponConf struct {
	Name        string `json:"name" form:"name"`
	OffsetPrice int64  `json:"offset_price" form:"offset_price"`
	SalePrice   int64  `json:"sale_price" form:"sale_price"`
	Usage       string `json:"usage" form:"usage"`
	Images      string `json:"images" form:"images"`
}

// AnnualCardPlan .
// 年卡计划
type AnnualCardPlan struct {
	CPID         int64  `json:"cp_id" xorm:"int(10) autoincr 'cp_id'"`
	Channel      string `json:"channel" xorm:"varchar(255) 'channel'" form:"channel"`
	ExpiredStart int64  `json:"expired_start" xorm:"int(10) 'expired_start'" form:"expired_start"`
	ExpiredStop  int64  `json:"expired_stop" xorm:"int(10) 'expired_stop'" form:"expired_stop"`
	CreateNum    int64  `json:"create_num" xorm:"int(10) 'create_num'" form:"create_num"`
	ActiveNum    int64  `json:"active_num" xorm:"int(10) 'active_num'"`
	Created      int64  `json:"created" xorm:"int(10) created 'created'"`
	Updated      int64  `json:"updated" xorm:"int(10) updated 'updated'"`
	IsDisabled   int    `json:"is_disabled" xorm:"tinyint(1) 'is_disabled'"`
}

// AnnualCard .
// 年卡基础表
type AnnualCard struct {
	CardID       int64  `json:"card_id" xorm:"int(10) autoincr 'card_id'" form:"card_id"`
	PlanID       int64  `json:"plan_id" xorm:"int(10) 'plan_id'" form:"plan_id"`
	CardName     string `json:"card_name" xorm:"varchar(255) 'card_name'" form:"card_name"`
	CardNO       string `json:"card_no" xorm:"varchar(20) unique 'card_no'" form:"card_no"`
	RelationUID  int64  `json:"relation_uid" xorm:"int(10) 'relation_uid'" form:"relation_uid"`
	BindHeadimg  string `json:"bind_headimg" xorm:"text 'bind_headimg'" form:"bind_headimg"`
	BindName     string `json:"bind_name" xorm:"varchar(32) 'bind_name'" form:"bind_name"`
	BindContact  string `json:"bind_contact" xorm:"varchar(12) 'bind_contact'" form:"bind_contact"`
	BindIDCard   string `json:"bind_idcard" xorm:"char(18) 'bind_idcard'" form:"bind_idcard"`
	CardPasswd   string `json:"card_passwd" xorm:"char(10) 'card_passwd'" form:"card_passwd"`
	ExpiredStart int64  `json:"expired_start" xorm:"int(10) 'expired_start'" form:"expired_start"`
	ExpiredStop  int64  `json:"expired_stop" xorm:"int(10) 'expired_stop'" form:"expired_stop"`
	IsActive     int    `json:"is_active" xorm:"char(1) 'is_active'" form:"is_active"`
	IsDelete     int `json:"is_delete" xorm:"char(1) 'is_delete'"`
	Created      int64  `json:"created" xorm:"int(10) created 'created'"`
	Updated      int64  `json:"updated" xorm:"int(10) updated 'updated'"`
}

type CardInfo struct {
	CardID      int64  `json:"card_id" xorm:"int(10) autoincr 'card_id'" form:"card_id"`
	CardName    string `json:"card_name" xorm:"varchar(255) 'card_name'" form:"card_name"`
	CardNO      string `json:"card_no" xorm:"varchar(20) unique 'card_no'" form:"card_no"`
	BindHeadimg string `json:"bind_headimg" xorm:"text 'bind_headimg'" form:"bind_headimg"`
	BindName    string `json:"bind_name" xorm:"varchar(32) 'bind_name'" form:"bind_name"`
	BindContact string `json:"bind_contact" xorm:"varchar(12) 'bind_contact'" form:"bind_contact"`
	BindIDCard  string `json:"bind_idcard" xorm:"char(18) 'bind_idcard'" form:"bind_idcard"`
}

func (CardInfo) TableName() string {
	return "annual_card"
}

type CardUsageInfo struct {
	CardID      int64  `json:"card_id" xorm:"int(10) autoincr 'card_id'" form:"card_id"`
	CardName    string `json:"card_name" xorm:"varchar(255) 'card_name'" form:"card_name"`
	CardNO      string `json:"card_no" xorm:"varchar(20) unique 'card_no'" form:"card_no"`
	BindHeadimg string `json:"bind_headimg" xorm:"text 'bind_headimg'" form:"bind_headimg"`
	BindName    string `json:"bind_name" xorm:"varchar(32) 'bind_name'" form:"bind_name"`
	BindContact string `json:"bind_contact" xorm:"varchar(12) 'bind_contact'" form:"bind_contact"`
	BindIDCard  string `json:"bind_idcard" xorm:"char(18) 'bind_idcard'" form:"bind_idcard"`
	UsageNum    int64 `json:"usage_num" xorm:"'usage_num'"`
}

func (CardUsageInfo) TableName() string {
	return "annual_card"
}

// AnnualCardUsageLog .
// 年卡使用表
type AnnualCardUsageLog struct {
	UsageID   int64 `json:"usage_id" xorm:"int(10) autoincr 'usage_id'"`
	CardID    int64 `json:"card_id" xorm:"int(10) 'card_id'"`
	MchID     int64 `json:"mch_id" xorm:"int(10) 'mch_id'"`
	StoreID   int64 `json:"store_id" xorm:"int(10) 'store_id'"`
	UsageTime int64 `json:"usage_time" xorm:"int(10) created 'usage_time'"`
	CardInfo `json:"card" xorm:"extends"`
	StoreName  string `json:"store_name" xorm:"varchar(255) 'store_name'"`
}

func (AnnualCardUsageLog) TableName() string {
	return "annual_card_usage_log"
}

type CardUsageLogInsert struct {
	UsageID   int64 `json:"usage_id" xorm:"int(10) autoincr 'usage_id'"`
	CardID    int64 `json:"card_id" xorm:"int(10) 'card_id'"`
	MchID     int64 `json:"mch_id" xorm:"int(10) 'mch_id'"`
	StoreID   int64 `json:"store_id" xorm:"int(10) 'store_id'"`
	UsageTime int64 `json:"usage_time" xorm:"int(10) created 'usage_time'"`
}

func (CardUsageLogInsert) TableName() string {
	return "annual_card_usage_log"
}

// Merchant .
// 商户基础表
type Merchant struct {
	MchID     int64  `json:"mch_id" xorm:"int(10) autoincr 'mch_id'" form:"id"`
	MchName   string `json:"mch_name" xorm:"varchar(255) 'mch_name'" form:"name"`
	Value     string `json:"value" xorm:"varchar(10) 'value'" form:"value"`
	Consume   string `json:"consume" xorm:"text 'consume'" form:"consume"`
	Usage     string `json:"usage" xorm:"text 'usage'" form:"usage"`
	Contact   string `json:"contact" xorm:"varchar(255) 'contact'" form:"contact"`
	Address   string `json:"address" xorm:"text 'address'" form:"address"`
	Introduce string `json:"introduce" xorm:"text 'introduce'" form:"introduce"`
	Cover     string `json:"cover" xorm:"text 'cover'" form:"cover"`
	Imgs      string `json:"imgs" xorm:"text 'imgs'" form:"images"`
	Created   int64  `json:"created" xorm:"int(10) created 'created'"`
	Updated   int64  `json:"updated" xorm:"int(10) updated 'updated'"`
	State     int    `json:"state" xorm:"tinying(1) 'state'"`
}

// MerchantAccount
// 商户账户表
type MerchantAccount struct {
	MchID   int64 `json:"mch_id" xorm:"int(10) autoincr 'mch_id'" form:"id"`
	Name    string `json:"name" xorm:"varchar(255) 'name'" form:"name"`
	Account string `json:"account" xorm:"varchar(255) 'account'" form:"account"`
	Passwd  string `json:"-" xorm:"varchar(255) 'passwd'" form:"passwd"`
	Created int64  `json:"created" xorm:"int(10) created 'created'"`
	Updated int64  `json:"updated" xorm:"int(10) updated 'updated'"`
	State   int    `json:"state" xorm:"tinying(1) 'state'" form:"state"`
}

// MerchantStore
// 商户分店表
type MerchantStore struct {
	StoreID   int64 `json:"store_id" xorm:"int(10) autoincr 'store_id'" form:"id"`
	MchId     int64 `json:"mch_id" xorm:"int(10) 'mch_id'" form:"mch_id"`
	StoreName string `json:"store_name" xorm:"varchar(255) 'store_name'" form:"store_name"`
	Account   string `json:"account" xorm:"varchar(255) 'account'" form:"account"`
	Passwd    string `json:"-" xorm:"varchar(255) 'passwd'" form:"passwd"`
	Created   int64  `json:"created" xorm:"int(10) created 'created'"`
	Updated   int64  `json:"updated" xorm:"int(10) updated 'updated'"`
	State     int    `json:"state" xorm:"tinying(1) 'state'"`
}
