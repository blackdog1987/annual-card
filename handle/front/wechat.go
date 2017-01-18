package front

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/blackdog1987/annual-card/lib/data"
	"github.com/blackdog1987/annual-card/module"
	"time"
)

// message type.
const (
	TextMsg = "text"       // 文本消息
	ImageMsg = "image"      // 图片消息
	VoiceMsg = "voice"      // 语音消息
	VideoMsg = "video"      // 视频消息
	ShortVideoMsg = "shortvideo" // 短视频消息
	LocationMsg = "location"   // 位置消息
	LinkMsg = "link"       //链接消息
	EventMsg = "event"      // 事件消息
)

// event
const (
	SubscribeEvent = "subscribe"   // 关注事件
	UnSubscribeEvent = "unsubscribe" // 取消关注事件
	ScanEvent = "scan"        // 扫描
	LocationEvent = "location"    // 位置
	ClickEvent = "click"       // 点击
	ViewEvent = "view"        // 显示
)

func WxValid(ctx *gin.Context) {
	// valid request
	signature := ctx.Query("signature")
	timestamp := ctx.Query("timestamp")
	nonce := ctx.Query("nonce")
	echostr := ctx.Query("echostr")
	tmpstr := []string{data.BaseConf.WeChat.Token, timestamp, nonce}
	sort.Strings(tmpstr)
	tmp := strings.Join(tmpstr, "")
	sha := sha1.New()
	io.WriteString(sha, tmp)
	if fmt.Sprintf("%x", sha.Sum(nil)) != signature {
		ctx.String(http.StatusOK, "valid failed.")
		return
	}
	ctx.String(http.StatusOK, echostr)
}

// Receive .
// POST /wechat/receive
func Receive(ctx *gin.Context) {
	// valid request
	signature := ctx.Query("signature")
	timestamp := ctx.Query("timestamp")
	nonce := ctx.Query("nonce")
	tmpstr := []string{data.BaseConf.WeChat.Token, timestamp, nonce}
	sort.Strings(tmpstr)
	tmp := strings.Join(tmpstr, "")
	sha := sha1.New()
	io.WriteString(sha, tmp)
	if fmt.Sprintf("%x", sha.Sum(nil)) != signature {
		ctx.String(http.StatusOK, "valid failed.")
		return
	}
	var msg = data.EventMessage{}
	if err := ctx.BindWith(&msg, binding.XML); err == nil {
		switch msg.MsgType {
		case EventMsg:
			switch msg.Event {
			case SubscribeEvent:
				ctx.XML(http.StatusOK, module.Wechat.Subscribe(msg))
				return
			default:
				ctx.String(http.StatusOK, "")
				return
			}
		default: // 暂不处理其他消息
			ctx.String(http.StatusOK, "")
			return
		}
	}
	ctx.String(http.StatusOK, "")
}

// Authorize .
func Authorize(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	from := ctx.Query("from")
	if state != data.BaseConf.WeChat.State {
		ctx.String(400, "Request valid failed.")
		ctx.Abort()
		return
	}
	info, err := module.Wechat.OauthInfo(code)
	if err != nil {
		ctx.String(400, err.Error(), nil)
		ctx.Abort()
		return
	}
	session := sessions.Default(ctx)
	session.Set("u", info.UID)
	session.Save()
	ctx.Redirect(http.StatusMovedPermanently, data.BaseConf.Server.Host + "/?#/" + from)
}

func WxPayNotify(ctx *gin.Context) {
	var req = &data.WXPayNotifyReq{}
	var resp = data.WXPayNotifyResp{
		ReturnCode: "FAIL",
	}
	err := ctx.BindWith(&req, binding.XML)
	fmt.Println(req, err)
	if err != nil {
		ctx.XML(http.StatusOK, resp)
		return
	}
	// 获取订单
	var order = &data.OrderLog{}
	if has, _ := data.Db.Where("order_no = ?", req.OutTradeNo).Get(order); !has || order == nil {
		ctx.XML(http.StatusOK, resp)
		return
	}
	// 修改状态
	order.IsPay = 1
	order.TransactionID = req.TransactionID
	data.Db.Where("order_id = ?", order.OrderID).Update(order)
	// 返佣或者积分
	has, userInfo := module.User.GetInfoByUID(order.UID)
	if has {
		// 添加商品记录
		switch order.Category {
		case "annual":
			goods := make([]data.AnnualCard, order.GoodsNum)
			var info data.AnnualCardConf
			module.Conf.GetObject("annual_card", &info)
			for i := 0; i < int(order.GoodsNum); i++ {
				goods[i] = data.AnnualCard{
					PlanID:      0,
					CardName:    info.Name,
					CardNO:      fmt.Sprintf("%s%d", "L", time.Now().UnixNano()),
					RelationUID: order.UID,
				}
			}
			data.Db.Insert(&goods)
		case "coupon":
			goods := make([]data.CouponLog, order.GoodsNum)
			var info data.CouponConf
			module.Conf.GetObject("coupon", &info)
			for i := 0; i < int(order.GoodsNum); i++ {
				goods[i] = data.CouponLog{
					UID:         userInfo.UID,
					OffsetPrice: info.OffsetPrice,
					IsUsage:     0,
				}
			}
			data.Db.Insert(&goods)
		}
		// 处理推广
		if userInfo.SpreadUID > 0 {
			module.Spread.Commission(userInfo.SpreadUID, order.Price, userInfo.UID, order.OrderID)
		}
	}
	resp.ReturnCode = "SUCCESS"
	ctx.XML(http.StatusOK, resp)
}

// 获取js签名
func GetWxJsSign(ctx *gin.Context) {
	uri := ctx.Query("uri")
	uris := strings.Split(uri, "#")
	sign := module.Wechat.GenerateJsSign(uris[0])
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"info": "success",
		"data": sign,
	})
}
