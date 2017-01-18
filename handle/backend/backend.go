package backend

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/blackdog1987/annual-card/lib/data"
	"github.com/blackdog1987/annual-card/lib/errors"
	"github.com/blackdog1987/annual-card/module"
	"github.com/molibei/page"
	"github.com/molibei/validate"
	"github.com/tealeg/xlsx"
	"crypto/md5"
	"encoding/hex"
)

type loginInput struct {
	Phone  string `form:"phone"`
	Passwd string `form:"passwd"`
}

// Login .
// /v1/manage/login POST
// 管理员登陆
func Login(ctx *gin.Context) {
	var li loginInput
	err := ctx.Bind(&li)
	if err != nil {
		ctx.JSON(http.StatusOK, errors.PARAM_PARSE_ERR)
		return
	}
	if err := validate.Phone(li.Phone); err != nil {
		ctx.JSON(http.StatusOK, errors.PHONE_VALID_ERR)
		return
	}
	info, ok := module.Manager.Login(li.Phone, li.Passwd)
	if !ok {
		ctx.JSON(http.StatusOK, errors.PASSWD_VALID_ERR)
		return
	}
	secret, _ := ctx.Get("secret")
	secretStr := secret.(string)
	token := data.Token{
		UID:       info.ManagerID,
		Identity:  secretStr,
		ExpiredIn: time.Now().Unix() + 7200,
	}
	tk, _ := module.Token.Encode(token, secretStr)
	// rule
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "SUCCESS",
		"data": gin.H{
			"name":       info.Name,
			"token":      tk,
			"expired_in": token.ExpiredIn,
		},
	})
}

// SpreadPlans .
// /v1/manage/spread/plans GET
// 获取推广计划列表
func SpreadPlans(ctx *gin.Context) {
	pstr := ctx.Query("p")
	p, _ := strconv.Atoi(pstr)
	statestr := ctx.Query("state")
	state, _ := strconv.Atoi(statestr)
	total, _ := module.Spread.PlansCount(state)
	offset, length := page.Page(int(total), 10, p, 10)
	var items []data.SpreadPlan
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
	items, _ = module.Spread.Plans(state, offset, length)
	ctx.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": len(items),
		"total": total,
		"p":     p,
		"data":  items,
	})
}

type spreadPlanInput struct {
	ID             int64  `form:"id"`
	Name           string `form:"name"`
	Channel        string `form:"channel"`
	Contact        string `form:"contact"`
	RegCommission  int64  `form:"reg_commission"`
	SaleCommission int64  `form:"sale_commission"`
}

// SaveSpreadPlan .
// /v1/manage/spread/plan POST | PUT
// 新增| 修改 推广计划
func SaveSpreadPlan(ctx *gin.Context) {
	var spi spreadPlanInput
	err := ctx.Bind(&spi)
	if err != nil {
		ctx.JSON(http.StatusOK, errors.PARAM_PARSE_ERR)
		return
	}
	if len(spi.Name) < 1 {
		ctx.JSON(http.StatusOK, errors.SPREAD_PLAN_NAME_EMPTY)
		return
	}
	// check spread plan name repeat
	if has := module.Spread.IsPlanNameRepeat(spi.ID, spi.Name); has {
		ctx.JSON(http.StatusOK, errors.SPREAD_PLAN_NAME_REPEAT)
		return
	}
	if len(spi.Channel) < 1 {
		ctx.JSON(http.StatusOK, errors.SPREAD_CHANNEL_EMPTY)
		return
	}
	if len(spi.Contact) < 1 {
		ctx.JSON(http.StatusOK, errors.CONTACT_EMPTY)
		return
	}
	plan := &data.SpreadPlan{
		Name:           spi.Name,
		Channel:        spi.Channel,
		Contact:        spi.Contact,
		RegCommission:  spi.RegCommission,
		SaleCommission: spi.SaleCommission,
	}
	var dberr error
	if spi.ID > 0 {
		_, dberr = data.Db.Where("sp_id = ?", spi.ID).Update(plan)
	} else {
		_, dberr = data.Db.Insert(plan)
	}
	if dberr != nil {
		ctx.JSON(http.StatusOK, errors.ERROR)
	} else {
		ctx.JSON(http.StatusOK, errors.SUCCESS)
	}
}

func DownXls(ctx *gin.Context) {
	fi := ctx.Param("fi")
	ctx.File("./xlsx/" + fi)
}

// SpreadPlan .
// /v1/manage/spread/plan GET
// 获取单个推广计划
func SpreadPlan(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, _ := strconv.ParseInt(idstr, 10, 64)
	if id < 1 {
		ctx.JSON(http.StatusOK, errors.ID_EMPTY)
		return
	}
	var info = &data.SpreadPlan{}
	has, _ := data.Db.Where("sp_id = ?", id).Get(info)
	if !has {
		ctx.JSON(http.StatusOK, errors.ID_EMPTY)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "SUCCESS",
		"data": info,
	})
}

