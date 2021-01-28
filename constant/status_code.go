package constant

const (
	CodeNotLoggedIn   = 101  // 未登录
	CodeLoginExpired  = 102  // 登录已过期
	CodeLoginInvalid  = 103  // 登录已失效
	CodeNoPermission  = 201  // 没有权限
	CodeBusinessError = 2998 // 常规业务错误
)

// 状态对应的文本内容
var codeText = map[int]string{
	CodeNotLoggedIn:  "未登录，请先登录",
	CodeLoginExpired: "登录已过期，请重新登录",
	CodeLoginInvalid: "登录已失效，请重新登录",
	CodeNoPermission: "没有权限，拒绝访问",
}

// 获取状态文本内容
func CodeText(code int) string {
	return codeText[code]
}
