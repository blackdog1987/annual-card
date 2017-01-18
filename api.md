# 年卡后台接口文档
> version 0.1
> author 廖才新
### 接口规范
接口使用restful API设计方式,安全方面使用`token`做校验,除了登陆接口之外，其余均需显式或者隐式传递`token`进行鉴权.

#### 接口定义
**接口地址** `http://120.76.223.14:8080`
****
### 基础接口
##### 1.管理员登录 /v1/manager/login
- method

```
post
```

- parameter

``` javascript
phone   // 必传 手机号
passwd  // 必传 密码 需传递密码的md5值
```

- response

``` javascript
{
  "code": 0,
  "msg": "SUCCESS",
  "data": {
    "expired_in": 1475610149,
    "name": "超级管理员",
    "token": "twxDIpZMECla16XKqaOeUJ.0000.1z5kPJosZywCjr7eh2eWS28X.00000.zuoNc6jqjYO3.00000.a.0000.Up0sp.0000.Yi78fKgx0A6IV4a1A.000000..000000."
  }
}
```

##### 2.上传图片 /v1/manager/upload/image
- method

```
post
```

- parameter

``` javascript
token     // 必传 令牌
code_type   // 必传 文件内容编码格式 如: base64
body    // 必传 文件内容
```

- response

``` javascript
{
  "code": 0,
  "data": "media/1475685707869297752.jpg",
  "info": "success"
}
```

##### 3.获取配置 /v1/manager/config/db
- method

```
get
```

- parameter

``` javascript
token     // 必传 令牌
```

- response

``` javascript
{
  "code": 0,
  "data": {
    "integral_rate": "1", // 积分兑换比例
    "spread_reg": "1",    // 推广注册送积分
    "spread_sale": "2"    // 推广销售送积分
  },
  "msg": "SUCCESS"
}
```

##### 4.修改配置 /v1/manager/config/db
- method

```
post
```

- parameter

``` javascript
token         // 必传 令牌
integral_rate // 积分兑换比例
spread_reg    // 推广注册送积分
spread_sale   // 推广销售送积分
```

- response

``` javascript
{
  "code": 0,
  "msg": "SUCCESS"
}
```

### 推广计划接口

##### 1.获取推广计划列表 /v1/manager/spread/plans
- method

```
get
```

- parameter

``` javascript
token     // 必传 令牌
p       // 页码  不传默认为1
state   // 是否关闭 -1 全部 0 未关闭 1 已关闭
```

- response

``` javascript
{
  "code": 0,
  "count": 1,
  "data": [
    {
    "sp_id": 1,                                // 计划ID
      "name": "测试推广计划",
      "channel": "微信",//渠道
      "contact": "023-11111111",// 联系电话
      "reg_commission": 100,// 推广注册
      "sale_commission": 10,// 推广销售
      "created": "2016-10-08T00:15:17+08:00",
      "updated": "2016-10-08T00:15:17+08:00",
      "is_disabled": 0,
      "qrcode": ""
    }
  ],
  "p": 0,
  "total": 1
}
```

##### 2.新增/编辑推广计划 /v1/manager/spread/plan
- method

```
post
```

- parameter

``` javascript
token         // 必传 令牌
id            // 计划id  编辑的时候传
name        // 必传 推广计划
channel       // 必传 渠道
contact       // 必传 联系电话
reg_commission    // 注册提成 默认为0 金额单位为分
sale_commission   // 销售提成 默认为0
```
- response

``` javascript
{
  "code": 0,
  "msg": "SUCCESS"
}
```

##### 3.获取单个推广计划详情 /v1/manager/spread/plan/:id
- method

```
get
```

- parameter

``` javascript
token         // 必传 令牌
id            // 必传 计划id
```

- response

``` javascript
{
  "code": 0,
  "data": {
    "sp_id": 1,
    "name": "测试推广计划",
    "channel": "微信",
    "contact": "023-11111111",
    "reg_commission": 10,
    "sale_commission": 1,
    "created": "2016-10-08T00:15:17+08:00",
    "updated": "2016-10-08T01:26:15+08:00",
    "is_disabled": 0,
    "qrcode": ""
  },
  "msg": "SUCCESS"
}
```

