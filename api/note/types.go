package note

type NoteDetails struct {
	ID      string `json:"id"`
	Server  string `json:"server"`
	Author  string `json:"author"`
	Content string `json:"content"`
	Private bool   `json:"private"`
	//	Pinned    bool   `json:"pinned"`
	//	UpdatedAt string `json:"updated_at"`
}

type NoteCreateRequest struct {
	Server  string `json:"server"`
	Content string `json:"content"`
	Private bool   `json:"private"`
	Pinned  bool   `json:"pinned"`
}