// SpreadPlanLogs .
// /v1/manage/spread/logs GET
// 获取推广明细
func SpreadPlanLogs(ctx *gin.Context) {
	pstr := ctx.Query("p")
	p, _ := strconv.Atoi(pstr)
	idstr := ctx.Param("id")
	planID, _ := strconv.ParseInt(idstr, 10, 64)
	total, _ := module.Spread.LogsCount(planID, 0, 0)
	offset, length := page.Page(int(total), 10, p, 10)
	var items []data.SpreadLogItem
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
	items, _ = module.Spread.Logs(planID, 0, 0, offset, length)
	ctx.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": len(items),
		"total": total,
		"p":     p,
		"data":  items,
	})
}
// 推广明细导出xls
func SpreadPlanLogs2xls(ctx *gin.Context) {
	idstr := ctx.Param("id")
	planID, _ := strconv.ParseInt(idstr, 10, 64)
	startstr := strings.Trim(ctx.Query("start"), " ")
	endstr := strings.Trim(ctx.Query("end"), " ")
	var (
		start, end int64
	)
	if startstr != "" {
		tm, _ := time.Parse("2006-01-02", startstr)
		start = tm.Unix()
	}
	if endstr != "" {
		tm, _ := time.Parse("2006-01-02", endstr)
		end = tm.Unix()
	}
	total, _ := module.Spread.LogsCount(planID, start, end)
	if total < 1 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "没有数据需要导出",
		})
		return
	}
	items, _ := module.Spread.Logs(planID, start, end, 0, int(total))
	var (
		file                                                                                            *xlsx.File
		sheet                                                                                           *xlsx.Sheet
		row                                                                                             *xlsx.Row
		created, nickname, realname, phone, isBuy, buyTotal, saleComm, regComm *xlsx.Cell
		err error
	)

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("推广明细")
	if err != nil {
		fmt.Printf(err.Error())
	}
	row = sheet.AddRow()
	// 表头
	created = row.AddCell()
	created.Value = "发生时间"
	nickname = row.AddCell()
	nickname.Value = "微信昵称"
	realname = row.AddCell()
	realname.Value = "姓名"
	phone = row.AddCell()
	phone.Value = "手机号"
	isBuy = row.AddCell()
	isBuy.Value = "是否购买"
	buyTotal = row.AddCell()
	buyTotal.Value = "购买金额"
	regComm = row.AddCell()
	regComm.Value = "推广返佣"
	saleComm = row.AddCell()
	saleComm.Value = "销售返佣"
	l := len(items)
	for i := 0; i < l; i++ {
		rows := sheet.AddRow()
		createdVal := rows.AddCell()
		createdVal.Value = time.Unix(items[i].Created, 0).Format("2006-01-02 15:04:05")
		nicknameVal := rows.AddCell()
		nicknameVal.Value = items[i].Consumer.Nickname
		realnameVal := rows.AddCell()
		realnameVal.Value = items[i].Consumer.RealName
		phoneVal := rows.AddCell()
		phoneVal.Value = items[i].Consumer.Phone
		isBuyVal := rows.AddCell()
		if items[i].Category == 1 {
			isBuyVal.Value = "注册"
		} else {
			isBuyVal.Value = "购买"
		}
		buyTotalVal := rows.AddCell()
		buyTotalVal.Value = fmt.Sprintf("%8.2f", float64(items[i].OrderTotal) * 0.01)
		regCommVal := rows.AddCell()
		saleCommVal := rows.AddCell()
		if items[i].Category == 1 {
			regCommVal.Value = fmt.Sprintf("%8.2f", float64(items[i].Commission) * 0.01)
			saleCommVal.Value = "0"
		} else {
			regCommVal.Value = "0"
			saleCommVal.Value = fmt.Sprintf("%8.2f", float64(items[i].Commission) * 0.01)
		}
	}
	fi := fmt.Sprintf("推广明细-[%s-%s].xlsx", time.Unix(start, 0).Format("20060102"), time.Unix(end, 0).Format("20060102"))
	err = file.Save("./xlsx/" + fi)
	if err != nil {
		ctx.JSON(http.StatusOK, errors.ERROR)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": fi,
	})
}

// SpreadPlanQrcode .
// /v1/manage/spread/qrcode GET
// 获取推广二维码
func SpreadPlanQrcode(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id < 1 {
		ctx.JSON(http.StatusOK, errors.ID_EMPTY)
		return
	}
	_, has := module.Spread.Plan(id)
	if !has {
		ctx.JSON(http.StatusOK, errors.ID_EMPTY)
		return
	}
	qrcode, err := module.Spread.QRCode(id, "channel")
	fmt.Println(err)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取二维码失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": "data:image/jpeg;base64," + qrcode,
	})
}

