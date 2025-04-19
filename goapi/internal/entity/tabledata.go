package entity

type InfoTable struct {
	TableName         string   `json:"name"`
	TableColumns      []string `json:"columns"`
	TableColumnsTypes []string `json:"type"`
}
