package front

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/molibei/annual-card/lib/data"
	"github.com/molibei/annual-card/module"
	"github.com/molibei/page"
	"github.com/molibei/annual-card/lib/errors"
	"bytes"
	"os"
	"image/jpeg"
	"encoding/base64"
	"image"
	"golang.org/x/image/draw"
	"strings"
	"github.com/nfnt/resize"
	qrcode "github.com/skip2/go-qrcode"
)

func Goods(ctx *gin.Context) {
	var annual data.AnnualCardConf
	module.Conf.GetObject("ANNUAL_CARD", &annual)
	var coupon data.CouponConf
	module.Conf.GetObject("COUPON", &coupon)
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"annual": annual,
			"coupon": coupon,
		},
	})
}

// Merchants .
// /v1/manage/merchants GET
// 获取商户列表
func Merchants(ctx *gin.Context) {
	pstr := ctx.Query("p")
	p, _ := strconv.Atoi(pstr)
	total, _ := module.Merchant.Count("")
	//offset, length := page.Page(int(total), 10, p, 10)
	var items []data.Merchant
	//if offset == 0 && p >= 2 {
	//	ctx.JSON(http.StatusOK, gin.H{
	//		"code":  0,
	//		"count": 0,
	//		"total": total,
	//		"p":     p,
	//		"data":  items,
	//	})
	//	return
	//}
	items, _ = module.Merchant.Search("", 0, int(total))
	ctx.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": len(items),
		"total": total,
		"p":     p,
		"data":  items,
	})
}
func Merchant(ctx *gin.Context) {
	idStr := ctx.Query("mch_id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	has, info := module.Merchant.Get(id)
	if !has {
		ctx.JSON(http.StatusOK, errors.ID_EMPTY)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": info,
	})
}

func UserInfo(ctx *gin.Context) {
	u, _ := ctx.Get("uid")
	uid := u.(int64)
	has, info := module.User.GetInfoByUID(uid)
	if !has {
		ctx.JSON(http.StatusOK, errors.ID_EMPTY)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": info,
	})
}

func Spreads(ctx *gin.Context) {
	u, _ := ctx.Get("uid")
	uid := u.(int64)
	pstr := ctx.Query("p")
	p, _ := strconv.Atoi(pstr)
	total, _ := module.User.SpreadCount(uid)
	offset, length := page.Page(int(total), 20, p, 20)
	var items []data.User
	if offset == 0 && p >= 2 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":  0,
			"count": 0,
			"total": total,
			"p":     p,
			"data":  items,
		})
		return
	}
	items, _ = module.User.SpreadSearch(uid, offset, length)
	ctx.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": len(items),
		"total": total,
		"p":     p,
		"data":  items,
	})
}

func Picture(ctx *gin.Context) {
	u, _ := ctx.Get("uid")
	uid := u.(int64)
	_, info := module.User.GetInfoByUID(uid)
	qrcode, err := module.Spread.QRCode(uid, "user")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":-1,
			"msg":"获取二维码失败",
		})
		return
	}
	// 合并图片
	// step 1 读取背景
	bg, _ := os.Open(data.BaseConf.Server.Picture)
	bgImg, _ := jpeg.Decode(bg)
	defer bg.Close()
	// 头像
	client := http.DefaultClient
	headUri := strings.TrimRight(info.HeadImageURL, "0")
	resp, err := client.Get(headUri + "64")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "读取头像失败",
		})
		return
	}
	defer resp.Body.Close()
	head, _ := jpeg.Decode(resp.Body)
	// 二维码
	qrByte, _ := base64.StdEncoding.DecodeString(qrcode)
	qr, _ := jpeg.Decode(bytes.NewReader(qrByte))
	// 缩略二维码到180*180
	qrThumb := resize.Thumbnail(180, 180, qr, resize.Bicubic)
	//把水印写到右下角，并向0坐标各偏移10个像素
	headOffset := image.Pt(22, 596)
	bgBounds := bgImg.Bounds()
	bgRgba := image.NewRGBA(bgBounds)
	// 写入原图
	draw.Draw(bgRgba, bgBounds, bgImg, image.ZP, draw.Src)
	// 写入头像
	draw.Draw(bgRgba, head.Bounds().Add(headOffset), head, image.ZP, draw.Src)
	// 第二个头像
	draw.Draw(bgRgba, head.Bounds().Add(image.Pt(22, 732)), head, image.ZP, draw.Src)
	// 写入二维码
	draw.Draw(bgRgba, qr.Bounds().Add(image.Pt(337, 588)), qrThumb, image.ZP, draw.Src)
	emptyBuff := bytes.NewBuffer(nil)                  //开辟一个新的空buff
	jpeg.Encode(emptyBuff, bgRgba, nil)
	//fmt.Println(emptyBuff.Len())                  //开辟存储空间
	pic := base64.StdEncoding.EncodeToString(emptyBuff.Bytes())

	ctx.JSON(http.StatusOK, gin.H{
		"code":0,
		"data": "data:image/jpeg;base64," + pic,
	})
}

