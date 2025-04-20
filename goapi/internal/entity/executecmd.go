package entity

type ResultExecute struct {
	Status  string     `json:"status"`
	Data    [][]string `json:"data"`
	Columns []string   `json:"columns"`
}
