package gateway

// Post
type LoginBody struct {
	Id       string `json:"id"`
	Password string `json:"password"`
	Offline  bool   `json:"offline"`
}

type LoginResult TokenResponse

type TokenResponse struct {
	AccessToken      string `json:"accessToken"`
	RefreshToken     string `json:"refreshToken"`
	ExpiresIn        int    `json:"expiresIn"`
	RefreshExpiresIn int    `json:"refreshExpiresIn"`
	TokenType        string `json:"tokenType"`
	SessionState     string `json:"sessionState"`
	Scope            string `json:"scope"`
}

type RefreshUserBody struct {
	Id           string `json:"id"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshUserResult TokenResponse
