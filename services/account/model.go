package account

import mysql_entity "saidakbar.origin/db/mysql/entity"

// region Verify Code
type (
	VerifyCodeModel struct {
		Phone string
		Code  string
	}
)

// region Update Account
type (
	UpdateAccountModel struct {
		ID       string
		FullName string
	}

	UpdateAccountResult struct {
		Account *mysql_entity.AccountModel
	}
)