##### 4.获取推广明细 /v1/manager/spread/logs/:id
- method

```
get
```

- parameter

``` javascript
token         // 必传 令牌
id            // 必传 计划id
p         // 页码 默认1
```

- response

``` javascript
{
  "code": 0,
  "data": [
    {
          "log_id": 1,              // 记录ID
          "plan_id": 1,             // 计划ID
          "uid": 1,                 // 用户ID
          "consumer": {             // 消费者信息
            "nickname": "廖才新",
            "realname": "",
            "phone": ""
          },
          "category": 1,            // 1 注册 2 销售
          "commission": 10,         // 推广金额
          "order_total": 0,         // 订单金额  注册为0
          "created": 1476324026
        }
  ],
  "msg": "SUCCESS"
}
```

##### 5.修改推广计划状态 /v1/manager/spread/state/:id
- method

```
put
```

- parameter

``` javascript
token         // 必传 令牌
id            // 必传 计划id
state       // 状态 0 未关闭 1 关闭 默认为0
```

- response

``` javascript
{
  "code": 0,
  "msg": "SUCCESS"
}
```

### 商品管理

##### 1.获取年卡配置 /v1/manager/goods/annual
- method

```
get
```

- parameter

``` javascript
token         // 必传 令牌
```

- response

``` javascript
{
  "code": 0,
  "data": {
    "name": "年卡",
    "origin_price": 10000,
    "sale_price": 8000,
    "usage": "全城皆可使用",
    "images": "[{\"http://file.cbda.cn/uploadfile/2015/0330/20150330041852447.jpg\"}]"
  },
  "msg": "SUCCESS"
}
```

##### 2.保存年卡配置 /v1/manager/goods/annual
- method

```
post
```

- parameter

``` javascript
token         // 必传 令牌
name        // 必传 年卡名称
origin_price    // 必传 原价 单位为分
sale_price      // 必传 销售价 单位为分
usage         // 使用说明
images        // 图片列表 json字符串
```

- response

``` javascript
{
  "code": 0,
  "msg": "SUCCESS"
}
```

##### 3.获取优惠券配置 /v1/manager/goods/coupon
- method

```
get
```

- parameter

``` javascript
token         // 必传 令牌
```

- response

``` javascript
{
  "code": 0,
  "data": {
    "name": "优惠券",
    "offset_price": 10000,
    "sale_price": 8000,
    "usage": "全城皆可使用",
    "images": "[{\"http://file.cbda.cn/uploadfile/2015/0330/20150330041852447.jpg\"}]"
  },
  "msg": "SUCCESS"
}
```

##### 4.保存优惠券配置 /v1/manager/goods/coupon
- method

```
post
```
- parameter

``` javascript
token         // 必传 令牌
name        // 必传 年卡名称
offset_price    // 必传 原价 单位为分
sale_price      // 必传 销售价 单位为分
usage         // 使用说明
images        // 图片列表 json字符串
```
- response

``` javascript
{
  "code": 0,
  "msg": "SUCCESS"
}
```

### 商户信息

##### 1.获取商户列表 /v1/manager/merchants
- method

```
get
```
- parameter

``` javascript
token         // 必传 令牌
keyword       // 关键字
p             // 页码
```

- response

