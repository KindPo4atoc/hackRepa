package entity

type ContextTask struct {
	Data []Task `json:"tasks"`
}

func NewContextTask(d []Task) *ContextTask {
	return &ContextTask{
		Data: d,
	}
}
