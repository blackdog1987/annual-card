package data

// EventMessage .
type EventMessage struct {
	XMLName      string `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Event        string
	Content      string  `xml:"Content,omitempty"`
	EventKey     string  `xml:"EventKey,omitempty"`
	Ticket       string  `xml:"Ticket,omitempty"`
	Latitude     float64 `xml:"Latitude,omitempty"`
	Longitude    float64 `xml:"Longitude,omitempty"`
	Precision    float64 `xml:"Precision,omitempty"`
}

// AccessToken .
type AccessToken struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int64  `json:"expires_in,omitempty"`
	ErrCode     int    `json:"errcode,omitempty"`
	ErrMsg      string `json:"errmsg,omitempty"`
}

// JSAPITicket .
type JSAPITicket struct {
	Ticket    string `json:"ticket,omitempty"`
	ExpiresIn int64  `json:"expires_in,omitempty"`
	ErrCode   int    `json:"errcode,omitempty"`
	ErrMsg    string `json:"errmsg,omitempty"`
}

// JsSignPackage .
// js api sign package
type JsSignPackage struct {
	AppID     string `json:"app_id"`
	NonceStr  string `json:"nonce_str"`
	Sign      string `json:"sign"`
	Timestamp int64  `json:"timestamp"`
	URL       string `json:"url"`
}

// UnifiedOrderReq .
// 下单请求
type UnifiedOrderReq struct {
	AppID          string `xml:"appid"`
	Body           string `xml:"body"`
	MchID          string `xml:"mch_id"`
	NonceStr       string `xml:"nonce_str"`
	NotifyURL      string `xml:"notify_url"`
	TradeType      string `xml:"trade_type"`
	SpbillCreateIP string `xml:"spbill_create_ip"`
	TotalFee       int64  `xml:"total_fee"`
	OutTradeNo     string `xml:"out_trade_no"`
	Sign           string `xml:"sign"`
	OpenID         string `xml:"openid"`
}

// UnifiedOrderResp .
// 下单返回
type UnifiedOrderResp struct {
	ReturnCode string `xml:"return_code" json:"-"`
	ReturnMsg  string `xml:"return_msg" json:"-"`
	AppID      string `xml:"appid" json:"-"`
	MchID      string `xml:"mch_id" json:"-"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code" json:"-"`
	PrepayID   string `xml:"prepay_id"`
	TradeType  string `xml:"trade_type"`
}

// WXPayNotifyReq .
// 支付回调
type WXPayNotifyReq struct {
	ReturnCode    string `xml:"return_code"`
	ReturnMsg     string `xml:"return_msg"`
	AppID         string `xml:"appid"`
	MchID         string `xml:"mch_id"`
	NonceStr      string `xml:"nonce_str"`
	Sign          string `xml:"sign"`
	ResultCode    string `xml:"result_code"`
	OpenID        string `xml:"openid"`
	IsSubscribe   string `xml:"is_subscribe"`
	TradeType     string `xml:"trade_type"`
	BankType      string `xml:"bank_type"`
	TotalFee      int    `xml:"total_fee"`
	FeeType       string `xml:"fee_type"`
	CashFee       int    `xml:"cash_fee"`
	CashFeeType   string `xml:"cash_fee_type"`
	TransactionID string `xml:"transaction_id"`
	OutTradeNo    string `xml:"out_trade_no"`
	Attach        string `xml:"attach"`
	TimeEnd       string `xml:"time_end"`
}

// WXPayNotifyResp .
type WXPayNotifyResp struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

// WxPayRedPackReq .
// 发红包请求结构
type WxPayRedPackReq struct {
	NonceStr    string `xml:"nonce_str"`    // 随机字符串
	Sign        string `xml:"sign"`         // 签名
	MchBillNo   string `xml:"mch_billno"`   // 商户订单号
	MchID       string `xml:"mch_id"`       // 商户号
	WxAppID     string `xml:"wxappid"`      // app id
	SendName    string `xml:"send_name"`    // 商户名称
	ReOpenID    string `xml:"re_openid"`    // 接受者 openid
	TotalAmount int64  `xml:"total_amount"` // 金额
	TotalNum    int    `xml:"total_num"`    // 总人数
	Wishing     string `xml:"wishing"`      // 祝福语
	ClientIP    string `xml:"client_ip"`    // host
	ActName     string `xml:"act_name"`     // 活动名称
	Remark      string `xml:"remark"`       // 备注
}

// WxPayRedPackResp .
// 发红包返回结构
type WxPayRedPackResp struct {
	ReturnCode  string `xml:"return_code"`
	ReturnMsg   string `xml:"return_msg"`
	Sign        string `xml:"sign"`
	ResultCode  string `xml:"result_code"`
	ErrCode     string `xml:"err_code"`
	ErrCodeDes  string `xml:"err_code_des"`
	MchBillNo   string `xml:"mch_billno"`
	MchID       string `xml:"mch_id"`       // 商户号
	WxAppID     string `xml:"wxappid"`      // app id
	ReOpenID    string `xml:"re_openid"`    // 接受者 openid
	TotalAmount int64  `xml:"total_amount"` // 金额
	SendTime    int64  `xml:"send_time"`    //红包发送时间
	SendListID  string `xml:"send_listid"`  //红包订单的微信单号
}

// UserOauth .
type UserOauth struct {
	UoID         int64  `json:"-" xorm:"int(11) autoincr pk 'uoid'"`
	UID          int64  `json:"uid,omitempty" xorm:"int(11) 'uid'"`
	IsSubscriber int    `json:"subscribe" xorm:"-"`              // 用户是否订阅该公众号标识, 值为0时, 代表此用户没有关注该公众号, 拉取不到其余信息
	OpenID       string `json:"openid" xorm:"char(32) 'openid'"` // 用户的标识, 对当前公众号唯一
	Nickname     string `json:"nickname"`                        // 用户的昵称
	Sex          int    `json:"sex"`                             // 用户的性别, 值为1时是男性, 值为2时是女性, 值为0时是未知
	Language     string `json:"language" xorm:"-"`               // 用户的语言, zh_CN, zh_TW, en
	City         string `json:"city"`                            // 用户所在城市
	Province     string `json:"province"`                        // 用户所在省份
	Country      string `json:"country"`                         // 用户所在国家
	Source       string `json:"source" xorm:"'source'"`
	// 用户头像, 最后一个数值代表正方形头像大小(有0, 46, 64, 96, 132数值可选, 0代表640*640正方形头像), 用户没有头像时该项为空
	HeadImageURL string `json:"headimgurl" xorm:"text 'headimgurl'"`

	SubscribeTime int64  `json:"subscribe_time" xorm:"-"`                  // 用户关注时间, 为时间戳. 如果用户曾多次关注, 则取最后关注时间
	UnionID       string `json:"unionid" xorm:"char(32) unique 'unionid'"` // 只有在用户将公众号绑定到微信开放平台帐号后, 才会出现该字段.
	Remark        string `json:"remark" xorm:"-"`                          // 公众号运营者对粉丝的备注, 公众号运营者可在微信公众平台用户管理界面对粉丝添加备注
	GroupID       int64  `json:"groupid" xorm:"-"`                         // 用户所在的分组ID
	ErrCode       int    `json:"errcode,omitempty" xorm:"-"`
	ErrMsg        string `json:"errmsg,omitempty" xorm:"-"`
}
