package websh

type SessionRequest struct {
	Rows   int    `json:"rows"`
	Cols   int    `json:"cols"`
	Server string `json:"server"` // server id
	Root   bool   `json:"root"`
}

type SessionResponse struct {
	ID           string `json:"id"`
	Rows         int    `json:"rows"`
	Cols         int    `json:"cols"`
	Server       string `json:"server"`
	User         string `json:"user"`
	Root         bool   `json:"root"`
	UserAgent    string `json:"user_agent"`
	RemoteIP     string `json:"remote_ip"`
	WebsocketURL string `json:"websocket_url"`
}
