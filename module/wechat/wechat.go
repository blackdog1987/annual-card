package wechat

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/molibei/annual-card/lib/data"
	"github.com/molibei/annual-card/module"
	"github.com/molibei/huoda/model/wechat"
	"strconv"
)

// WechatService extend module.
type WechatService struct {
	token struct {
		token   string
		expires int64
		isFetch bool
		rw      sync.RWMutex
	}
	jsTicket struct {
		ticket  string
		expires int64
		isFetch bool
		rw      sync.RWMutex
	}
}

const (
	accessTokenURI      = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	oauthAccessTokenURI = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	oauthInfoURI        = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN"
	jsTicketURI         = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
	infoURI = "https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN"
	downMeidaURI        = "http://file.api.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s"
	snsInfoURI          = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=%s#wechat_redirect"
)

func init() {
	module.Wechat = &WechatService{}
}

// GetAccessToken .
func (p *WechatService) GetAccessToken() (token string, err error) {

	now := time.Now().Unix()
	if !p.token.isFetch || p.token.expires < now {
		uri := fmt.Sprintf(accessTokenURI, data.BaseConf.WeChat.AppID, data.BaseConf.WeChat.Secret)
		client := http.DefaultClient
		resp, err := client.Get(uri)
		if err != nil {
			return token, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return token, fmt.Errorf("http.Status: %s", resp.Status)
		}
		tk := wechat.AccessToken{}
		if err = json.NewDecoder(resp.Body).Decode(&tk); err != nil {
			return token, err
		}
		p.token.rw.Lock()
		defer p.token.rw.Unlock()
		p.token.token = tk.AccessToken
		p.token.expires = now + tk.ExpiresIn - 100
		p.token.isFetch = true

	}
	return p.token.token, nil
}

// 获取jsapi ticket
func (p *WechatService) GetJsAPITicket() (ticket string, err error) {
	now := time.Now().Unix()
	if !p.jsTicket.isFetch || p.jsTicket.expires < now {
		token, err := p.GetAccessToken()
		if err != nil {
			return ticket, errors.New("get access_token failed.")
		}
		uri := fmt.Sprintf(jsTicketURI, token)
		client := http.DefaultClient
		resp, err := client.Get(uri)
		if err != nil {
			return ticket, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return token, fmt.Errorf("http.Status: %s", resp.Status)
		}
		tk := data.JSAPITicket{}
		if err = json.NewDecoder(resp.Body).Decode(&tk); err != nil {
			return ticket, err
		}
		p.jsTicket.rw.Lock()
		defer p.jsTicket.rw.Unlock()
		p.jsTicket.ticket = tk.Ticket
		p.jsTicket.expires = now + tk.ExpiresIn - 100
		p.jsTicket.isFetch = true
	}
	return p.jsTicket.ticket, nil
}

// 生成jsapi 签名
func (p *WechatService) GenerateJsSign(url string) (sign *data.JsSignPackage) {
	sign = &data.JsSignPackage{
		AppID:     data.BaseConf.WeChat.AppID,
		NonceStr:  p.GenerateNonceStr(8),
		Timestamp: time.Now().Unix(),
		URL:       url,
	}
	ticket, _ := p.GetJsAPITicket()
	temp := []string{
		fmt.Sprintf("noncestr=%s", sign.NonceStr),
		fmt.Sprintf("jsapi_ticket=%s", ticket),
		fmt.Sprintf("timestamp=%d", sign.Timestamp),
		fmt.Sprintf("url=%s", url),
	}
	sort.Strings(temp)
	tmp := strings.Join(temp, "&")
	sha := sha1.New()
	io.WriteString(sha, tmp)
	sign.Sign = fmt.Sprintf("%x", sha.Sum(nil))
	return sign
}

