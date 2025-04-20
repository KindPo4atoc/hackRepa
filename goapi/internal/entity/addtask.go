package entity

type AddTask struct {
	Header      string   `json:"header"`
	Description string   `json:"description"`
	FilesName   []string `json:"filenames"`
	Contents    []string `json:"contents"`
	Author      string   `json:"author"`
}