func CardNum(ctx *gin.Context) {
	u, _ := ctx.Get("uid")
	uid := u.(int64)
	active, _ := module.AnnualCard.CountBind(uid, 1)
	noActive, _ := module.AnnualCard.CountBind(uid, 0)
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"active": active,
			"noActive":noActive,
		},
	})
}

func Cards(ctx *gin.Context) {
	u, _ := ctx.Get("uid")
	uid := u.(int64)
	items, _ := module.AnnualCard.SearchBind(uid, -1, 0, 20)
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": items,
	})
}

// 获取年卡信息
func Card(ctx *gin.Context) {
	u, _ := ctx.Get("uid")
	uid := u.(int64)
	cardIdStr := ctx.Query("card_id")
	cardId, _ := strconv.ParseInt(cardIdStr, 10, 64)
	info, has := module.AnnualCard.CardByID(cardId)
	if !has {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "该年卡不存在",
		})
		return
	}
	if info.RelationUID != uid {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "这不是您的年卡",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": info,
	})
}
// 绑定年卡
func BindCard(ctx *gin.Context) {
	u, _ := ctx.Get("uid")
	uid := u.(int64)
	var in data.AnnualCard
	if err := ctx.Bind(&in); err != nil {
		ctx.JSON(http.StatusOK, errors.PARAM_PARSE_ERR)
		return
	}
	info, has := module.AnnualCard.CardByID(in.CardID)
	if !has {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "你要操作的年卡不存在",
		})
		return
	}
	if info.RelationUID != uid {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "这不是您的年卡",
		})
		return
	}
	//if info.IsActive == 1 {
	//	ctx.JSON(http.StatusOK, gin.H{
	//		"code": -1,
	//		"msg": "该年卡已经激活过",
	//	})
	//	return
	//}
	_, has = module.AnnualCard.IsExistIDCard(info.CardID, in.BindIDCard)
	if has {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "该身份证已经绑定过年卡",
		})
		return
	}
	if info.IsActive != 1 {
		now := time.Now()
		// 检查年卡的激活有效期
		if info.PlanID > 0 && info.ExpiredStop > 0 {
			if now.Unix() > info.ExpiredStop {
				ctx.JSON(http.StatusOK, gin.H{
					"code": -1,
					"msg": "该年卡已经过了有效激活期",
				})
				return
			}
		}
		z, _ := time.Parse("2006-01-02", now.Format("2006-01-02"))
		info.ExpiredStart = z.Unix()
		z.AddDate(1, 0, 0)
		info.ExpiredStop = z.Unix()
		info.BindIDCard = in.BindIDCard
	}
	info.BindHeadimg = in.BindHeadimg
	info.BindContact = in.BindContact
	info.BindName = in.BindName
	// 处理head
	if !strings.Contains(in.BindHeadimg, ".jpg") {
		// 去微信下载图片
		fi := fmt.Sprintf("./head/%s.jpg", in.CardNO)
		ok := module.Wechat.Download(in.BindHeadimg, fi)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg": "处理认证头像失败",
			})
			return
		}
		info.BindHeadimg = strings.TrimLeft(fi, ".")
	}
	info.IsActive = 1
	_, err := data.Db.Where("card_id = ?", info.CardID).Update(info)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "提交资料失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, errors.SUCCESS)
}

func AddCard(ctx *gin.Context) {
	u, _ := ctx.Get("uid")
	uid := u.(int64)
	cardNo := ctx.PostForm("cardNo")
	cardPass := ctx.PostForm("cardPass")
	info, has := module.AnnualCard.CardByCardNo(cardNo)
	if !has {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -2400,
			"msg": "卡号不存在",
		})
		return
	}
	if info.CardPasswd != cardPass {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -2401,
			"msg": "卡密错误",
		})
		return
	}
	info.RelationUID = uid
	affected, err := data.Db.Where("card_id = ?", info.CardID).Update(info)
	if affected < 1 || err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -2402,
			"msg": "绑定失败,稍后再试",
		})
		return
	}
	ctx.JSON(http.StatusOK, errors.SUCCESS)
}

func QrCode(ctx *gin.Context) {
	u, _ := ctx.Get("uid")
	uid := u.(int64)
	var png []byte
	uri := fmt.Sprintf("%s/wxmerchant/#/usage/QR100%d", data.BaseConf.Server.Host, uid)
	png, err := qrcode.Encode(uri, qrcode.Medium, 256)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "二维码处理失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data":gin.H{
			"qrcode":base64.StdEncoding.EncodeToString(png),
			"uid": fmt.Sprintf("QR100%d", uid),
		},
	})
}