// SpreadPlanState .
// /v1/manage/spread/state PUT
// 关闭推广计划
func SpreadPlanState(ctx *gin.Context) {
	stateStr := ctx.PostForm("state")
	state, _ := strconv.Atoi(stateStr)
	idstr := ctx.Param("id")
	planID, _ := strconv.ParseInt(idstr, 10, 64)
	plan, has := module.Spread.Plan(planID)
	if !has {
		ctx.JSON(http.StatusOK, errors.ID_EMPTY)
		return
	}
	plan.IsDisabled = state
	if ok := module.Spread.UpdatePlan(plan); !ok {
		ctx.JSON(http.StatusOK, errors.ERROR)
		return
	}
	ctx.JSON(http.StatusOK, errors.SUCCESS)
}

// AnnualCard .
// /v1/manage/goods/annual GET
// 获取年卡配置
func AnnualCard(ctx *gin.Context) {
	var out data.AnnualCardConf
	if has := module.Conf.GetObject("ANNUAL_CARD", &out); !has {
		ctx.JSON(http.StatusOK, errors.ID_EMPTY)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "SUCCESS",
		"data": out,
	})
}

// Coupon .
// /v1/manage/goods/coupon GET
// 获取优惠券配置
func Coupon(ctx *gin.Context) {
	var out data.CouponConf
	if has := module.Conf.GetObject("COUPON", &out); !has {
		ctx.JSON(http.StatusOK, errors.ID_EMPTY)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "SUCCESS",
		"data": out,
	})
}

func ActiveHelp(ctx *gin.Context) {
	if has, activeHelp := module.Conf.Get("ACTIVE_HELP"); has {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg": "SUCCESS",
			"data": activeHelp,
		})
		return
	}
	ctx.JSON(http.StatusOK, errors.ID_EMPTY)
}

func SaveActiveHelp(ctx *gin.Context) {
	active_help:= ctx.PostForm("active_help")
	ok := module.Conf.Set("ACTIVE_HELP", active_help)
	if ok {
		ctx.JSON(http.StatusOK, errors.SUCCESS)
		return
	}
	ctx.JSON(http.StatusOK, errors.ERROR)
}


// SaveAnnual .
// /v1/manage/goods/annual POST
// 修改年卡
func SaveAnnual(ctx *gin.Context) {
	var in data.AnnualCardConf
	err := ctx.Bind(&in)
	if err != nil {
		ctx.JSON(http.StatusOK, errors.PARAM_PARSE_ERR)
		return
	}
	if len(in.Name) < 1 {
		ctx.JSON(http.StatusOK, errors.NAME_EMPTY)
		return
	}
	if in.SalePrice < 1 {
		ctx.JSON(http.StatusOK, errors.SALE_PRICE_ZERO)
		return
	}
	ok := module.Conf.SetObject("ANNUAL_CARD", in)
	if ok {
		ctx.JSON(http.StatusOK, errors.SUCCESS)
		return
	}
	ctx.JSON(http.StatusOK, errors.ERROR)
}

// SaveCoupon .
// /v1/manage/goods/coupon POST
// 修改优惠券
func SaveCoupon(ctx *gin.Context) {
	var in data.CouponConf
	err := ctx.Bind(&in)
	if err != nil {
		ctx.JSON(http.StatusOK, errors.PARAM_PARSE_ERR)
		return
	}
	if len(in.Name) < 1 {
		ctx.JSON(http.StatusOK, errors.NAME_EMPTY)
		return
	}
	if in.SalePrice < 1 {
		ctx.JSON(http.StatusOK, errors.SALE_PRICE_ZERO)
		return
	}
	ok := module.Conf.SetObject("COUPON", in)
	if ok {
		ctx.JSON(http.StatusOK, errors.SUCCESS)
		return
	}
	ctx.JSON(http.StatusOK, errors.ERROR)
}

type uploadInput struct {
	CodeType string `form:"code_type"`
	Body     string `form:"body"`
}

// UploadImage .
// /v1/manage/upload/image POST
// 上传图片
func UploadImage(ctx *gin.Context) {
	var (
		up uploadInput
		mime = map[string]string{
			"image/png;":  ".png",
			"image/bmp;":  ".bmp",
			"image/jpeg;": ".jpg",
		}
	)
	if err := ctx.Bind(&up); err != nil {
		fmt.Println("upload image failed:", err.Error())
		ctx.JSON(http.StatusOK, errors.UPLOAD_FAILED)
		return
	}
	// 保存到阿里云
	body := strings.Split(up.Body, "base64,")
	data, _ := base64.StdEncoding.DecodeString(body[1])
	// 取扩展名
	ext := strings.Split(body[0], ":")
	n := time.Now().UnixNano()
	fn := fmt.Sprintf("media/%d%s", n, mime[ext[1]])
	err := ioutil.WriteFile("./" + fn, data, 0666)
	if err != nil {
		fmt.Println("write failed:", err.Error())
		ctx.JSON(http.StatusOK, errors.UPLOAD_FAILED)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"info": "success",
		"data": fn,
	})
}

