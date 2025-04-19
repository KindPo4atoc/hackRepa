package entity

type Command struct {
	Cmd        string `json:"command"`
	TaskNumber int    `json:"task_number"`
}
