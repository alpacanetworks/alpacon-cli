package note

type NoteListResponse struct {
	Count    int           `json:"count"`
	Current  int           `json:"current"`
	Next     int           `json:"next"`
	Previous string        `json:"previous"`
	Last     int           `json:"last"`
	Results  []NoteDetails `json:"results"`
}

type NoteDetails struct {
	ID      string `json:"id"`
	Server  string `json:"server"`
	Author  string `json:"author"`
	Content string `json:"content"`
	Private bool   `json:"private"`
	//	Pinned    bool   `json:"pinned"`
	//	UpdatedAt string `json:"updated_at"`
}
