package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/molibei/annual-card/handle"
	"github.com/molibei/annual-card/handle/backend"
	"github.com/molibei/annual-card/handle/front"
	"github.com/molibei/annual-card/lib/data"
	_ "github.com/molibei/annual-card/module/annual"
	_ "github.com/molibei/annual-card/module/conf"
	_ "github.com/molibei/annual-card/module/manager"
	_ "github.com/molibei/annual-card/module/merchant"
	_ "github.com/molibei/annual-card/module/order"
	_ "github.com/molibei/annual-card/module/spread"
	_ "github.com/molibei/annual-card/module/token"
	_ "github.com/molibei/annual-card/module/user"
	_ "github.com/molibei/annual-card/module/wechat"
	_ "github.com/molibei/annual-card/module/wxpay"
	"strings"
	"image/jpeg"
	"github.com/nfnt/resize"
	"net/http"
	"strconv"
)

var (
	err      error
	confPath = flag.String("c", "./conf.json", "Usage:include conf file path")
)

func init() {
	flag.Parse()
	err = data.LoadBaseConf(*confPath)
	if err != nil {
		fmt.Println("load conf file failed. err of ", err.Error())
		os.Exit(1)
	}
	err = data.ConnectMySQL()
	if err != nil {
		fmt.Println("connect sql failed. err of ", err.Error())
		os.Exit(1)
	}
}

