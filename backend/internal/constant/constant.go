package constant

const (
	// 展示名称
	DisplayName = "Robot Shop"

	// 超级管理员
	AdminRole   = "admin"
	AdminUserID = "1"

	// 角色标识
	OperatorRole = "operator"
	UserRole     = "user"
	GuestRole    = "guest"

	// 资源前缀、分隔符、权限
	MenuResourcePrefix = "menu:"
	ApiResourcePrefix  = "api:"
	PermSep            = ","
	PermRead           = "read"

	// 默认的日期时间展示格式
	DateTimeLayout = "2006-01-02 15:04:05"

	// 重置密码
	ResetPasswordSubject      = "重置密码"
	ResetPasswordTextTemplate = `尊敬的%s：

		我们收到了您的密码重置请求。请点击以下链接重置您的密码：
		%s

		如果您没有请求重置密码，请忽略此邮件。

		此链接将在24小时后失效。

		%s 团队
		`

	// 账号状态
	UserStatusPending = 0 // 待激活
	UserStatusActive  = 1 // 正常
	UserStatusBanned  = 2 // 禁用
)
