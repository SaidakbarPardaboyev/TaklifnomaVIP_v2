package account

import "saidakbar.origin/dto"

type Service interface {
	VerifyCode(model VerifyCodeModel) (result *dto.VerifyCodeResult, err error)
}