func PreOrder(ctx *gin.Context) {
	u, _ := ctx.Get("uid")
	uid := u.(int64)
	// 获取用户信息
	_, userInfo := module.User.GetInfoByUID(uid)
	_, rate := module.Conf.GetFloat64("integral_rate")
	// 获取未使用的优惠券
	d := gin.H{
		"user": userInfo,
		"rate": rate,
		"coupon": false,
	}
	coupon, has := module.Order.Coupon(uid)
	if has {
		d["coupon"] = coupon
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": d,
	})
}

func GetPaySign(ctx *gin.Context) {
	category := ctx.PostForm("product")
	numStr := ctx.PostForm("num")
	num, _ := strconv.ParseInt(numStr, 10, 64)
	u, _ := ctx.Get("uid")
	uid := u.(int64)
	_, userInfo := module.User.GetInfoByUID(uid)
	_, rate := module.Conf.GetFloat64("integral_rate")
	// 构造订单
	order := data.OrderLog{
		OrderNO:  fmt.Sprintf("%d", time.Now().UnixNano()),
		UID:      uid,
		Category: category,
		GoodsNum: num,
	}
	// 获取产品
	switch category {
	case "coupon":
		var info data.CouponConf
		module.Conf.GetObject("coupon", &info)
		// 计算要支付的金额
		price := info.SalePrice * num
		order.Price = price
		order.GoodsName = info.Name
		order.Total = price
		if userInfo.PointsBalance > 0 {
			// 有积分可以使用
			// 计算所需积分
			points := order.Total / int64((rate * float64(100)))
			expend := data.PointsLog{
				UID:           uid,
				RelationUID:   uid,
				FriendlyIntro: "积分兑换" + info.Name,
				Type:          0,
				Total:         points,
			}
			if userInfo.PointsBalance >= points {
				// 积分可以全部支付
				// 修改订单
				order.Points = points
				order.IsPay = 1
				order.PointsPrice = order.Total
				order.Total = 0 // 积分全部抵扣
				// 更新余额
				userInfo.PointsBalance -= points
				userInfo.PointsExpend += points
				// 增加记录
				expend.Total = points
				goods := make([]data.CouponLog, num)
				for i := 0; i < int(num); i++ {
					goods[i] = data.CouponLog{
						UID:         uid,
						OffsetPrice: info.OffsetPrice,
						IsUsage:     0,
					}
				}
				// 事务
				session := data.Db.NewSession()
				defer session.Close()
				err := session.Begin()
				if err != nil {
					fmt.Println(err)
				}
				session.Insert(&order)
				expend.RelationLogID = order.OrderID
				session.Where("uid = ?", uid).Cols("points_earning", "points_expend", "points_balance").Update(userInfo)
				session.Insert(&expend)
				session.Insert(&goods)
				err = session.Commit()
				if err != nil {
					ctx.JSON(http.StatusOK, gin.H{
						"code": -1,
						"msg": "订单处理失败,稍后再试",
					})
					return
				}
				// 处理推广
				if userInfo.SpreadUID > 0 {
					module.Spread.Commission(userInfo.SpreadUID, order.Price, userInfo.UID, order.OrderID)
				}
				ctx.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg": "支付成功",
				})
				return
			}
			// 积分不够 重新计算金额
			order.Points = userInfo.PointsBalance
			expend.Total = userInfo.PointsBalance
			userInfo.PointsExpend += userInfo.PointsBalance
			userInfo.PointsBalance = 0
			// 计算需要支付的金额
			// 抵扣金额
			pointsPrice := order.Points * int64(rate * float64(100))
			order.Total = order.Total - pointsPrice
			order.PointsPrice = pointsPrice
			// 事务
			session := data.Db.NewSession()
			defer session.Close()
			err := session.Begin()
			if err != nil {
				fmt.Println(err)
			}
			session.Insert(&order)
			expend.RelationLogID = order.OrderID
			session.Where("uid = ?", uid).Cols("points_earning", "points_expend", "points_balance").Update(userInfo)
			session.Insert(&expend)
			err = session.Commit()
			if err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"code": -1,
					"msg":"订单处理失败,稍后再试",
				})
				return
			}
		} else {
			affected, err := data.Db.Insert(&order)
			if affected < 1 || err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"code": -1,
					"msg":"订单处理失败,稍后再试",
				})
				return
			}
		}
	case "annual":
		var info data.AnnualCardConf
		module.Conf.GetObject("annual_card", &info)
		order.GoodsName = info.Name
		total := info.SalePrice * num
		order.Price = total
		// 检测有没有能使用的优惠券
		coupon, has := module.Order.Coupon(uid)
		if has {
			total -= coupon.OffsetPrice
			order.IsCoupon = coupon.LogID
			order.CouponPrice = coupon.OffsetPrice
			coupon.IsUsage = 1
			data.Db.Where("log_id = ?", coupon.LogID).Update(coupon)
		}
		order.Total = total
		if userInfo.PointsBalance > 0 {
			// 有积分可用
			points := order.Total / int64((rate * float64(100)))
			expend := data.PointsLog{
				UID:           uid,
				RelationUID:   uid,
				FriendlyIntro: "积分兑换" + info.Name,
				Type:          0,
				Total:         points,
			}
			if userInfo.PointsBalance >= points {
				// 积分可以全部支付
				order.Points = points
				order.IsPay = 1
				order.PointsPrice = order.Total
				// 更新余额
				userInfo.PointsBalance -= points
				userInfo.PointsExpend += points
				// 增加记录
				expend.Total = 0
				// 商品
				goods := make([]data.AnnualCard, num)
				for i := 0; i < int(num); i++ {
					goods[i] = data.AnnualCard{
						PlanID:      0,
						CardName:    info.Name,
						CardNO:      fmt.Sprintf("%s%d", "L", time.Now().UnixNano()),
						RelationUID: uid,
					}
				}
				// 事务
				session := data.Db.NewSession()
				defer session.Close()
				err := session.Begin()
				if err != nil {
					fmt.Println(err)
				}
				session.Insert(&order)
				expend.RelationLogID = order.OrderID
				session.Where("uid = ?", uid).Update(userInfo)
				session.Insert(&expend)
				session.Insert(&goods)
				err = session.Commit()
				if err != nil {
					ctx.JSON(http.StatusOK, gin.H{
						"code": -1,
						"msg": "订单处理失败,稍后再试",
					})
					return
				}
				// 处理推广
				if userInfo.SpreadUID > 0 {
					module.Spread.Commission(userInfo.SpreadUID, order.Price, userInfo.UID, order.OrderID)
				}
				ctx.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg": "支付成功",
				})
				return
			}
			// 积分不够 重新计算金额
			order.Points = userInfo.PointsBalance
			expend.Total = userInfo.PointsBalance
			userInfo.PointsExpend += userInfo.PointsBalance
			userInfo.PointsBalance = 0
			// 计算需要支付的金额
			pointsPrice := order.Points * int64(rate * float64(100))
			order.Total = order.Total - pointsPrice
			order.PointsPrice = pointsPrice
			// 事务
			session := data.Db.NewSession()
			defer session.Close()
			err := session.Begin()
			if err != nil {
				fmt.Println(err)
			}
			session.Insert(&order)
			expend.RelationLogID = order.OrderID
			session.Where("uid = ?", uid).Cols("points_balance", "points_expend").Update(userInfo)
			session.Insert(&expend)
			err = session.Commit()
			if err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"code": -1,
					"msg":"订单处理失败,稍后再试",
				})
				return
			}
		} else {
			affected, err := data.Db.Insert(&order)
			if affected < 1 || err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"code": -1,
					"msg":"订单处理失败,稍后再试",
				})
				return
			}
		}
	}
	// 构造签名
	resp, err := module.WxPay.UnifiedOrder(order.GoodsName, ctx.ClientIP(), order.OrderNO, userInfo.WxOpenID, order.Total)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "微信支付请求失败",
		})
		return
	}
	resp["orderNo"] = order.OrderNO
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"info": "success",
		"data": resp,
	})
}

