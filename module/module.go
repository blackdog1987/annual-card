package module

import "github.com/blackdog1987/annual-card/lib/data"

// ManagerModule .
type ManagerModule interface {
	Login(phone, passwd string) (info *data.Manager, ok bool)
}

// TokenModule 令牌.
type TokenModule interface {
	Encode(tk data.Token, key string) (token string, err error)
	Decode(secret, key string) (tk data.Token, ok bool)
}

// SpreadModule  .
type SpreadModule interface {
	// 获取单个计划
	Plan(planID int64) (plan *data.SpreadPlan, has bool)
	// 更新
	UpdatePlan(plan *data.SpreadPlan) (ok bool)
	// 重复检测
	IsPlanNameRepeat(id int64, name string) (has bool)
	// 获取二维码
	QRCode(rid int64, category string) (qrcode string, err error)
	// 处理提成
	Commission(rid, total, uid, orderId int64) (ok bool)
	PlansCount(state int) (count int64, err error)
	Plans(state int, offset, length int) ([]data.SpreadPlan, error)
	LogsCount(planID, start, end int64) (count int64, err error)
	Logs(planID, start, end int64, offset, length int) (items []data.SpreadLogItem, err error)
}

type ConfModule interface {
	SetObject(key string, val interface{}) (ok bool)
	GetObject(key string, val interface{}) (has bool)
	Set(key string, val interface{}) (ok bool)
	Get(key string) (has bool, val string)
	GetInt64(key string) (has bool, val int64)
	GetFloat64(key string) (has bool, val float64)
}

type OrderModule interface {
	Coupon(uid int64) (info *data.CouponLog, has bool)
	Count(start, end int64) (count int64, err error)
	Search(start, end int64, offset, length int) (items []data.OrderLog, err error)
}

type AnnualCardModule interface {
	CardByID(cardID int64) (info *data.AnnualCard, has bool)
	IsExistIDCard(cardId int64, idCard string) (info *data.AnnualCard, has bool)
	CardByCardNo(cardNo string) (info *data.AnnualCard, has bool)
	Usages(uid, mchID int64) (items []data.CardUsageInfo, err error)
	UsagesByCardID(cardID int64) (info data.CardUsageInfo, has bool)
	CountBind(uid int64, active int) (total int64, err error)
	SearchBind(uid int64, active, offset, length int) (items []data.AnnualCard, err error)
	Count(isActive int, planID, start, end int64, name, idcard string) (count int64, err error)
	Search(isActive int, planID, start, end int64, name, idcard string, offset, length int) (items []data.AnnualCard, err error)
	PlanCount() (count int64, err error)
	PlanSearch(offset, length int) (items []data.AnnualCardPlan, err error)
}

type MerchantModule interface {
	Get(mid int64) (has bool, info data.Merchant)
	IsMchNameRepeat(mchID int64, name string) (has bool)
	Count(keyword string) (count int64, err error)
	Search(keyword string, offset, length int) (items []data.Merchant, err error)
}

// WechatModule .
type WechatModule interface {
	GetAccessToken() (token string, err error)
	GetUserInfo(data.EventMessage) data.UserOauth
	OauthInfo(code string) (info *data.User, err error)
	Subscribe(data.EventMessage) data.EventMessage
	UnSubscribe(data.EventMessage)
	GenerateJsSign(url string) *data.JsSignPackage
	GetJsAPITicket() (ticket string, err error)
	GenerateNonceStr(length int) (nonce string)
	Download(mediaID, filename string) bool
	GenerateSnsInfoURI(redirect string) (uri string)
}
type WxPayModule interface {
	UnifiedOrder(msg, clientIP string, orderID, openID string, total int64) (map[string]string, error)
	SendRedPack(openID string, total int64, orderId string) (ok bool, err error)
}

type UserModule interface {
	SpreadCount(uid int64) (total int64, err error)
	SpreadSearch(uid int64, offset, length int) (items []data.User, err error)
	GetInfoByOpenID(openID string) (has bool, info *data.User)
	Sync(info *data.User) (newinfo *data.User, err error)
	GetInfoByUID(uid int64) (has bool, info *data.User)
}

type MerchatAccountModule interface {
	Get(mid int64) (has bool, info data.MerchantAccount)
	GetByAccount(account string) (has bool, info data.MerchantAccount)
	IsAccountRepeat(mchID int64, name string) (has bool)
	IsMchAccountRepeat(mchID int64, account string) (has bool)
	Count(keyword string) (count int64, err error)
	Search(keyword string, offset, length int) (items []data.MerchantAccount, err error)
}

type MerchantStoreModule interface {
	Get(mid int64) (has bool, info data.MerchantStore)
	GetByAccount(account string) (has bool, info data.MerchantStore)
	IsAccountRepeat(mchID int64, name string) (has bool)
	Count(mchID int64, keyword string) (count int64, err error)
	Search(mchID int64, keyword string, offset, length int) (items []data.MerchantStore, err error)
}

type CardUsageModule interface {
	Count(mchID, storeID, start, end int64, isStore bool) (total int64, err error)
	Search(mchID, storeID, start, end int64, isStore bool, offset, length int) (items []data.AnnualCardUsageLog, err error)
}

var (
	Manager ManagerModule
	Token TokenModule
	Spread SpreadModule
	Conf ConfModule
	Order OrderModule
	AnnualCard AnnualCardModule
	Merchant MerchantModule
	MerchantAccount MerchatAccountModule
	MerchantStore MerchantStoreModule
	WxPay WxPayModule
	Wechat WechatModule
	User UserModule
	CardUsage CardUsageModule
)
