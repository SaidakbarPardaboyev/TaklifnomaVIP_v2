package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Account struct {
	ID                 primitive.ObjectID `json:"id"`
	Username           string             `json:"username"`
	Organizations      []*Organization    `json:"organizations"`
	Name               string             `json:"name"`
	Type               byte               `json:"type"`
	ActiveOrganization *Organization      `json:"active_organization"`
	TokenType          byte               `json:"token_type"`
}

type Organization struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Inn   *string `json:"inn"`
	Pinfl *string `json:"pinfl"`
}

func (o Organization) GetInn() string {
	if o.Inn != nil {
		return *o.Inn
	}

	return ""
}

func (o Organization) GetPinfl() string {
	if o.Pinfl != nil {
		return *o.Pinfl
	}

	return ""
}

func (o Organization) IsSameInnOrPinfl(innOrPinfl string) bool {
	return o.GetInn() == innOrPinfl || o.GetPinfl() == innOrPinfl
}
