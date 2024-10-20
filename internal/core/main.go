package core

type Author struct {
	Name         string `json:"name"`
	Handle       string `json:"handle"`
	Image        string `json:"image"`
	BlueVerified bool   `json:"blue_verified"`
}

type Post struct {
	Id     string `json:"id"`
	Author Author `json:"author"`
	Text   string `json:"text"`
	Likes  int    `json:"likes"`
}

type XtimatedPost struct {
	Post Post   `json:"post"`
	Type string `json:"type"`
}

type Thread struct {
	Start  Post   `json:"start"`
	Thread []Post `json:"thread"`
}

type XtimatedThread struct {
	Start  Post           `json:"start"`
	Thread []XtimatedPost `json:"thread"`
}