```java
{
  "code": 0,
  "count": 1,
  "data": [
    {
      "mch_id": 1,              // 商户ID
      "mch_name": "乡村基",      
      "value": "100",
      "consume": "全天可用",     // 消费特惠
      "usage": "1元买巨无霸",    // 使用说明
      "contact": "023-12345678",// 联系电话
      "address": "[\"重庆市南岸区\"]", // 地址
      "introduce": "著名餐饮品牌", // 详细介绍
      "cover": "http://desk.fd.zol-img.com.cn/t_s960x600c5/g4/M01/0D/04/Cg-4WVP_npmIY6GRAKcKYPPMR3wAAQ8LgNIuTMApwp4015.jpg",
      "imgs": "[\"http://desk.fd.zol-img.com.cn/t_s960x600c5/g4/M01/0D/04/Cg-4WVP_npmIY6GRAKcKYPPMR3wAAQ8LgNIuTMApwp4015.jpg\",\"http://desk.fd.zol-img.com.cn/t_s960x600c5/g4/M01/0D/04/Cg-4WVP_npmIY6GRAKcKYPPMR3wAAQ8LgNIuTMApwp4015.jpg\"]",
      "created": 1475920908,
      "updated": 1475920908,
      "state": 0
    }
  ],
  "p": 0,
  "total": 1
}
```
##### 2.获取商户详情 /v1/manager/merchant/:id
- method

```
get
```

- parameter

``` javascript
token         // 必传 令牌
id        // 必传 商户ID
```

- response

``` javascript
{
  "code": 0,
  "msg": "SUCCESS",
  "data":{
      "mch_id": 1,              // 商户ID
      "mch_name": "乡村基",      
      "value": "100",
      "consume": "全天可用",     // 消费特惠
      "usage": "1元买巨无霸",    // 使用说明
      "contact": "023-12345678",// 联系电话
      "address": "[\"重庆市南岸区\"]", // 地址
      "introduce": "著名餐饮品牌", // 详细介绍
      "cover": "http://desk.fd.zol-img.com.cn/t_s960x600c5/g4/M01/0D/04/Cg-4WVP_npmIY6GRAKcKYPPMR3wAAQ8LgNIuTMApwp4015.jpg",
      "imgs": "[\"http://desk.fd.zol-img.com.cn/t_s960x600c5/g4/M01/0D/04/Cg-4WVP_npmIY6GRAKcKYPPMR3wAAQ8LgNIuTMApwp4015.jpg\",\"http://desk.fd.zol-img.com.cn/t_s960x600c5/g4/M01/0D/04/Cg-4WVP_npmIY6GRAKcKYPPMR3wAAQ8LgNIuTMApwp4015.jpg\"]",
      "created": 1475920908,
      "updated": 1475920908,
      "state": 0
    }
}
```
##### 3.新增／编辑商户 /v1/manager/merchant
- method

```
post
```

- parameter

``` javascript
token       // 必传 令牌
name          // 必传 商户名称
value         // 价值
consume       // 消费特惠
usage         // 使用说明
contact       // 联系电话
address       // 地址 json字符串
introduce     // 详细描述
cover         // 封面图
images        // 图片数组 json字符串
```

- response

``` javascript
{
  "code": 0,
  "msg": “SUCCESS”
}
```
### 销售相关
##### 1.获取订单列表 /v1/manager/orders
- method

```
get
```

- parameter

``` javascript
token       // 必传 令牌
p           // 页码
start       // 起始时间 unix时间戳
end         // 结束时间 unix时间戳
```

- response

``` javascript
{
  "code": 0,
  "count": 0,
  "data": [
    {
      "order_id": 1,                        // 订单ID
      "order_no": "1476338291178123734",    // 订单号
      "uid": 1,                             // 用户ID
      "goods_name": "优惠券",                // 商品名称
      "transaction_id": "",                 // 外部交易流水号
      "price": 0,                           // 商品价格
      "total": 0,                           // 实际支付金额
      "category": "coupon",                 // 商品品种 coupon 优惠券 annual 年卡
      "points": 2,                          // 使用积分数量
      "points_price": 2,                    // 积分抵扣金额
      "is_pay": 1,                          // 是否支付 1 支付成功
      "is_coupon": 0,                       // 是否使用优惠券
      "coupon_price": 0,                    // 优惠券抵扣金额
      "created": 1476338291,                // 订单产生时间
      "updated": 1476338291                 // 订单更新时间
    }
  ],
  "p": 0,
  "total": 0
}
```
##### 2.获取年卡会员列表 /v1/manager/members
- method

```
get
```

- parameter

``` javascript
token       // 必传 令牌
p           // 页码
start       // 起始时间 unix时间戳
end         // 结束时间 unix时间戳
```