// Orders .
// /v1/manage/orders GET
// 订单列表
func Orders(ctx *gin.Context) {
	pstr := ctx.Query("p")
	p, _ := strconv.Atoi(pstr)
	startstr := strings.Trim(ctx.Query("start"), " ")
	endstr := strings.Trim(ctx.Query("end"), " ")
	var (
		start, end int64
	)
	if startstr != "" {
		tm, _ := time.Parse("2006-01-02", startstr)
		start = tm.Unix()
	}
	if endstr != "" {
		tm, _ := time.Parse("2006-01-02", endstr)
		end = tm.Unix()
	}
	total, _ := module.Order.Count(start, end)
	offset, length := page.Page(int(total), 10, p, 10)
	var items []data.OrderLog
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
	items, _ = module.Order.Search(start, end, offset, length)
	ctx.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": len(items),
		"total": total,
		"p":     p,
		"data":  items,
	})
}
// 导出订单列表
func OrdersExport(ctx *gin.Context) {
	startstr := strings.Trim(ctx.Query("start"), " ")
	endstr := strings.Trim(ctx.Query("end"), " ")
	var (
		start, end int64
	)
	if startstr != "" {
		tm, _ := time.Parse("2006-01-02", startstr)
		start = tm.Unix()
	}
	if endstr != "" {
		tm, _ := time.Parse("2006-01-02", endstr)
		end = tm.Unix()
	}
	total, _ := module.Order.Count(start, end)
	if total < 1 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":"没有要导出的数据",
		})
		return
	}
	var items []data.OrderLog
	items, _ = module.Order.Search(start, end, 0, int(total))
	var (
		file                                                                                            *xlsx.File
		sheet                                                                                           *xlsx.Sheet
		row                                                                                             *xlsx.Row
		orderNo, goodsName, goodsNum, price, points, coupon, totalFee  *xlsx.Cell
		err error
	)

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("实体卡明细")
	if err != nil {
		fmt.Printf(err.Error())
	}
	row = sheet.AddRow()
	// 表头
	orderNo = row.AddCell()
	orderNo.Value = "订单号"
	goodsName = row.AddCell()
	goodsName.Value = "产品名称"
	goodsNum = row.AddCell()
	goodsNum.Value = "商品数量"
	price = row.AddCell()
	price.Value = "订单价格"
	points = row.AddCell()
	points.Value = "积分抵扣"
	coupon = row.AddCell()
	coupon.Value = "优惠券抵扣"
	totalFee = row.AddCell()
	totalFee.Value = "实收金额"
	l := len(items)
	for i := 0; i < l; i++ {
		rows := sheet.AddRow()
		orderNoVal := rows.AddCell()
		orderNoVal.Value = items[i].OrderNO
		goodsNameVal := rows.AddCell()
		goodsNameVal.Value = items[i].GoodsName
		goodsNumVal := rows.AddCell()
		goodsNumVal.Value = fmt.Sprintf("%d", items[i].GoodsNum)
		priceVal := rows.AddCell()
		priceVal.Value = fmt.Sprintf("%8.2f", float64(items[i].Price) * 0.01)
		pointsVal := rows.AddCell()
		pointsVal.Value = fmt.Sprintf("%d", items[i].Points)
		couponVal := rows.AddCell()
		couponVal.Value = fmt.Sprintf("%8.2f", float64(items[i].CouponPrice) * 0.01)
		totalFeeVal := rows.AddCell()
		totalFeeVal.Value = fmt.Sprintf("%8.2f", float64(items[i].Total) * 0.01)
	}
	fi := fmt.Sprintf("订单明细-[%s-%s].xlsx", time.Unix(start, 0).Format("20060102"), time.Unix(end, 0).Format("20060102"))
	err = file.Save("./xlsx/" + fi)
	if err != nil {
		ctx.JSON(http.StatusOK, errors.ERROR)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": fi,
	})
}
// Members .
// /v1/manage/members GET
// 年卡会员列表
func Members(ctx *gin.Context) {
	pstr := ctx.Query("p")
	p, _ := strconv.Atoi(pstr)
	name:= ctx.Query("name")
	idcard:=ctx.Query("idcard")
	startstr := strings.Trim(ctx.Query("start"), " ")
	endstr := strings.Trim(ctx.Query("end"), " ")
	var (
		start, end int64
	)
	if startstr != "" {
		tm, _ := time.Parse("2006-01-02", startstr)
		start = tm.Unix()
	}
	if endstr != "" {
		tm, _ := time.Parse("2006-01-02", endstr)
		end = tm.Unix()
	}
	total, _ := module.AnnualCard.Count(1, 0, start, end, name, idcard)
	offset, length := page.Page(int(total), 10, p, 10)
	var items []data.AnnualCard
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
	items, _ = module.AnnualCard.Search(1, 0, start, end,name,idcard, offset, length)
	ctx.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": len(items),
		"total": total,
		"p":     p,
		"data":  items,
	})
}
// 删除年卡
func DeleteCard(ctx *gin.Context) {
	card_id := ctx.Param("card_id")
	cardID, _ := strconv.ParseInt(card_id, 10, 64)
	info, has := module.AnnualCard.UsagesByCardID(cardID)
	if !has {
		ctx.JSON(http.StatusOK, errors.CARD_NOT_FOUND)
		return
	}
	if info.UsageNum > 0 {
		ctx.JSON(http.StatusOK, errors.CARD_USAGED)
		return
	}
	card := data.AnnualCard{
		IsDelete:1,
	}
	affected, err :=data.Db.Where("card_id = ?", cardID).Update(&card)
	if affected < 1 || err != nil {
		ctx.JSON(http.StatusOK, errors.ERROR)
		return
	}
	ctx.JSON(http.StatusOK, errors.SUCCESS)
}

