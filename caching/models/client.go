package models

type ClientToken struct {
	TokenType          byte         `json:"token_type"`
	ActiveOrganization Organization `json:"active_organization"`
}

type ClientCredentials struct {
	ClientID     string       `json:"client_id"`
	ClientSecret string       `json:"client_secret"`
	Organization Organization `json:"organization"`
}
