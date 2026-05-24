package contracts

import "saidakbar.origin/caching/models"

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
		ID:    acc.ActiveOrganization.ID,
		Name:  acc.Name,
		Phone: acc.Username,
	}
}
