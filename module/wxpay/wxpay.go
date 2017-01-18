package wxpay

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"hash"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/molibei/annual-card/lib/data"
	"github.com/molibei/annual-card/module"
)

type WxPayService struct{}

const (
	unifiedOrderURL = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	sendRedPackURL  = "https://api.mch.weixin.qq.com/mmpaymkttransfers/sendredpack"
)

func init() {
	module.WxPay = &WxPayService{}
}

// 统一下单
func (p *WxPayService) UnifiedOrder(msg, clientIP string, orderID, openID string, total int64) (res map[string]string, err error) {
	nonceStr := module.Wechat.GenerateNonceStr(32)
	var ureq = data.UnifiedOrderReq{
		AppID:          data.BaseConf.WeChat.AppID,
		Body:           msg,
		MchID:          data.BaseConf.WeChat.MchID,
		NonceStr:       nonceStr,
		NotifyURL:      data.BaseConf.WeChat.PayNotifyURL,
		TradeType:      "JSAPI",
		SpbillCreateIP: clientIP,
		TotalFee:       total,
		OutTradeNo:     orderID,
		OpenID:         openID,
	}
	m := make(map[string]string, 0)
	m["body"] = msg
	m["out_trade_no"] = orderID
	m["total_fee"] = fmt.Sprintf("%d", total)
	m["trade_type"] = "JSAPI"
	m["notify_url"] = ureq.NotifyURL
	m["appid"] = ureq.AppID
	m["mch_id"] = ureq.MchID
	m["spbill_create_ip"] = clientIP
	m["nonce_str"] = nonceStr
	m["openid"] = openID
	ureq.Sign = sign(m, data.BaseConf.WeChat.Key, md5.New)
	bytes_req, err := xml.Marshal(ureq)
	if err != nil {
		return
	}
	str_req := string(bytes_req)
	str_req = strings.Replace(str_req, "UnifiedOrderReq", "xml", -1)
	request, err := http.NewRequest("POST", unifiedOrderURL, bytes.NewReader(bytes_req))
	if err != nil {
		return
	}
	request.Header.Set("Accept", "application/xml")
	request.Header.Set("Content-Type", "application/xml;charset=utf-8")
	c := http.Client{}
	result, _err := c.Do(request)
	if _err != nil {
		return
	}
	defer result.Body.Close()
	resp := &data.UnifiedOrderResp{}
	err = xml.NewDecoder(result.Body).Decode(resp)
	if resp == nil {
		return
	}
	var m1 = make(map[string]string, 0)
	m1["appId"] = data.BaseConf.WeChat.AppID
	m1["package"] = "prepay_id=" + resp.PrepayID
	m1["nonceStr"] = resp.NonceStr
	m1["signType"] = "MD5"
	m1["timeStamp"] = fmt.Sprintf("%d", time.Now().Unix())
	m1["sign"] = sign(m1, data.BaseConf.WeChat.Key, md5.New)
	return m1, nil
}

// 发红包
func (p *WxPayService) SendRedPack(openID string, total int64, orderID string) (ok bool, err error) {
	req := &data.WxPayRedPackReq{
		NonceStr:    module.Wechat.GenerateNonceStr(32),
		MchBillNo:   orderID,
		MchID:       data.BaseConf.WeChat.MchID,
		WxAppID:     data.BaseConf.WeChat.AppID,
		SendName:    data.BaseConf.WeChat.Name,
		ReOpenID:    openID,
		TotalAmount: total,
		TotalNum:    1,
		Wishing:     "感谢您在[听健]的贡献",
		ClientIP:    data.BaseConf.Server.IP,
		ActName:     "[听健]收益提现",
		Remark:      "答题质量越高,分享收益也越高",
	}
	m := make(map[string]string, 0)
	m["nonce_str"] = req.NonceStr
	m["mch_billno"] = req.MchBillNo
	m["mch_id"] = req.MchID
	m["wxappid"] = req.WxAppID
	m["send_name"] = req.SendName
	m["re_openid"] = req.ReOpenID
	m["total_amount"] = fmt.Sprintf("%d", req.TotalAmount)
	m["total_num"] = fmt.Sprintf("%d", req.TotalNum)
	m["wishing"] = req.Wishing
	m["client_ip"] = req.ClientIP
	m["act_name"] = req.ActName
	m["remark"] = req.Remark
	req.Sign = sign(m, data.BaseConf.WeChat.Key, md5.New)
	bytes_req, err := xml.Marshal(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	str_req := string(bytes_req)
	str_req = strings.Replace(str_req, "WxPayRedPackReq", "xml", -1)
	wtls, err := wxtls()
	if err != nil {
		return
	}
	tr := &http.Transport{TLSClientConfig: wtls}
	client := &http.Client{Transport: tr}
	result, err := client.Post(sendRedPackURL, "text/xml", bytes.NewBuffer(bytes_req))
	if err != nil {
		return
	}
	defer result.Body.Close()
	resp := &data.WxPayRedPackResp{}
	err = xml.NewDecoder(result.Body).Decode(resp)
	if resp == nil {
		return
	}
	if resp.ReturnCode != "SUCCESS" {
		return ok, errors.New(resp.ReturnMsg)
	}
	if resp.ResultCode == "FAIL" {
		return ok, errors.New(resp.ErrCodeDes)
	}
	// 处理账户
	orderInfo := &data.OrderLog{}
	if has, err := data.Db.Where("order_id = ?", resp.MchBillNo).Get(orderInfo); !has || err != nil {
		// TODO Log
		return ok, err
	}
	// orderInfo.IsPay = true
	// orderInfo.TransactionID = resp.SendListID
	// orderInfo.TimeEnd = fmt.Sprintf("%d", resp.SendTime)
	// data.Db.Where("order_id = ?", orderInfo.OrderID).Update(orderInfo)
	// // 修改余额
	// module.User.Expend(orderInfo.UID, orderInfo.Total)
	// module.User.NewCashLog(orderInfo.UID, orderInfo.Total, orderInfo.OrderID, "cash")
	return true, nil
}

// 签名
func sign(params map[string]string, apiKey string, fn func() hash.Hash) string {
	if fn == nil {
		fn = md5.New
	}
	h := fn()
	bufw := bufio.NewWriterSize(h, 128)

	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "sign" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := params[k]
		if v == "" {
			continue
		}
		bufw.WriteString(k)
		bufw.WriteByte('=')
		bufw.WriteString(v)
		bufw.WriteByte('&')
	}
	bufw.WriteString("key=")
	bufw.WriteString(apiKey)

	bufw.Flush()
	signature := make([]byte, hex.EncodedLen(h.Size()))
	hex.Encode(signature, h.Sum(nil))
	return string(bytes.ToUpper(signature))
}

func wxtls() (wt *tls.Config, err error) {
	cert, err := tls.LoadX509KeyPair(data.BaseConf.WeChat.CertCertPath, data.BaseConf.WeChat.CertKeyPath)
	if err != nil {
		return nil, err
	}
	// load root ca
	caData, err := ioutil.ReadFile(data.BaseConf.WeChat.CertCaPath)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}, nil
}
