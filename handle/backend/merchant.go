package backend

import (
	"github.com/gin-gonic/gin"
	"time"
	"net/http"
	"github.com/blackdog1987/annual-card/lib/errors"
	"github.com/blackdog1987/annual-card/module"
	"github.com/blackdog1987/annual-card/lib/data"
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"github.com/molibei/page"
	"strings"
	"github.com/tealeg/xlsx"
	"fmt"
)

type merchantLoginInput struct {
	Account string `form:"account"`
	Passwd  string `form:"passwd"`
}

// 商户|门店 登录
func MerchantLogin(ctx *gin.Context) {
	var li merchantLoginInput
	err := ctx.Bind(&li)
	if err != nil {
		ctx.JSON(http.StatusOK, errors.PARAM_PARSE_ERR)
		return
	}
	var (
		id int64
		isStore bool
		name, pwd string
	)
	md5Ctx := md5.New()

	md5Ctx.Write([]byte(li.Passwd))
	cipherStr := md5Ctx.Sum(nil)
	pwd = hex.EncodeToString(cipherStr)

	if has, info := module.MerchantAccount.GetByAccount(li.Account); has {
		if pwd == info.Passwd {
			id = info.MchID
		} else {
			ctx.JSON(http.StatusOK, errors.PASSWD_VALID_ERR)
			return
		}
		if info.State == -1 {
			ctx.JSON(http.StatusOK, errors.ACCOUNT_DISABLED)
			return
		}
		name = info.Name
	} else {
		has, info := module.MerchantStore.GetByAccount(li.Account)
		if !has {
			ctx.JSON(http.StatusOK, errors.PASSWD_VALID_ERR)
			return
		}
		if pwd != info.Passwd {
			ctx.JSON(http.StatusOK, errors.PASSWD_VALID_ERR)
			return
		}
		if info.State == -1 {
			ctx.JSON(http.StatusOK, errors.ACCOUNT_DISABLED)
			return
		}
		id = info.StoreID
		isStore = true
		name = info.StoreName
	}
	secret, _ := ctx.Get("secret")
	secretStr := secret.(string)
	token := data.Token{
		UID:       id,
		IsStore:   isStore,
		Identity:  secretStr,
		ExpiredIn: time.Now().Unix() + 7200,
	}
	tk, _ := module.Token.Encode(token, secretStr)
	// rule
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "SUCCESS",
		"data": gin.H{
			"is_store": isStore,
			"token":      tk,
			"name": name,
			"expired_in": token.ExpiredIn,
		},
	})
}

// 修改密码
func ResetMerchantAccountPwd(ctx *gin.Context) {
	tks, _ := ctx.Get("token")
	tk := tks.(data.Token)
	old := ctx.PostForm("old_passwd")
	if len(old) < 1 {
		ctx.JSON(http.StatusOK, errors.PASSWD_VALID_ERR)
		return
	}
	newPwd := ctx.PostForm("new_passwd")
	if len(newPwd) < 1 {
		ctx.JSON(http.StatusOK, errors.PASSWD_EMPTY)
		return
	}
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(old))
	cipherStr := md5Ctx.Sum(nil)
	old = hex.EncodeToString(cipherStr)
	md5Ctx.Reset()
	md5Ctx.Write([]byte(newPwd))
	cipherStr = md5Ctx.Sum(nil)
	newPwd = hex.EncodeToString(cipherStr)
	if len(old) < 1 {
		ctx.JSON(http.StatusOK, errors.PASSWD_VALID_ERR)
		return
	}
	if tk.IsStore {
		if has, info := module.MerchantStore.Get(tk.UID); has {
			if old != info.Passwd {
				ctx.JSON(http.StatusOK, errors.PASSWD_VALID_ERR)
				return
			}
			info.Passwd = newPwd
			if _, err := data.Db.Where("store_id = ?", info.StoreID).Update(&info); err != nil {
				ctx.JSON(http.StatusOK, errors.ERROR)
				return
			}
			ctx.JSON(http.StatusOK, errors.SUCCESS)
			return
		}
		ctx.JSON(http.StatusOK, errors.ID_EMPTY)
		return
	} else {
		if has, info := module.MerchantAccount.Get(tk.UID); has {
			if old != info.Passwd {
				ctx.JSON(http.StatusOK, errors.PASSWD_VALID_ERR)
				return
			}
			info.Passwd = newPwd
			if _, err := data.Db.Where("mch_id = ?", info.MchID).Update(&info); err != nil {
				ctx.JSON(http.StatusOK, errors.ERROR)
				return
			}
			ctx.JSON(http.StatusOK, errors.SUCCESS)
			return
		}
		ctx.JSON(http.StatusOK, errors.ID_EMPTY)
		return
	}
}