func main() {
	gin.SetMode(data.BaseConf.Server.Mode)
	router := gin.New()
	store := sessions.NewCookieStore([]byte(data.BaseConf.Server.Secret))
	router.Use(sessions.Sessions("me", store))
	router.Use(handle.Identity())
	router.GET("/media/:filename", func(ctx *gin.Context) {
		fi := ctx.Param("filename")
		ctx.File("./media/" + fi)
	})
	router.GET("/head/:filename", func(ctx *gin.Context) {
		fi := ctx.Param("filename")
		// 处理缩略图
		fis := strings.Split(fi, "@")
		fp := "./head/" + fis[0]
		if len(fis) == 1 {
			// 原图
			ctx.File(fp)
			return
		}
		size, _:= strconv.Atoi(fis[1])
		if size > 500 {
			// 原图
			ctx.File(fp)
			return
		}
		// 压缩图片
		file, err := os.Open(fp)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg": "获取图像资源失败",
			})
			return
		}

		// decode jpeg into image.Image
		img, err := jpeg.Decode(file)
		bounds := img.Bounds()
		dx:= bounds.Dx()
		dy:= bounds.Dy()
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg": "获取图像资源失败",
			})
			return
		}
		defer file.Close()
		m := resize.Thumbnail(uint(size), uint(size*dy/dx), img, resize.MitchellNetravali)
		// write new image to file
		jpeg.Encode(ctx.Writer, m, nil)
	})
	wechat := router.Group("/wechat")
	{
		// 接收事件
		wechat.POST("/receive", front.Receive)
		wechat.GET("/receive", front.WxValid)
		// 授权
		wechat.GET("/authorize", front.Authorize)
	}
	api := router.Group("/v1")
	{
		admin := api.Group("/manage")
		admin.POST("/login", backend.Login)
		admin.Use(handle.Authorization())
		{
			admin.GET("/download/xlx/:fi", backend.DownXls)
			// 上传图片
			admin.POST("/upload/image", backend.UploadImage)
			/*    推广    */
			// 获取推广计划列表
			admin.GET("/spread/plans", backend.SpreadPlans)
			// 新增|编辑 推广计划
			admin.POST("/spread/plan", backend.SaveSpreadPlan)
			// 获取单个推广计划
			admin.GET("/spread/plan/:id", backend.SpreadPlan)
			// 获取推广明细
			admin.GET("/spread/logs/:id", backend.SpreadPlanLogs)
			// 推广明细导出xls
			admin.GET("/spread/logs-export/:id", backend.SpreadPlanLogs2xls)
			// 获取二维码
			admin.GET("/spread/qrcode/:id", backend.SpreadPlanQrcode)
			// 关闭推广计划
			admin.PUT("/spread/state/:id", backend.SpreadPlanState)
			/*    商品    */
			// 获取年卡配置
			admin.GET("/goods/annual", backend.AnnualCard)
			// 修改年卡
			admin.POST("/goods/annual", backend.SaveAnnual)
			// 获取优惠券
			admin.GET("/goods/coupon", backend.Coupon)
			// 修改优惠券
			admin.POST("/goods/coupon", backend.SaveCoupon)
			/*    订单    */
			// 获取激活须知
			admin.GET("/setting/active-help", backend.ActiveHelp)
			// 保存激活须知
			admin.POST("/setting/active-help",backend.SaveActiveHelp)
			// 订单列表
			admin.GET("/orders", backend.Orders)
			// 导出订单列表
			admin.GET("/orders/export", backend.OrdersExport)
			// 年卡会员列表
			admin.GET("/members", backend.Members)
			// 删除年卡
			admin.DELETE("/member/:card_id", backend.DeleteCard)
			// 取消激活
			admin.PUT("/member/unactive/:card_id", backend.UnActive)
			// 导出年卡会员
			admin.GET("/members/export", backend.MembersExport)
			// 新建年卡推广计划
			admin.POST("/card/plan", backend.AddAnnualCardPlan)
			// 年卡生成计划列表
			admin.GET("/card/plans", backend.AnnualCardPlans)
			// 年卡计划明细
			admin.GET("/card/detail/:id", backend.AnnualCardDetail)
			// 导出年卡明细
			admin.GET("/card/detail-export/:id", backend.AnnualCardDetailExport)
			/*    商户       */
			// 获取单个商户详情
			admin.GET("/merchant/:id", backend.Merchant)
			// 获取商户列表
			admin.GET("/merchants", backend.Merchants)
			// 新增|编辑商户
			admin.POST("/merchant", backend.SaveMerchant)
			/*    配置       */
			// 修改配置
			admin.POST("/config/db", backend.SaveDbConfig)
			// 获取配置
			admin.GET("/config/db", backend.DbConfig)
			// 商户账号管理
			// 获取商户账号列表
			admin.GET("/merchant-accounts", backend.MerchantAccounts)
			// 获取单个商户详情
			admin.GET("/merchant-account/:id", backend.MerchantAccount)
			// 新增 | 编辑商户账户
			admin.POST("/merchant-account", backend.SaveMerchantAccount)
			// 锁定|解锁
			admin.PUT("/merchant-account-state/:id", backend.ToggleMerchantAccountState)
			// 修改密码
			admin.POST("/setting/reset-pwd", backend.ResetPwd)
		}
		merchant_account := api.Group("/merchant")
		merchant_account.POST("/login",backend.MerchantLogin)
		merchant_account.Use(handle.Authorization())
		{
			// 修改密码
			merchant_account.POST("/reset-pwd", backend.ResetMerchantAccountPwd)
			// 门店列表
			merchant_account.GET("/stores", backend.Stores)
			// 单个门店账号
			merchant_account.GET("/store/:id", backend.Store)
			// 新建| 编辑门店
			merchant_account.POST("/store", backend.SaveMerchantStore)
			// 禁用|启用门店
			merchant_account.PUT("/store-state/:id", backend.ToggleMerchantStoreState)
			// 年卡用户使用列表
			merchant_account.GET("/card-visitors", backend.CardUsages)
			merchant_account.GET("/card-visitors-export", backend.CardUsagesExport)
			// 获取用户年卡详情
			merchant_account.GET("/cards/:usage_id", backend.Cards)
			// 核销年卡
			merchant_account.POST("/usage/:card_id", backend.Usage)
		}
		// 获取产品列表
		api.GET("/goods", front.Goods)
		// 获取商户列表
		api.GET("/merchants", front.Merchants)
		api.GET("/merchant", front.Merchant)
		api.GET("/setting/active-help", backend.ActiveHelp)
		api.GET("/getWxJsConfig", front.GetWxJsSign)
		api.POST("/wxpay_notify", front.WxPayNotify)
		auth := api.Group("/")
		auth.Use(handle.Auth())
		{
			auth.GET("/preorder", front.PreOrder)
			auth.POST("/getPaySign", front.GetPaySign)
			auth.GET("/getUserInfo", front.UserInfo)
			auth.GET("/spreads", front.Spreads)
			auth.GET("/picture", front.Picture)
			auth.GET("/cardNum", front.CardNum)
			auth.POST("/card/add", front.AddCard)
			auth.GET("/cards", front.Cards)
			auth.GET("/card", front.Card)
			auth.PUT("/card", front.BindCard)
			auth.GET("/qrcode", front.QrCode)
			auth.PUT("/order/reset/:orderNo", front.ResetOrder)
		}
	}
	router.Run(data.BaseConf.Server.Port)
}
