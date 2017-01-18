package errors

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

var (
	SUCCESS = &Error{0, "SUCCESS"}
	ERROR = &Error{-1, "服务开小差了"}
	PARAM_PARSE_ERR = &Error{-2000, "参数解析错误"}
	PHONE_VALID_ERR = &Error{-2001, "手机号格式错误"}
	PASSWD_VALID_ERR = &Error{-2002, "密码错误"}
	TOKEN_VALID_ERR = &Error{-2003, "令牌校验错误"}
	AUTHORIFY_VALID_ERR = &Error{-2004, "您没有该操作权限"}
	TOKEN_EXPIRED = &Error{-2005, "令牌已过期"}
	PASSWD_EMPTY = &Error{-2006, "密码不能为空"}
	PHONE_REPEAT = &Error{-2007, "手机号已经存在"}
	NAME_EMPTY = &Error{-2008, "名称不能为空"}
	SALE_PRICE_ZERO = &Error{-2009, "销售价格不能小于1"}
	CHANNEL_EMPTY = &Error{-2010, "渠道不能为空"}
	EXPIRED_START_ZERO = &Error{-2011, "生效时间为空"}
	EXPIRED_STOP_ZERO = &Error{-2012, "失效时间为空"}
	CREATE_NUM_ZERO = &Error{-2014, "生成数量不能为0"}
	ID_EMPTY = &Error{-2015, "没有需要操作的数据"}
	CAPTCHA_SEND_ERROR = &Error{-2016, "验证码发送失败"}
	CAPTCHA_VALID_ERR = &Error{-2017, "验证码校验失败"}
	UPLOAD_FAILED = &Error{-2018, "上传失败"}
	SPREAD_PLAN_NAME_EMPTY = &Error{-2019, "推广计划不能为空"}
	SPREAD_PLAN_NAME_REPEAT = &Error{-2020, "推广计划已经存在"}
	SPREAD_CHANNEL_EMPTY = &Error{-2021, "推广渠道不能为空"}
	CONTACT_EMPTY = &Error{-2022, "联系电话不能为空"}
	CREATE_NUM_MAX = &Error{-2023, "单次最多生成500"}
	CARD_PLAN_ID_ZERO = &Error{-2024, "计划ID不能为0"}
	NAME_REPEAT = &Error{-2025, "名称已经存在"}
	ACCOUNT_REPEAT = &Error{-2027, "账号已经存在"}
	IMAGE_EMPTY = &Error{-2026, "图片不能为空"}
	ACCOUNT_EMPTY = &Error{-2028, "账号不能为空"}
	ACCOUNT_DISABLED = &Error{-2029, "账号已被禁用"}
	STORE_NOT_FOUND = &Error{-2030, "门店不存在"}
	CARD_NOT_FOUND = &Error{-2031, "年卡不存在"}
	CARD_NOT_ACTIVE = &Error{-2032, "年卡还没有激活"}
	CARD_USAGED = &Error{-2033, "年卡已经使用过"}
)