- response

``` javascript
{
  "code": 0,
  "count": 0,
  "data": [
    {
      "card_id": 1,                       // 年卡ID
      "plan_id": 1,                       // 计划ID
      "card_name": "年卡",                 // 卡名称
      "card_no": "11475913761377258330",  // 卡号
      "relation_uid": 0,                  // 绑定用户
      "bind_headimg": "",                 // 绑定半身照片
      "bind_name": "",                    // 绑定姓名
      "bind_contact": "",                 // 绑定联系电话
      "bind_idcard": "",                  // 绑定身份证号
      "card_passwd": "805634",            // 卡密
      "expired_start": 1475884800,        // 有效期 起始
      "expired_stop": 1476057600,         // 有效期 结束
      "is_active": 0,                     // 是否激活 1:激活 0: 待激活
      "created": 1475913761,              // 创建时间
      "updated": 1475913761               // 更新时间
    }
  ],
  "p": 0,
  "total": 0
}
```
##### 3.新建年卡推广计划 /v1/manager/card/plan
- method

```
post
```

- parameter

``` javascript
token               // 必传 令牌
channel             // 渠道
expired_start       // 起始时间 unix时间戳
expired_stop        // 结束时间 unix时间戳
create_num          // 年卡生成数量
```

- response

``` javascript
{
  "code": 0,
  "msg": "SUCCESS"
}
```
##### 4.年卡推广计划列表 /v1/manager/card/plans
- method

```
get
```

- parameter

``` javascript
token       // 必传 令牌
p           // 页码
```

- response

``` javascript
{
  "code": 0,
  "count": 5,
  "data": [
    {
      "cp_id": 1,                   // 计划ID
      "channel": "测试渠道",         // 渠道 
      "expired_start": 1475884800,  // 有效期 起始
      "expired_stop": 1476057600,   // 有效期 结束
      "create_num": 100,            // 生产年卡数量
      "active_num": 0,              // 激活年卡数量
      "created": 1475913761,
      "updated": 1475913761,
      "is_disabled": 0
    },
    {
      "cp_id": 2,
      "channel": "测试渠道",
      "expired_start": 1475884800,
      "expired_stop": 1476057600,
      "create_num": 100,
      "active_num": 0,
      "created": 1475914076,
      "updated": 1475914076,
      "is_disabled": 0
    },
    {
      "cp_id": 3,
      "channel": "测试渠道",
      "expired_start": 1475884800,
      "expired_stop": 1476057600,
      "create_num": 100,
      "active_num": 0,
      "created": 1475916631,
      "updated": 1475916631,
      "is_disabled": 0
    },
    {
      "cp_id": 4,
      "channel": "测试渠道",
      "expired_start": 1475884800,
      "expired_stop": 1476057600,
      "create_num": 2,
      "active_num": 0,
      "created": 1475916779,
      "updated": 1475916779,
      "is_disabled": 0
    },
    {
      "cp_id": 5,
      "channel": "测试渠道",
      "expired_start": 1475884800,
      "expired_stop": 1476057600,
      "create_num": 2,
      "active_num": 0,
      "created": 1475916854,
      "updated": 1475916854,
      "is_disabled": 0
    }
  ],
  "p": 0,
  "total": 5
}
```
##### 5.获取年卡计划明细 /v1/manager/card/detail/:id
- method

```
get
```

- parameter

``` javascript
token       // 必传 令牌
p           // 页码
id          // 年卡计划ID
start       // 起始时间 unix时间戳
end         // 结束时间 unix时间戳
```

- response