func UnActive(ctx *gin.Context) {
	card_id := ctx.Param("card_id")
	cardID, _ := strconv.ParseInt(card_id, 10, 64)
	info, has := module.AnnualCard.UsagesByCardID(cardID)
	if !has {
		ctx.JSON(http.StatusOK, errors.CARD_NOT_FOUND)
		return
	}
	if info.UsageNum > 0 {
		ctx.JSON(http.StatusOK, errors.CARD_USAGED)
		return
	}
	now := time.Now().Unix()
	card := data.AnnualCard{
		IsActive:0,
		BindContact: "",
		BindHeadimg: "",
		BindIDCard: "",
		BindName: "",
		ExpiredStart:now,
		ExpiredStop: now + 31536000,
	}
	affected, err :=data.Db.Where("card_id = ?", cardID).Cols("is_active,bind_contact,bind_headimg,bind_idcard,bind_name,expired_start, expired_stop").Update(&card)
	if affected < 1 || err != nil {
		ctx.JSON(http.StatusOK, errors.ERROR)
		return
	}
	ctx.JSON(http.StatusOK, errors.SUCCESS)
}

// 导出年卡会员明细
func MembersExport(ctx *gin.Context) {
	startstr := strings.Trim(ctx.Query("start"), " ")
	endstr := strings.Trim(ctx.Query("end"), " ")
	name:= ctx.Query("name")
	idcard:=ctx.Query("idcard")
	var (
		start, end int64
	)
	if startstr != "" {
		tm, _ := time.Parse("2006-01-02", startstr)
		start = tm.Unix()
	}
	if endstr != "" {
		tm, _ := time.Parse("2006-01-02", endstr)
		end = tm.Unix()
	}
	total, _ := module.AnnualCard.Count(1, 0, start, end, name, idcard)
	if total < 1 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "没有要导出的数据",
		})
		return
	}
	var items []data.AnnualCard
	items, _ = module.AnnualCard.Search(1, 0, start, end,name, idcard, 0, int(total))
	var (
		file                                                                                            *xlsx.File
		sheet                                                                                           *xlsx.Sheet
		row                                                                                             *xlsx.Row
		headImg, realname, phone, expired_in, idcardCell  *xlsx.Cell
		err error
	)

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("实体卡明细")
	if err != nil {
		fmt.Printf(err.Error())
	}
	row = sheet.AddRow()
	// 表头
	headImg = row.AddCell()
	headImg.Value = "照片"
	realname = row.AddCell()
	realname.Value = "姓名"
	idcardCell = row.AddCell()
	idcardCell.Value = "身份证号"
	phone = row.AddCell()
	phone.Value = "手机号"
	expired_in = row.AddCell()
	expired_in.Value = "有效期"
	l := len(items)
	for i := 0; i < l; i++ {
		rows := sheet.AddRow()
		headImgVal := rows.AddCell()
		headImgVal.Value = data.BaseConf.Server.Host + items[i].BindHeadimg
		realnameVal := rows.AddCell()
		realnameVal.Value = items[i].BindName
		idcardVal := rows.AddCell()
		idcardVal.Value = items[i].BindIDCard
		phoneVal := rows.AddCell()
		phoneVal.Value = items[i].BindContact
		expiredInVal := rows.AddCell()
		expiredInVal.Value = time.Unix(items[i].ExpiredStart, 0).Format("20060102") + "-" + time.Unix(items[i].ExpiredStop, 0).Format("20060102")
	}
	fi := fmt.Sprintf("年卡会员明细-[%s-%s].xlsx", time.Unix(start, 0).Format("20060102"), time.Unix(end, 0).Format("20060102"))
	err = file.Save("./xlsx/" + fi)
	if err != nil {
		ctx.JSON(http.StatusOK, errors.ERROR)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": fi,
	})
}
// AddAnnualCardPlan .
// /v1/manage/card/plan POST
// 新建年卡推广计划
func AddAnnualCardPlan(ctx *gin.Context) {
	var in data.AnnualCardPlan
	err := ctx.Bind(&in)
	if err != nil {
		ctx.JSON(http.StatusOK, errors.PARAM_PARSE_ERR)
		return
	}
	if len(in.Channel) < 1 {
		ctx.JSON(http.StatusOK, errors.CHANNEL_EMPTY)
		return
	}
	if in.ExpiredStart < 1 {
		ctx.JSON(http.StatusOK, errors.EXPIRED_START_ZERO)
		return
	}
	if in.ExpiredStop < 1 {
		ctx.JSON(http.StatusOK, errors.EXPIRED_STOP_ZERO)
		return
	}
	if in.CreateNum < 1 {
		ctx.JSON(http.StatusOK, errors.CREATE_NUM_ZERO)
		return
	}
	if in.CreateNum > 500 {
		ctx.JSON(http.StatusOK, errors.CREATE_NUM_MAX)
	}
	// generate card
	affected, err := data.Db.InsertOne(&in)
	if err != nil {
		ctx.JSON(http.StatusOK, errors.ERROR)
		return
	}
	items := make([]data.AnnualCard, in.CreateNum)
	cardConf := data.AnnualCardConf{}
	module.Conf.GetObject("ANNUAL_CARD", &cardConf)
	for i := 0; i < int(in.CreateNum); i++ {
		items[i].PlanID = in.CPID
		items[i].CardName = cardConf.Name
		items[i].ExpiredStart = in.ExpiredStart
		items[i].ExpiredStop = in.ExpiredStop
		now := time.Now().UnixNano()
		items[i].CardNO = generateCardNo(in.CPID, now)
		items[i].CardPasswd = generateCaptcha(now)
	}
	affected, err = data.Db.Insert(items)
	if err != nil {
		if affected == 0 {
			ctx.JSON(http.StatusOK, errors.ERROR)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  fmt.Sprintf("共有%d创建成功", affected),
		})
		return
	}
	ctx.JSON(http.StatusOK, errors.SUCCESS)
}

