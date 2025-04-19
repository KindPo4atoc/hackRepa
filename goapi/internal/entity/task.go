package entity

type Task struct {
	IdTask          int    `json:"id"`
	Header          string `json:"header"`
	DescriptionTask string `json:"description"`
	LevelTask       string `json:"level"`
	LoginAuthor     string `json:"author"`
	PathTotables    string `json:"-"`
}
