package v1

type Doc struct {
	ID       string    `json:id`
	Title    string    `json:"title"`
	Desc     string    `json:"desc"`
	Url      string    `json:"url"`
	Img      []Doc_img `json:"img"`
	Source   int       `json:"source"`
	Keys     string    `json:"keys"`
	Upload   string    `json:"upload"`
	Chanid   int       `json:"chanid"`
	Isparsed bool      `json:"isparsed"`
}
type Doc_img struct {
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