func generateCaptcha(now int64) (captcha string) {
	r := rand.New(rand.NewSource(now))
	code := r.Perm(10)
	s := code[0]
	var capt, l []int
	if s < 4 {
		capt = code[s : s + 6]
	} else {
		capt = code[s:]
		l = code[0:(6 - (9 - s))]
		capt = append(capt, l...)
	}
	for i := 0; i < 6; i++ {
		captcha = fmt.Sprintf("%v%v", captcha, capt[i])
	}
	return captcha
}

func generateCardNo(prefix, now int64) (s string) {
	r := rand.New(rand.NewSource(now))
	code := r.Perm(9)
	s = fmt.Sprintf("%d", prefix)
	for i := 0; i < 9; i++ {
		s = fmt.Sprintf("%v%v", s, code[i])
	}
	return
}

// AnnualCardPlans .
// /v1/manage/card/plans GET
// 年卡计划列表
func AnnualCardPlans(ctx *gin.Context) {
	pstr := ctx.Query("p")
	p, _ := strconv.Atoi(pstr)
	total, _ := module.AnnualCard.PlanCount()
	offset, length := page.Page(int(total), 10, p, 10)
	var items []data.AnnualCardPlan
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
	items, _ = module.AnnualCard.PlanSearch(offset, length)
	ctx.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": len(items),
		"total": total,
		"p":     p,
		"data":  items,
	})
}

// AnnualCardDetail .
// /v1/manage/card/detail GET
// 年卡计划明细
func AnnualCardDetail(ctx *gin.Context) {
	pstr := ctx.Query("p")
	p, _ := strconv.Atoi(pstr)
	planIDStr := ctx.Param("id")
	planID, _ := strconv.ParseInt(planIDStr, 10, 64)
	if planID < 1 {
		ctx.JSON(http.StatusOK, errors.CARD_PLAN_ID_ZERO)
		return
	}
	startstr := strings.Trim(ctx.Query("start"), " ")
	endstr := strings.Trim(ctx.Query("end"), " ")
	var (
		start, end int64
	)
	if startstr != "" {
		tm, _ := time.Parse("2006-01-02", startstr)
		start = tm.Unix()
	}
	if endstr != "" {
		tm, _ := time.Parse("2006-01-02", endstr)
		end = tm.Unix()
	}
	total, _ := module.AnnualCard.Count(-1, planID, start, end, "", "")
	offset, length := page.Page(int(total), 10, p, 10)
	var items []data.AnnualCard
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
	items, _ = module.AnnualCard.Search(-1, planID, start, end, "", "", offset, length)
	ctx.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": len(items),
		"total": total,
		"p":     p,
		"data":  items,
	})
}

