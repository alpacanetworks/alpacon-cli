package auth0

type DeviceCodeResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
	VerificationURIComplete string `json:"verification_uri_complete"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Error        string `json:"error,omitempty"`
	ErrorDesc    string `json:"error_description,omitempty"`
}

type Auth0Config struct {
	Method   string `json:"method"`
	ClientID string `json:"client_id,omitempty"`
	Domain   string `json:"domain,omitempty"`
	Audience string `json:"audience,omitempty"`
}

type AuthEnvResponse struct {
	Auth0    Auth0Config `json:"auth0"`
	Language string      `json:"language"`
}
