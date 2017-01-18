package data

// OrderLog .
// 支付订单表
type OrderLog struct {
	OrderID       int64  `json:"order_id" xorm:"int(10) autoincr 'order_id'"`
	OrderNO       string `json:"order_no" xorm:"char(20) 'order_no'"`
	UID           int64  `json:"uid" xorm:"int(10) 'uid'"`
	GoodsName     string `json:"goods_name" xorm:"varchar(255) 'goods_name'"`
	GoodsNum      int64 `json:"goods_num" xorm:"int(10) 'goods_num'"`
	TransactionID string `json:"transaction_id" xorm:"varchar(64) 'transaction_id'"`
	Price 		  int64  `json:"price" xorm:"int(10) 'price'"`
	Total         int64  `json:"total" xorm:"int(10) 'total'"`
	Category      string `json:"category" xorm:"varchar(10) 'category'"`
	Points        int64  `json:"points" xorm:"int(10) 'points'"`
	PointsPrice   int64  `json:"points_price" xorm:"int(10) 'points_price'"`
	IsPay         int    `json:"is_pay" xorm:"tinyint(1) 'is_pay'"`
	IsCoupon	  int64   `json:"is_coupon" xorm:"int(10) 'is_coupon'"`
	CouponPrice   int64  `json:"coupon_price" xorm:"int(10) 'coupon_price'"`
	Created       int64  `json:"created" xorm:"int(10) created 'created'"`
	Updated       int64  `json:"updated" xorm:"int(10) updated 'updated'"`
}

// CouponLog .
// 优惠券记录
type CouponLog struct {
	LogID       int64 `json:"log_id" xorm:"int(10) autoincr 'log_id'" form:"log_id"`
	UID         int64 `json:"uid" xorm:"int(10) 'uid'" form:"uid"`
	OffsetPrice int64 `json:"offset_price" xorm:"int(10) 'offset_price'"`
	IsUsage     int   `json:"is_usage" xorm:"tinyint(1) 'is_usage'"`
	UsageID     int64 `json:"usage_id" xorm:"int(10) 'usage_id'"`
	Created     int64 `json:"created" xorm:"int(10) created"`
	Updated     int64 `json:"updated" xorm:"int(10) updated"`
}

func (CouponLog) TableName() string {
	return "coupon_logs"
}

// Config .
// 配置表
type Config struct {
	Key   string `json:"key" xorm:"char(10) 'key'"`
	Value string `json:"value" xorm:"text 'value'"`
	Name  string `json:"name" xorm:"varchar(64) 'name'"`
}

// Banner .
// Banner配置表
type Banner struct {
	BannerID int64  `json:"banner_id" xorm:"int(10) autoincr 'banner_id'" form:"id"`
	Name     string `json:"name" xorm:"varchar(64) 'name'" form:"name"`
	Img      string `json:"img" xorm:"text 'img'" form:"img"`
	Link     string `json:"link" xorm:"text 'link'" form:"link"`
	Ordid    int64  `json:"ordid" xorm:"int(10) 'ordid'" form:"ordid"`
	Created  int64  `json:"created" xorm:"int(10) created 'created'"`
}
