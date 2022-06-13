package errcode

//Common  code
var (
	Success = NewError(200, "成功")
)

var (
	ServerError               = NewError(10000, "服务内部错误")
	InvalidParams             = NewError(10001, "入参错误")
	NoToken                   = NewError(100011, "未传递Token")
	NotFound                  = NewError(10002, "找不到")
	UnauthorizedAuthNotExist  = NewError(10003, "鉴权失败，找不到对应的AppKey和AppSecret")
	UnauthorizedTokenError    = NewError(10004, "鉴权失败，Token错误")
	UnauthorizedTokenTimeout  = NewError(10005, "鉴权失败，Token超时")
	UnauthorizedTokenGenerate = NewError(10006, "鉴权失败，Token生成失败")
	TooManyRequests           = NewError(10007, "请求过多")
)

var (
	LoginFailed         = NewError(20001, "登录失败")
	UserNotExistError   = NewError(20002, "用户不存在")
	UserExistsError     = NewError(20003, "用户已存在")
	MobileCodeError     = NewError(20004, "验证码错误")
	GetSMSCodeError     = NewError(20005, "验证码发送失败")
	GetSMSCodeTooFast   = NewError(20006, "请求验证码过于频繁")
	ModifyUserInfoError = NewError(20007, "修改用户信息失败")

	NotAllowFileExts = NewError(20008, "不支持上传的文件类型")
	NotAllowFileSize = NewError(20009, "文件大小超过限制")
)

var (
	VideoPublishFailed = NewError(30001, "视频上传出错")
	DuplicateLike      = NewError(30002, "重复点赞")
	UnLike             = NewError(30003, "未点赞")
)

var (
	DuplicateFollow = NewError(40001, "重复关注")
	UnFollow        = NewError(40002, "未关注")
)
