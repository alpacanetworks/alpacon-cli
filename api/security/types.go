package security

type CommandAclRequest struct {
	Token   string `json:"token"`
	Command string `json:"command"`
}

type CommandAclResponse struct {
	Id        string `json:"id"`
	Token     string `json:"token"`
	TokenName string `json:"token_name"`
	Command   string `json:"command"`
}