func Stores(ctx *gin.Context) {
	tks, _ := ctx.Get("token")
	tk := tks.(data.Token)
	if tk.IsStore {
		ctx.JSON(http.StatusOK, errors.AUTHORIFY_VALID_ERR)
		return
	}
	pstr := ctx.Query("p")
	p, _ := strconv.Atoi(pstr)
	keyword := ctx.Query("keyword")
	total, _ := module.MerchantStore.Count(tk.UID, keyword)
	offset, length := page.Page(int(total), 10, p, 10)
	var items []data.MerchantStore
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
	items, _ = module.MerchantStore.Search(tk.UID, keyword, offset, length)
	ctx.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": len(items),
		"total": total,
		"p":     p,
		"data":  items,
	})
}

// 获取门店账号详情
func Store(ctx *gin.Context) {
	tks, _ := ctx.Get("token")
	tk := tks.(data.Token)
	idstr := ctx.Param("id")
	id, _ := strconv.ParseInt(idstr, 10, 64)
	if id < 1 {
		ctx.JSON(http.StatusOK, errors.ID_EMPTY)
		return
	}
	if tk.IsStore && tk.UID != id {
		ctx.JSON(http.StatusOK, errors.AUTHORIFY_VALID_ERR)
		return
	}
	has, info := module.MerchantStore.Get(id)
	if !has {
		ctx.JSON(http.StatusOK, errors.ID_EMPTY)
		return
	}
	if !tk.IsStore {
		if tk.UID != info.MchId {
			ctx.JSON(http.StatusOK, errors.AUTHORIFY_VALID_ERR)
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "SUCCESS",
		"data": info,
	})
}