func AnnualCardDetailExport(ctx *gin.Context) {
	planIDStr := ctx.Param("id")
	planID, _ := strconv.ParseInt(planIDStr, 10, 64)
	if planID < 1 {
		ctx.JSON(http.StatusOK, errors.CARD_PLAN_ID_ZERO)
		return
	}
	startstr := strings.Trim(ctx.Query("start"), " ")
	endstr := strings.Trim(ctx.Query("end"), " ")
	var (
		start, end int64
	)
	if startstr != "" {
		tm, _ := time.Parse("2006-01-02", startstr)
		start = tm.Unix()
	}
	if endstr != "" {
		tm, _ := time.Parse("2006-01-02", endstr)
		end = tm.Unix()
	}
	total, _ := module.AnnualCard.Count(-1, planID, start, end, "", "")
	if total < 1 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "没有需要导出的数据",
		})
		return
	}
	var items []data.AnnualCard
	items, _ = module.AnnualCard.Search(-1, planID, start, end, "", "", 0, int(total))
	var (
		file                                                                                            *xlsx.File
		sheet                                                                                           *xlsx.Sheet
		row                                                                                             *xlsx.Row
		number, pass, realname, phone,idcard  *xlsx.Cell
		err error
	)

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("实体卡明细")
	if err != nil {
		fmt.Printf(err.Error())
	}
	row = sheet.AddRow()
	// 表头
	number = row.AddCell()
	number.Value = "卡号"
	pass = row.AddCell()
	pass.Value = "卡密"
	realname = row.AddCell()
	realname.Value = "姓名"
	idcard = row.AddCell()
	idcard.Value = "身份证号"
	phone = row.AddCell()
	phone.Value = "手机号"
	l := len(items)
	for i := 0; i < l; i++ {
		rows := sheet.AddRow()
		numberVal := rows.AddCell()
		numberVal.Value = items[i].CardNO
		passVal := rows.AddCell()
		passVal.Value = items[i].CardPasswd
		realnameVal := rows.AddCell()
		realnameVal.Value = items[i].BindName
		idcardVal := rows.AddCell()
		idcardVal.Value = items[i].BindIDCard
		phoneVal := rows.AddCell()
		phoneVal.Value = items[i].BindContact
	}
	fi := fmt.Sprintf("实体卡明细-%d-[%s-%s].xlsx", planID, time.Unix(start, 0).Format("20060102"), time.Unix(end, 0).Format("20060102"))
	err = file.Save("./xlsx/" + fi)
	if err != nil {
		ctx.JSON(http.StatusOK, errors.ERROR)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": fi,
	})
}

// Merchants .
// /v1/manage/merchants GET
// 获取商户列表
func Merchants(ctx *gin.Context) {
	pstr := ctx.Query("p")
	p, _ := strconv.Atoi(pstr)
	keyword := ctx.Query("keyword")
	total, _ := module.Merchant.Count(keyword)
	offset, length := page.Page(int(total), 10, p, 10)
	var items []data.Merchant
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
	items, _ = module.Merchant.Search(keyword, offset, length)
	ctx.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": len(items),
		"total": total,
		"p":     p,
		"data":  items,
	})
}

// Merchant .
// /v1/manage/merchant GET
// 获取单个商户信息
func Merchant(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, _ := strconv.ParseInt(idstr, 10, 64)
	if id < 1 {
		ctx.JSON(http.StatusOK, errors.ID_EMPTY)
		return
	}
	has, info := module.Merchant.Get(id)
	if !has {
		ctx.JSON(http.StatusOK, errors.ID_EMPTY)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "SUCCESS",
		"data": info,
	})
}

// SaveMerchant .
// /v1/manage/merchant POST
// 新增 | 编辑商户信息
func SaveMerchant(ctx *gin.Context) {
	var in data.Merchant
	err := ctx.Bind(&in)
	if err != nil {
		ctx.JSON(http.StatusOK, errors.PARAM_PARSE_ERR)
		return
	}
	if len(in.MchName) < 1 {
		ctx.JSON(http.StatusOK, errors.NAME_EMPTY)
		return
	}
	// repeat
	if module.Merchant.IsMchNameRepeat(in.MchID, in.MchName) {
		ctx.JSON(http.StatusOK, errors.NAME_REPEAT)
		return
	}
	if in.MchID > 0 {
		_, err = data.Db.Where("mch_id = ?", in.MchID).Update(&in)
	} else {
		_, err = data.Db.Insert(&in)
	}
	if err != nil {
		ctx.JSON(http.StatusOK, errors.ERROR)
		return
	}
	ctx.JSON(http.StatusOK, errors.SUCCESS)
}