``` javascript
{
  "code": 0,
  "count": 1,
  "data": [
    {
      "card_id": 1,                       // 年卡ID
      "plan_id": 1,                       // 计划ID
      "card_name": "年卡",                 // 卡名称
      "card_no": "11475913761377258330",  // 卡号
      "relation_uid": 0,                  // 绑定用户
      "bind_headimg": "",                 // 绑定半身照片
      "bind_name": "",                    // 绑定姓名
      "bind_contact": "",                 // 绑定联系电话
      "bind_idcard": "",                  // 绑定身份证号
      "card_passwd": "805634",            // 卡密
      "expired_start": 1475884800,        // 有效期 起始
      "expired_stop": 1476057600,         // 有效期 结束
      "is_active": 0,                     // 是否激活 1:激活 0: 待激活
      "created": 1475913761,              // 创建时间
      "updated": 1475913761               // 更新时间
    }
  ],
  "p": 0,
  "total": 1
}
```
### Banner
##### 1.获取商品主题图列表 /v1/manager/banners
- method

```
get
```

- parameter

``` javascript
token       // 必传 令牌
```

- response

``` javascript
{
  "code": 0,
  "data": [
    {
      "banner_id": 2,
      "name": "至尊pissa",
      "img": "http://b.hiphotos.baidu.com/image/h%3D200/sign=066d5b2977cf3bc7f700caece100babd/f636afc379310a55d27d79d0b04543a9822610bc.jpg",
      "link": "",
      "ordid": 1,
      "created": 1476033381
    }
  ],
  "msg": "SUCCESS"
}
```
##### 2.新增／编辑 banner /v1/manager/banner
- method

```
post
```

- parameter

``` javascript
token       // 必传 令牌
id          // banner ID 编辑的时候传
name        // 名称
img         // 图片地址
link        // 跳转地址
ordid       // 排序 正序
```

- response

``` javascript
{
  "code": 0,
  "msg": "SUCCESS"
}
```
##### 3.获取单个 banner /v1/manager/banner/:id
- method

```
get
```

- parameter

``` javascript
token       // 必传 令牌
id          // banner ID 
```

- response

``` javascript
{
  "code": 0,
  "data": {
    "banner_id": 2,
    "name": "至尊pissa",
    "img": "http://b.hiphotos.baidu.com/image/h%3D200/sign=066d5b2977cf3bc7f700caece100babd/f636afc379310a55d27d79d0b04543a9822610bc.jpg",
    "link": "",
    "ordid": 1,
    "created": 1476033381
  },
  "msg": "SUCCESS"
}
```
##### 4.删除 banner /v1/manager/banner/:id
- method

```
delete
```

- parameter

``` javascript
token       // 必传 令牌
id          // banner ID 
```

- response

``` javascript
{
  "code": 0,
  "msg": "SUCCESS"
}
```

### 商户账户

##### 1.获取商户账户列表 /v1/manager/merchant-accounts
- method

```
get
```
- parameter

``` javascript
token         // 必传 令牌
keyword       // 关键字
p             // 页码


```
- response

``` javascript
{
  "code": 0,
  "count": 1,
  "data": [
    {
      "mch_id": 43,        // 商户ID
      "name": "乡村基",     // 商户名称
      "account": "乡村基1", // 商户账号
      "created": 1481962822,
      "updated": 1481962927,
      "state": 0            // 是否禁用 -1
    }
  ],
  "p": 0,
  "total": 1
}
```

##### 2.获取商户详情 /v1/manager/merchant-account/:id
- method

```
get
```

- parameter

``` javascript
token         // 必传 令牌
id        // 必传 商户ID
```

- response

``` javascript
{
  "code": 0,
  "data": {
    "mch_id": 43,
    "name": "乡村基",
    "account": "乡村基",
    "created": 1481962822,
    "updated": 1481962822,
    "state": 0
  },
  "msg": "SUCCESS"
}
```

##### 3.新增／编辑商户 /v1/manager/merchant-account
- method

```
post
```

- parameter

``` javascript
id            // 商户id 不传为新增
token       // 必传 令牌
name          // 必传 商户名称
account       // 必传 商户账号
passwd        // 必传 密码
```

- response

``` javascript
{
  "code": 0,
  "msg": “SUCCESS”
}
```

##### 4.禁用／启用商户 /v1/manager/merchant-account-state/:id
- method

```
post
```

- parameter

``` javascript
token       // 必传 令牌
id            // 必传 商户名称
state         // 状态 -1 禁用 0 正常
```

