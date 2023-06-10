package codec

import "zura/pkg/errors"

const (
	UnSupportedRegisterTypeStatus = iota + 200001
	UnSupportedLoginTypeStatus
	UserRegisterdStatus
	PasswordNotConsistentStatus
	VerifyStatusWrongStatus
	UserNotFoundStatus
	PasswordWrongStatus
	OldPasswordWrongStatus

	PhoneEmptyStatus
	EmailEmptyStatus
	UsernameEmptyStatus
)

func init() {
	errors.NewCode(UnSupportedRegisterTypeStatus, "不支持该注册方式")
	errors.NewCode(UnSupportedLoginTypeStatus, "不支持该登录方式")
	errors.NewCode(UserRegisterdStatus, "用户已注册")
	errors.NewCode(PasswordNotConsistentStatus, "密码不一致")
	errors.NewCode(VerifyStatusWrongStatus, "验证码错误")
	errors.NewCode(UserNotFoundStatus, "用户不存在")
	errors.NewCode(PasswordWrongStatus, "密码错误")
	errors.NewCode(OldPasswordWrongStatus, "原密码错误")
	errors.NewCode(PhoneEmptyStatus, "手机号不能为空")
	errors.NewCode(EmailEmptyStatus, "邮箱不能为空")
	errors.NewCode(UsernameEmptyStatus, "用户名不能为空")
}
