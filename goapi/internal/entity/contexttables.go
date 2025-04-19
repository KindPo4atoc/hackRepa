package entity

type ContextTables struct {
	Tables []InfoTable `json:"tables"`
}

func NewContextTables(t []InfoTable) *ContextTables {
	return &ContextTables{
		Tables: t,
	}
}
