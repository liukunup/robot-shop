package constant

const (
	AdminRole          = "admin"
	AdminUserID        = "1"
	MenuResourcePrefix = "menu:"
	ApiResourcePrefix  = "api:"
	PermSep            = ","
	OperatorRole       = "operator"
	GuestRole          = "guest"
	DateTimeLayout     = "2006-01-02 15:04:05"

	ResetPasswordSubject      = "重置密码"
	ResetPasswordTextTemplate = `尊敬的%s：

		我们收到了您的密码重置请求。请点击以下链接重置您的密码：
		%s

		如果您没有请求重置密码，请忽略此邮件。

		此链接将在24小时后失效。

		Robot Shop 团队
		`
)
