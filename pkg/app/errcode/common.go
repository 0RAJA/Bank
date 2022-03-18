package errcode

//以 1 开头表示公共错误码
var (
	Success                      = NewError(0, "成功")
	ServerErr                    = NewError(10000000, "服务内部错误")
	InvalidParamsErr             = NewError(10000001, "入参错误")
	NotFoundErr                  = NewError(10000002, "无结果")
	UnauthorizedAuthNotExistErr  = NewError(10000003, "鉴权失败, 无法解析")
	UnauthorizedTokenErr         = NewError(10000004, "鉴权失败，Token 错误")
	UnauthorizedTokenTimeoutErr  = NewError(10000005, "鉴权失败，Token 超时")
	UnauthorizedTokenGenerateErr = NewError(10000006, "鉴权失败，Token 生成失败")
	TooManyRequestsErr           = NewError(10000007, "请求过多")
	TimeOutErr                   = NewError(10000008, "请求超时")
	UnauthorizedNotLoginErr      = NewError(10000009, "鉴权失败,未登录")
	LoginErr                     = NewError(10000010, "登录失败")
	InsufficientPermissionsErr   = NewError(10000011, "权限不足")
)