- response

``` javascript
{
  "code": 0,
  "msg": “SUCCESS”
}
```

##### 5.菜单创建
```
{
     "button":[
     {
          "name":"了解年卡",
          "sub_button":[
          {
              "type":"view",
              "name":"年卡介绍",
              "url":"http://mp.weixin.qq.com/s?__biz=MzIxNjYyMTQ0Ng==&mid=2247483819&idx=2&sn=b495e4fe61c055c93e7cc467fbfbdf44&chksm=97870942a0f080549ca0aabfd5fe238243ac95ea009c7aa2612bcf3947b216ffdfd5587736d3&mpshare=1&scene=1&srcid=010110GrOEobKA2bvs1XlXcl#rd"
           },
           {
              "type":"view",
              "name":"年卡攻略",
              "url":"http://mp.weixin.qq.com/s?__biz=MzIxNjYyMTQ0Ng==&mid=2247483861&idx=2&sn=de18ac898053e8f9fb3907cecf74680c&chksm=9787093ca0f0802a983d929c5c82fe5370796ff7b543040351edf602d3950c1bfd2feb2bb58e&mpshare=1&scene=1&srcid=0101Iy200ypcoeRN19hj5FlH#rd"
           },
           {
              "type":"view",
              "name":"视频集锦",
              "url":"http://v.qq.com/x/page/i03433nhmjq.html"
           },
           {
              "type":"view",
              "name":"常见问答",
              "url":"http://mp.weixin.qq.com/s?__biz=MzIxNjYyMTQ0Ng==&mid=2247483734&idx=2&sn=d3aa72866a89f569144f73e62a910bc0&chksm=978709bfa0f080a90fc9783b1f24fba314b2ac322c3d7e84eebd58a6bfc8c67da0b6b100b85e&mpshare=1&scene=1&srcid=1119pVw88HZMmUESjyMB3qFX#rd"
           }]
      },
      {
           "name":"年卡购买",
           "sub_button":[
           {
              "type":"view",
              "name":"积分换卡",
              "url":"http://mp.weixin.qq.com/s?__biz=MzIxNjYyMTQ0Ng==&mid=2247483861&idx=1&sn=c8db63bea8f8e418932797775e2b9900&chksm=9787093ca0f0802a7a20547d5e913d6c6d0670d75e2c324fb451685df78ff2a862bae22ca800&mpshare=1&scene=1&srcid=0101ta0M4qQ7bqRMeAMsBDSd#rd"
           },
           {
               "type":"view",
               "name":"年卡购买",
               "url":"https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx2df87c500b984c2d&redirect_uri=http%3a%2f%2fwx.qinzinianka.com%2fwechat%2fauthorize%3ffrom%3dgoods/annual&response_type=code&scope=snsapi_userinfo&state=annual#wechat_redirect"
            }]
       },
       {
           "name":"会员中心",
           "sub_button":[
           {
               "type":"view",
               "name":"我的年卡",
               "url":"https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx2df87c500b984c2d&redirect_uri=http%3a%2f%2fwx.qinzinianka.com%2fwechat%2fauthorize%3ffrom%3dmycard&response_type=code&scope=snsapi_userinfo&state=annual#wechat_redirect"
            },
            {
               "type":"view",
               "name":"获得积分",
               "url":"https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx2df87c500b984c2d&redirect_uri=http%3a%2f%2fwx.qinzinianka.com%2fwechat%2fauthorize%3ffrom%3dmypoints&response_type=code&scope=snsapi_userinfo&state=annual#wechat_redirect"
            },
            {
               "type":"view",
               "name":"年卡激活",
               "url":"https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx2df87c500b984c2d&redirect_uri=http%3a%2f%2fwx.qinzinianka.com%2fwechat%2fauthorize%3ffrom%3dcard/list&response_type=code&scope=snsapi_userinfo&state=annual#wechat_redirect"
            },
            {
               "type":"view",
               "name":"商户登录",
               "url":"http://wx.qinzinianka.com/wxmerchant/#/checkcard"
            }]
       }]
 }
```