package contracts

import (
	mysql_entity "saidakbar.origin/db/mysql/entity"
	"saidakbar.origin/caching/models"
)

type AccountContract struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func CreateAccountContract(acc *models.Account) AccountContract {
	if acc == nil {
		return AccountContract{}
	}

	return AccountContract{
		ID:    acc.ID,
		Name:  acc.FullName,
		Phone: acc.Phone,
	}
}

func CreateAccountContractFromEntity(acc *mysql_entity.AccountModel) AccountContract {
	if acc == nil {
		return AccountContract{}
	}

	return AccountContract{
		ID:    acc.ID,
		Name:  acc.FullName,
		Phone: acc.Phone,
	}
}
