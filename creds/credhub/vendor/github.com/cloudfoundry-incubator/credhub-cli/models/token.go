package models

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

func (t Token) ToYaml() string {
	return ""
}

func (t Token) ToJson() string {
	return ""
}
