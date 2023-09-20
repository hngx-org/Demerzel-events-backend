package types

type UserInfo struct {
	Sub        string `json:"sub"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email      string `json:"email"`
	Picture    string `json:"picture"`
	Locale     string `json:"locale"`
}

type UserData struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

type UserUpdatables struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type JwtCustomClaims struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}