func ResetOrder (ctx *gin.Context) {
	u, _ := ctx.Get("uid")
	uid := u.(int64)
	orderNo := ctx.Param("orderNo")
	order := &data.OrderLog{}
	has, err := data.Db.Where("order_no = ?", orderNo).Get(order)
	if !has || err != nil {
		ctx.JSON(http.StatusOK, errors.ERROR)
		return
	}
	if uid != order.UID {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "订单不是你的",
		})
		return
	}
	if order.IsPay == 1 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "订单已经支付成功",
		})
		return
	}
	if order.IsCoupon > 0 {
		coupon := &data.CouponLog{
			LogID:order.IsCoupon,
			IsUsage:0,
		}
		data.Db.Where("log_id = ?", order.IsCoupon).Cols("is_usage").Update(coupon)
	}
	if order.Points > 0 {
		has, info := module.User.GetInfoByUID(order.UID)
		if has {
			info.PointsBalance += order.Points
			info.PointsEarning += order.Points
			data.Db.Where("uid = ?", info.UID).Update(info)
			earn := &data.PointsLog{
				UID: info.UID,
				RelationUID:order.UID,
				RelationLogID:order.OrderID,
				FriendlyIntro:"取消订单退回积分",
				Total:order.Points,
				Type:1,

			}
			data.Db.InsertOne(earn)
		}
	}
	ctx.JSON(http.StatusOK, errors.SUCCESS)
}