// DbConfig .
// /v1/manage/config/db GET
// 获取配置 - 数据库
func DbConfig(ctx *gin.Context) {
	has, reg := module.Conf.Get("SPREAD_REG")
	if !has {
		reg = "0"
	}
	has, sale := module.Conf.Get("SPREAD_SALE")
	if !has {
		sale = "0"
	}
	has, rate := module.Conf.Get("INTEGRAL_RATE")
	if !has {
		rate = "0"
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "SUCCESS",
		"data": gin.H{
			"spread_reg":    reg,
			"spread_sale":   sale,
			"integral_rate": rate,
		},
	})
}

// SaveDbConfig .
// /v1/manage/config/db POST
// 保存配置 - 数据库
func SaveDbConfig(ctx *gin.Context) {
	reg := ctx.PostForm("spread_reg")
	sale := ctx.PostForm("spread_sale")
	rate := ctx.PostForm("integral_rate")
	module.Conf.Set("INTEGRAL_RATE", rate)
	module.Conf.Set("SPREAD_REG", reg)
	module.Conf.Set("SPREAD_SALE", sale)
	ctx.JSON(http.StatusOK, errors.SUCCESS)
}

// MerchantAccounts .
// /v1/manage/merchant-accounts GET
// 获取商户账户列表
func MerchantAccounts(ctx *gin.Context) {
	pstr := ctx.Query("p")
	p, _ := strconv.Atoi(pstr)
	keyword := ctx.Query("keyword")
	total, _ := module.MerchantAccount.Count(keyword)
	offset, length := page.Page(int(total), 10, p, 10)
	var items []data.MerchantAccount
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
	items, _ = module.MerchantAccount.Search(keyword, offset, length)
	ctx.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": len(items),
		"total": total,
		"p":     p,
		"data":  items,
	})
}

// MerchantAccount .
// /v1/manage/merchant-account GET
// 获取单个商户信息
func MerchantAccount(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, _ := strconv.ParseInt(idstr, 10, 64)
	if id < 1 {
		ctx.JSON(http.StatusOK, errors.ID_EMPTY)
		return
	}
	has, info := module.MerchantAccount.Get(id)
	if !has {
		ctx.JSON(http.StatusOK, errors.ID_EMPTY)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "SUCCESS",
		"data": info,
	})
}

// SaveMerchantAccount .
// /v1/manage/merchant-account POST
// 新增 | 编辑商户信息
func SaveMerchantAccount(ctx *gin.Context) {
	var in data.MerchantAccount
	err := ctx.Bind(&in)
	if err != nil {
		ctx.JSON(http.StatusOK, errors.PARAM_PARSE_ERR)
		return
	}
	if len(in.Name) < 1 {
		ctx.JSON(http.StatusOK, errors.NAME_EMPTY)
		return
	}
	if len(in.Account) < 1 {
		ctx.JSON(http.StatusOK, errors.ACCOUNT_EMPTY)
		return
	}
	if module.MerchantAccount.IsMchAccountRepeat(in.MchID, in.Account) {
		ctx.JSON(http.StatusOK, errors.ACCOUNT_REPEAT)
		return
	}
	if module.MerchantStore.IsAccountRepeat(0, in.Account) {
		ctx.JSON(http.StatusOK, errors.ACCOUNT_REPEAT)
		return
	}
	if len(in.Passwd) > 0 {
		md5Ctx := md5.New()
		md5Ctx.Write([]byte(in.Passwd))
		cipherStr := md5Ctx.Sum(nil)
		in.Passwd = hex.EncodeToString(cipherStr)
	}
	if in.MchID > 0 {
		_, err = data.Db.Where("mch_id = ?", in.MchID).Update(&in)
	} else {
		if len(in.Passwd) < 1 {
			ctx.JSON(http.StatusOK, errors.PASSWD_EMPTY)
			return
		}
		_, err = data.Db.Insert(&in)
	}
	if err != nil {
		ctx.JSON(http.StatusOK, errors.ERROR)
		return
	}
	ctx.JSON(http.StatusOK, errors.SUCCESS)
}

//  ToggleMerchantAccountState
// /v1/manage/merchant-account PUT
// 禁用|启用 商户账号

func ToggleMerchantAccountState(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, _ := strconv.ParseInt(idstr, 10, 64)
	stateStr := ctx.PostForm("state")
	state, _:= strconv.Atoi(stateStr)
	if id < 1 {
		ctx.JSON(http.StatusOK, errors.ID_EMPTY)
		return
	}
	has, info := module.MerchantAccount.Get(id)
	if !has {
		ctx.JSON(http.StatusOK, errors.ID_EMPTY)
		return
	}
	info.State = state
	_, err := data.Db.Where("mch_id = ?", id).Cols("state").Update(&info)
	if err != nil {
		ctx.JSON(http.StatusOK, errors.ERROR)
		return
	}
	ctx.JSON(http.StatusOK, errors.SUCCESS)
}

func ResetPwd(ctx *gin.Context) {

}