package data

// SpreadAccess
// 推广
type SpreadAccess struct {
	AccessID   int64 `json:"access_id" xorm:"int(10) autoincr 'access_id'"` // 推广记录ID
	RelationID int64 `json:"relation_id" xorm:"int(10) 'relation_id'"`      // 关联的推广来源ID
	Category   string `json:"category" xorm:"char(10) 'category'"`          // 推广的来源分类
	QRCode     string `json:"qrcode" xorm:"text 'qrcode'"`                  // 二维码
}

type SpreadPlan struct {
	SPID           int64  `json:"sp_id" xorm:"int(10) autoincr 'sp_id'"`
	Name           string `json:"name" xorm:"varchar(255) 'name'"`
	Channel        string `json:"channel" xorm:"varchar(64)"`
	Contact        string `json:"contact" xorm:"varchar(32)"`
	RegCommission  int64  `json:"reg_commission" xorm:"int(10) 'reg_commission'"`
	SaleCommission int64  `json:"sale_commission" xorm:"int(10) 'sale_commission'"`
	Created        int64  `json:"created" xorm:"int(10) created"`
	Updated        int64  `json:"updated" xorm:"int(10) updated"`
	IsDisabled     int    `json:"is_disabled" xorm:"tinyint(1) 'is_disabled'"`
}

type SpreadLogs struct {
	LogID      int64 `json:"log_id" xorm:"int(10) autoincr 'log_id'"`  // 明细ID
	PlanID     int64 `json:"plan_id" xorm:"int(10) 'plan_id'"`          // 计划ID
	UID        int64 `json:"uid" xorm:"int(10) 'uid'"`                 // 用户ID
	RelationID int64 `json:"-" xorm:"int(10) 'relation_id'"`           // 关联ID 推广时为用户ID 销售的为订单ID
	Category   int   `json:"category" xorm:"tinyint(1) 'category'"`    // 分类
	Commission int64 `json:"commission" xorm:"int(10) 'commission'"`   // 返佣
	OrderTotal int64 `json:"order_total" xorm:"int(10) 'order_total'"` // 订单金额
	Created    int64 `json:"created" xorm:"int(10) created"`           // 订单创建时间
}

type SpreadLogItem struct {
	LogID      int64    `json:"log_id" xorm:"int(10) autoincr 'log_id'"`  // 明细ID
	PlanID     int64    `json:"plan_id" xorm:"int(10) 'plan_id'"`         // 计划ID
	UID        int64    `json:"uid" xorm:"int(10) 'uid'"`                 // 用户ID
	Consumer   Consumer `json:"consumer" xorm:"extends"`                  // 消费者
	Category   int      `json:"category" xorm:"tinyint(1) 'category'"`    // 分类
	Commission int64    `json:"commission" xorm:"int(10) 'commission'"`   // 返佣
	OrderTotal int64    `json:"order_total" xorm:"int(10) 'order_total'"` // 订单金额
	Created    int64    `json:"created" xorm:"int(10) created"`           // 订单创建时间
}

type Consumer struct {
	Nickname string `json:"nickname" xorm:"nickname"`
	RealName string `json:"realname" xorm:"realname"`
	Phone    string `json:"phone" xorm:"phone"`
}