// 随机字符串
func (p *WechatService) GenerateNonceStr(length int) (nonce string) {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// GenerateSnsInfoURI .
func (p *WechatService) GenerateSnsInfoURI(redirect string) (uri string) {
	redirect = url.QueryEscape(redirect)
	return fmt.Sprintf(snsInfoURI, data.BaseConf.WeChat.AppID, redirect, data.BaseConf.WeChat.State)
}

// GetUserInfo .
func (p *WechatService) GetUserInfo(msg data.EventMessage) (useroauth data.UserOauth) {
	return useroauth
}

type OauthAccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `josn:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
}

// OauthInfo .
func (p *WechatService) OauthInfo(code string) (info *data.User, err error) {
	uri := fmt.Sprintf(oauthAccessTokenURI, data.BaseConf.WeChat.AppID, data.BaseConf.WeChat.Secret, code)
	client := http.DefaultClient
	resp, err := client.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return info, fmt.Errorf("http.Status: %s", resp.Status)
	}
	var otoken = OauthAccessToken{}
	if err = json.NewDecoder(resp.Body).Decode(&otoken); err != nil {
		return
	}
	// check register
	has, info := module.User.GetInfoByOpenID(otoken.OpenID)
	if has {
		return info, nil
	}

	uri = fmt.Sprintf(oauthInfoURI, otoken.AccessToken, otoken.OpenID)
	resp, err = client.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return info, fmt.Errorf("http.Status: %s", resp.Status)
	}
	useroauth := &data.UserOauth{}
	if err = json.NewDecoder(resp.Body).Decode(&useroauth); err != nil {
		return
	}
	if useroauth.ErrCode != 0 {
		return info, errors.New(useroauth.ErrMsg)
	}
	// sync user info.
	info.Nickname = useroauth.Nickname
	info.Sex = useroauth.Sex
	info.HeadImageURL = useroauth.HeadImageURL
	info.Country = useroauth.Country
	info.Province = useroauth.Province
	info.City = useroauth.City
	info.WxOpenID = useroauth.OpenID
	info.WxUnionID = useroauth.UnionID
	info, err = module.User.Sync(info)
	return
}

// Subscribe .
func (p *WechatService) Subscribe(msg data.EventMessage) data.EventMessage {
	// 获取用户信息
	var tmp = msg.ToUserName
	msg.ToUserName = msg.FromUserName
	msg.FromUserName = tmp
	msg.CreateTime = time.Now().Unix()
	msg.MsgType = "text"
	msg.Content = "欢迎来到公号"
	// 检测用户是否关注过
	has, info := module.User.GetInfoByOpenID(msg.ToUserName)
	if has {
		return msg
	}
	tk, err := p.GetAccessToken()
	if err != nil {
		return msg
	}
	uri := fmt.Sprintf(infoURI, tk, msg.ToUserName)
	client := http.DefaultClient
	resp, err := client.Get(uri)
	if err != nil {
		return msg
	}
	defer resp.Body.Close()
	uo :=  &data.UserOauth{}
	if err = json.NewDecoder(resp.Body).Decode(uo); err != nil {
		return msg
	}
	if uo.ErrCode != 0 {
		return msg
	}
	info.Nickname = uo.Nickname
	info.Sex = uo.Sex
	info.HeadImageURL = uo.HeadImageURL
	info.Country = uo.Country
	info.Province = uo.Province
	info.City = uo.City
	info.WxOpenID = uo.OpenID
	info.WxUnionID = uo.UnionID
	// 同步用户信息
	info, err =module.User.Sync(info)
	if err != nil {
		fmt.Println("insert info failed of ", info)
		return msg
	}
	if msg.EventKey != "" {
		ek := strings.Split(msg.EventKey, "_")
		sp, err := strconv.ParseInt(ek[1], 10, 64)
		if err != nil {
			return msg
		}
		// 获取推广信息
		var spreadAccess = &data.SpreadAccess{}
		has, _ := data.Db.Where("access_id = ?", sp).Get(spreadAccess)
		if !has {
			return msg
		}
		switch spreadAccess.Category {
		case "channel": // 渠道
			plan, has :=module.Spread.Plan(spreadAccess.RelationID)
			if has && plan.IsDisabled == 0 {
				spreadLogs := &data.SpreadLogs{
					PlanID:plan.SPID,
					UID:info.UID,
					Category:1,
					Commission:plan.RegCommission,
				}
				info.SpreadUID = spreadAccess.AccessID
				data.Db.Where("uid = ?", info.UID).Update(info)
				data.Db.Insert(spreadLogs)
			}
		case "user":	// 积分
			has, spreadInfo := module.User.GetInfoByUID(spreadAccess.RelationID)
			if has {
				fmt.Println(spreadInfo)
				_, reward := module.Conf.GetInt64("spread_reg")
				info.SpreadUID = spreadAccess.AccessID
				info.PointsContribute = reward
				// 给推广人添加积分
				spreadInfo.PointsBalance += reward
				spreadInfo.PointsEarning += reward
				// 积分日志
				integralLog := &data.PointsLog{
					UID:spreadAccess.RelationID,
					RelationUID:info.UID,
					FriendlyIntro:"推广「" + info.Nickname + "」获得积分",
					Total:reward,
					Type:1,
				}
				data.Db.Where("uid = ?", info.UID).Update(info)
				data.Db.Where("uid = ?", spreadInfo.UID).Update(spreadInfo)
				data.Db.Insert(integralLog)
			}
		}
	}
	return msg
}

// UnSubscribe .
func (p *WechatService) UnSubscribe(msg data.EventMessage) {

}

// 下载文件
func (p *WechatService) Download(mediaID, filename string) (ok bool) {
	tk, err := p.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf(downMeidaURI, tk, mediaID)
	fmt.Println(uri)
	client := http.DefaultClient
	resp, err := client.Get(uri)
	if err != nil {
		fmt.Println(err)
		return
	}
	//if resp.Header.Get("Content-Type") != "audio/amr" {
	//	return
	//}
	defer resp.Body.Close()
	//filename = filename + ".amr"
	fl, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fl.Close()
	written, err := io.Copy(fl, resp.Body)
	if err != nil || written == 0 {
		fmt.Println(err)
		return
	}
	return true

}