//  新增|编辑 门店账号
func SaveMerchantStore(ctx *gin.Context) {
	tks, _ := ctx.Get("token")
	tk := tks.(data.Token)
	if tk.IsStore {
		ctx.JSON(http.StatusOK, errors.AUTHORIFY_VALID_ERR)
		return
	}
	var in data.MerchantStore
	err := ctx.Bind(&in)
	if err != nil {
		ctx.JSON(http.StatusOK, errors.PARAM_PARSE_ERR)
		return
	}
	if len(in.StoreName) < 1 {
		ctx.JSON(http.StatusOK, errors.NAME_EMPTY)
		return
	}
	if len(in.Account) < 1 {
		ctx.JSON(http.StatusOK, errors.ACCOUNT_EMPTY)
		return
	}

	if module.MerchantAccount.IsMchAccountRepeat(0, in.Account) {
		ctx.JSON(http.StatusOK, errors.ACCOUNT_REPEAT)
		return
	}
	if module.MerchantStore.IsAccountRepeat(in.StoreID, in.Account) {
		ctx.JSON(http.StatusOK, errors.ACCOUNT_REPEAT)
		return
	}
	// 处理密码
	if len(in.Passwd) > 0 {
		md5Ctx := md5.New()
		md5Ctx.Write([]byte(in.Passwd))
		cipherStr := md5Ctx.Sum(nil)
		in.Passwd = hex.EncodeToString(cipherStr)
	}
	in.MchId = tk.UID
	if in.StoreID > 0 {
		_, err = data.Db.Where("store_id = ?", in.StoreID).Update(&in)
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

// 禁用|启用账号
func ToggleMerchantStoreState(ctx *gin.Context) {
	tks, _ := ctx.Get("token")
	tk := tks.(data.Token)
	if tk.IsStore {
		ctx.JSON(http.StatusOK, errors.AUTHORIFY_VALID_ERR)
		return
	}
	idstr := ctx.Param("id")
	id, _ := strconv.ParseInt(idstr, 10, 64)
	stateStr := ctx.PostForm("state")
	state, _ := strconv.Atoi(stateStr)
	if id < 1 {
		ctx.JSON(http.StatusOK, errors.ID_EMPTY)
		return
	}
	has, info := module.MerchantStore.Get(id)
	if !has {
		ctx.JSON(http.StatusOK, errors.ID_EMPTY)
		return
	}
	if info.MchId != tk.UID {
		ctx.JSON(http.StatusOK, errors.AUTHORIFY_VALID_ERR)
		return
	}
	info.State = state
	_, err := data.Db.Where("store_id = ?", id).Cols("state").Update(&info)
	if err != nil {
		ctx.JSON(http.StatusOK, errors.ERROR)
		return
	}
	ctx.JSON(http.StatusOK, errors.SUCCESS)
}

// 本商户|门店 年卡使用记录
func CardUsages(ctx *gin.Context) {
	tks, _ := ctx.Get("token")
	tk := tks.(data.Token)
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
	var (
		mchID, storeID int64
	)
	if tk.IsStore {
		storeID = tk.UID
		_, info := module.MerchantStore.Get(storeID)
		mchID = info.MchId
	} else {
		mchID = tk.UID
	}
	total, _ := module.CardUsage.Count(mchID, storeID, start, end, tk.IsStore)
	offset, length := page.Page(int(total), 10, p, 10)
	var items []data.AnnualCardUsageLog
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
	items, _ = module.CardUsage.Search(mchID, storeID, start, end, tk.IsStore, offset, length)
	ctx.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": len(items),
		"total": total,
		"p":     p,
		"data":  items,
	})
}
// 本商户|门店 年卡使用记录
func CardUsagesExport(ctx *gin.Context) {
	tks, _ := ctx.Get("token")
	tk := tks.(data.Token)
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
	var (
		mchID, storeID int64
	)
	if tk.IsStore {
		storeID = tk.UID
		_, info := module.MerchantStore.Get(storeID)
		mchID = info.MchId
	} else {
		mchID = tk.UID
	}
	total, _ := module.CardUsage.Count(mchID, storeID, start, end, tk.IsStore)
	if total < 1 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "没有需要导出的数据",
		})
		return
	}
	var items []data.AnnualCardUsageLog
	items, _ = module.CardUsage.Search(mchID, storeID, start, end, tk.IsStore, 0, int(total))
	var (
		file                                                                                            *xlsx.File
		sheet                                                                                           *xlsx.Sheet
		row                                                                                             *xlsx.Row
		pic, name, phone, idcard, storeName, usageTime  *xlsx.Cell
		err error
	)

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("实体卡明细")
	if err != nil {
		fmt.Printf(err.Error())
	}
	row = sheet.AddRow()
	// 表头
	pic = row.AddCell()
	pic.Value = "照片"
	name = row.AddCell()
	name.Value = "姓名"
	phone = row.AddCell()
	phone.Value = "手机号"
	idcard = row.AddCell()
	idcard.Value = "身份证"
	storeName = row.AddCell()
	storeName.Value = "核销门店"
	usageTime = row.AddCell()
	usageTime.Value = "核销时间"
	l := len(items)
	for i := 0; i < l; i++ {
		rows := sheet.AddRow()
		picVal := rows.AddCell()
		picVal.Value = "http://wx.qinzinianka.com" + items[i].BindHeadimg
		nameVal := rows.AddCell()
		nameVal.Value = items[i].BindName
		phoneVal := rows.AddCell()
		phoneVal.Value = items[i].BindContact
		idcardVal := rows.AddCell()
		idcardVal.Value = items[i].BindIDCard
		storeNameVal := rows.AddCell()
		storeNameVal.Value = items[i].StoreName
		usageTimeVal := rows.AddCell()
		usageTimeVal.Value = time.Unix(items[i].UsageTime, 0).Format("2006-01-02 15:04:05")
	}
	fi := fmt.Sprintf("核销明细-%d-[%s-%s].xlsx", tk.UID, time.Unix(start, 0).Format("20060102"), time.Unix(end, 0).Format("20060102"))
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
// 获取用户的所有激活年卡在本商户的消费情况
func Cards(ctx *gin.Context) {
	tks, _ := ctx.Get("token")
	tk := tks.(data.Token)
	usage_id_str := ctx.Param("usage_id")
	fmt.Println(usage_id_str)
	usage_id := strings.TrimPrefix(usage_id_str, "QR100")
	fmt.Println(usage_id)
	usageID, _ := strconv.ParseInt(usage_id, 10, 64)
	// 取用户的年卡列表
	mchID := tk.UID
	if tk.IsStore {
		has, info := module.MerchantStore.Get(tk.UID)
		if !has {
			ctx.JSON(http.StatusOK, errors.STORE_NOT_FOUND)
			return
		}
		mchID = info.MchId
	}
	items, _ := module.AnnualCard.Usages(usageID, mchID)
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"count": len(items),
		"data": items,
	})
}
// 核销年卡
func Usage(ctx *gin.Context) {
	tks, _ := ctx.Get("token")
	tk := tks.(data.Token)
	card_id := ctx.Param("card_id")
	cardID, _ := strconv.ParseInt(card_id, 10, 64)
	mchID := tk.UID
	storeID := tk.UID
	if tk.IsStore {
		has, info := module.MerchantStore.Get(tk.UID)
		if !has {
			ctx.JSON(http.StatusOK, errors.STORE_NOT_FOUND)
			return
		}
		mchID = info.MchId
	}
	cardInfo, has := module.AnnualCard.CardByID(cardID)
	if !has {
		ctx.JSON(http.StatusOK, errors.CARD_NOT_FOUND)
		return
	}
	if cardInfo.IsDelete == 1 {
		ctx.JSON(http.StatusOK, errors.CARD_NOT_FOUND)
		return
	}
	if cardInfo.IsActive != 1 {
		ctx.JSON(http.StatusOK, errors.CARD_NOT_ACTIVE)
		return
	}
	info := data.CardUsageLogInsert{
		CardID:cardInfo.CardID,
		MchID:mchID,
		StoreID:storeID,
	}
	affected, err := data.Db.Insert(&info)
	if affected < 1 || err != nil {
		ctx.JSON(http.StatusOK, errors.ERROR)
		return
	}
	ctx.JSON(http.StatusOK, errors.SUCCESS)
}